package scm

import (
	"github.com/pkg/errors"

	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var gitFileDiffRegexp = regexp.MustCompile(`^(.+?)\s+\|\s+(\d+) [\+|\-]*$`)
var gitSummaryDiffRegexp = regexp.MustCompile(`^(?:(\d+) files? changed)?(?:, )?(?:(\d+) insertions?\(\+\))?(?:, )?(?:(\d+) deletions?\(\-\))?$`)
var gitHookTmpl = template.Must(
	template.New("githooks").Parse(`{{ .ExecutablePath }} {{ with .OutputConfig.Filename }}{{ printf "-out=%q" . }}{{ end }} {{ with .OutputConfig.Pretty }}{{ printf "-pretty=%t" . }}{{ end }}; # installed by scm-status (github.com/jimmysawczuk/scm-status)`),
)

type gitParser struct {
	dir string
}

func newGitParser(dir string) (*gitParser, error) {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	if err != nil {
		return nil, errors.New("Not a git repository")
	}

	g := &gitParser{
		dir: dir,
	}

	return g, nil
}

// Parse returns a RevisionInfo struct after parsing the provided directory for working copy information.
func (p *gitParser) Parse() (Snapshot, error) {
	version, _ := runCommand(p.dir, "git", "--version")
	rawBranch, _ := runCommand(p.dir, "git", "branch", "--contains", "HEAD")
	tagsJoined, _ := runCommand(p.dir, "git", "tag", "--contains", "HEAD")
	metaJoined, _ := runCommand(p.dir, "git", "log", "-1", "--pretty=format:%h%n%h%n%H%n%ci%n%an%n%ae%n%p%n%P%n%s")

	useOldGitSyntax := !meetsVersion(version, "1.7.2")
	rawCommitMessage := ""
	if useOldGitSyntax {
		rawCommitMessage, _ = runCommand(p.dir, "git", "log", "-1", "--pretty=format:%s%n%b")
	} else {
		rawCommitMessage, _ = runCommand(p.dir, "git", "log", "-1", "--pretty=format:%B")
	}

	workingCopyStats, _ := runCommand(p.dir, "git", "diff", "HEAD", "--stat")

	rev := Snapshot{
		Type:  "git",
		Extra: make(map[string]interface{}),
	}

	version = strings.Replace(version, "git version ", "", -1)

	currentBranch, allBranches := extractBranches(rawBranch)

	tagsJoined = strings.TrimSpace(tagsJoined)
	tags := []string{}
	if len(tagsJoined) > 0 {
		tags = strings.Split(tagsJoined, "\n")
	}

	meta := strings.Split(metaJoined, "\n")

	commitMessage := strings.TrimSpace(rawCommitMessage)

	rev.Type = "git"
	rev.Message = commitMessage
	rev.Tags = tags
	rev.Branch = currentBranch
	rev.Extra["branches"] = allBranches

	rev.Dec = meta[0]
	rev.Hex = Hex{
		Short: meta[1],
		Full:  meta[2],
	}

	rev.Author = Author{
		Name:  meta[4],
		Email: meta[5],
	}

	if date, err := time.Parse("2006-01-02 15:04:05 -0700", meta[3]); err == nil {
		rev.CommitDate = CommitDate(date)
	}

	shortParents := strings.Split(strings.TrimSpace(meta[6]), " ")
	longParents := strings.Split(strings.TrimSpace(meta[7]), " ")

	subject := strings.TrimSpace(meta[8])

	parents := []map[string]interface{}{}
	for i := range shortParents {
		parent := make(map[string]interface{})
		parent["short"] = shortParents[i]
		parent["full"] = longParents[i]
		parents = append(parents, parent)
	}

	rev.Extra["parents"] = parents
	rev.Extra["subject"] = subject

	rev.UncommittedChanges = parseDiffStat(workingCopyStats)

	return rev, nil

}

func (p *gitParser) InstallHooks(config HooksConfig) error {

	buf := &bytes.Buffer{}
	gitHookTmpl.Execute(buf, config)
	hook := buf.String()

	hookDir := filepath.Join(p.dir, ".git", "hooks")

	filenames := []string{
		filepath.Join(hookDir, "post-checkout"),
		filepath.Join(hookDir, "post-merge"),
		filepath.Join(hookDir, "post-commit"),
	}

	for _, filename := range filenames {
		fp, _ := os.OpenFile(filename, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0775)
		fp.WriteString(hook)
		fp.Close()
	}

	return nil
}

func meetsVersion(test, req string) bool {

	intConvert := func(a []string) []int {
		b := make([]int, len(a))

		for i, v := range a {
			b[i], _ = strconv.Atoi(v)
		}

		return b
	}

	testVersion := intConvert(strings.Split(strings.TrimSpace(test), "."))
	reqVersion := intConvert(strings.Split(strings.TrimSpace(req), "."))

	// make both versions same length for easier comparison
	if len(testVersion) > len(reqVersion) {
		for i := len(reqVersion); i < len(testVersion); i++ {
			reqVersion = append(reqVersion, 0)
		}
	} else if len(reqVersion) > len(testVersion) {
		for i := len(testVersion); i < len(reqVersion); i++ {
			testVersion = append(testVersion, 0)
		}
	}

	// compare each line and see where we are
	series := make([]int, len(reqVersion))
	lastPos, lastNeg := -1, -1

	for i, j := len(reqVersion)-1, len(reqVersion)-1; i >= 0; i, j = i-1, j-1 {
		if testVersion[i] > reqVersion[i] {
			series[j] = 1
			lastPos = j
		} else if testVersion[i] < reqVersion[i] {
			series[j] = -1
			lastNeg = j
		} else {
			series[j] = 0
		}
	}

	return (lastPos < lastNeg && lastPos >= 0) || (lastPos == -1 && lastNeg == -1)
}

func extractBranches(rawBranch string) (primary string, all []string) {
	rawBranch = strings.Replace(rawBranch, "\r\n", "\n", -1)
	temp := strings.Split(rawBranch, "\n")
	for _, b := range temp {
		b = strings.TrimSpace(b)

		if strings.HasPrefix(b, "* ") {
			b = strings.Replace(b, "* ", "", -1)
			primary = b
		}

		if len(b) > 0 {
			all = append(all, b)
		}
	}

	return
}

func parseDiffStat(in string) UncommittedChanges {
	in = strings.Replace(in, "\r\n", "\n", -1)
	lines := strings.Split(in, "\n")

	uc := UncommittedChanges{
		Files: []UncommittedFile{},
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		} else if gitFileDiffRegexp.MatchString(line) {
			match := gitFileDiffRegexp.FindAllStringSubmatch(line, -1)

			changes, _ := strconv.ParseInt(match[0][2], 10, 32)

			uc.Files = append(uc.Files, UncommittedFile{
				Name:    match[0][1],
				Changes: int(changes),
			})

		} else if gitSummaryDiffRegexp.MatchString(line) {
			match := gitSummaryDiffRegexp.FindAllStringSubmatch(line, -1)

			additions, _ := strconv.ParseInt(match[0][2], 10, 32)
			deletions, _ := strconv.ParseInt(match[0][3], 10, 32)

			uc.Added = int(additions)
			uc.Deleted = int(deletions)
		}
	}

	return uc
}
