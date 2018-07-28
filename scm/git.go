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

func (g *gitParser) Build(dir string) error {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	if err != nil {
		switch errors.Cause(err).(type) {
		case *os.PathError:
			return errNotRepository
		default:
			return errors.Wrap(err, "access .git directory")
		}
	}

	g.dir = dir

	return nil
}

// Parse returns a RevisionInfo struct after parsing the provided directory for working copy information.
func (g *gitParser) Parse() (Snapshot, error) {
	rev := Snapshot{
		Type:  "git",
		Extra: make(map[string]interface{}),
	}

	version, _ := runCommand(g.dir, "git", "--version")
	version = strings.Replace(version, "git version ", "", -1)

	rawBranch, _ := runCommand(g.dir, "git", "branch", "--contains", "HEAD")
	tagsJoined, _ := runCommand(g.dir, "git", "tag", "--contains", "HEAD")
	metaJoined, _ := runCommand(g.dir, "git", "log", "-1", "--pretty=format:%h%n%h%n%H%n%ci%n%an%n%ae%n%p%n%P%n%s")

	useOldGitSyntax := versionNumber(version).meets(versionNumber("1.7.2"))
	rawCommitMessage := ""
	if useOldGitSyntax {
		rawCommitMessage, _ = runCommand(g.dir, "git", "log", "-1", "--pretty=format:%s%n%b")
	} else {
		rawCommitMessage, _ = runCommand(g.dir, "git", "log", "-1", "--pretty=format:%B")
	}

	workingCopyStats, _ := runCommand(g.dir, "git", "diff", "HEAD", "--stat")

	currentBranch, allBranches := extractBranches(rawBranch)

	tagsJoined = strings.TrimSpace(tagsJoined)
	tags := []string{}
	if len(tagsJoined) > 0 {
		tags = strings.Split(tagsJoined, "\n")
	}

	meta := strings.Split(metaJoined, "\n")
	if len(meta) < 9 {
		return Snapshot{}, errors.New("blank repository; no commits")
	}

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
		rev.Date = date
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

func (g *gitParser) InstallHooks(config HooksConfig) error {

	buf := &bytes.Buffer{}
	gitHookTmpl.Execute(buf, config)
	hook := buf.String()

	hookDir := filepath.Join(g.dir, ".git", "hooks")

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
