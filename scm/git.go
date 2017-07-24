package scm

import (
	"errors"
	"flag"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var gitFileDiffRegexp = regexp.MustCompile(`^(.+?)\s+\|\s+(\d+) [\+|\-]*$`)
var gitSummaryDiffRegexp = regexp.MustCompile(`^(?:(\d+) files? changed)?(?:, )?(?:(\d+) insertions?\(\+\))?(?:, )?(?:(\d+) deletions?\(\-\))?$`)

type GitParser struct {
	Dir  string
	Type ScmType
}

// NewGitParser returns a GitParser for the provided directory, if valid.
func NewGitParser(fqDir string) (*GitParser, error) {
	err := os.Chdir(fqDir)
	if err != nil {
		return nil, errors.New("Not a valid directory")
	}

	err = os.Chdir(".git")
	if err != nil {
		return nil, errors.New("Not a git repository")
	}

	g := &GitParser{
		Dir:  fqDir,
		Type: Git,
	}

	os.Chdir(fqDir)

	return g, nil
}

// Parse returns a RevisionInfo struct after parsing the provided directory for working copy information.
func (p *GitParser) Parse() RevisionInfo {
	var rev RevisionInfo
	rev.Extra = make(map[string]interface{})

	version, _ := runCommand("git", "--version")
	version = strings.Replace(version, "git version ", "", -1)

	// Need to make this smarter about parsing the active branch.
	rawBranch, _ := runCommand("git", "branch", "--contains", "HEAD")
	currentBranch, allBranches := extractBranches(rawBranch)

	tagsJoined, _ := runCommand("git", "tag", "--contains", "HEAD")
	tagsJoined = strings.TrimSpace(tagsJoined)
	tags := []string{}
	if len(tagsJoined) > 0 {
		tags = strings.Split(tagsJoined, "\n")
	}

	metaJoined, _ := runCommand("git", "log", "-1", "--pretty=format:%h%n%h%n%H%n%ci%n%an%n%ae%n%p%n%P%n%s")
	meta := strings.Split(metaJoined, "\n")

	useOldGitSyntax := !meetsVersion(version, "1.7.2")
	rawCommitMessage := ""
	if useOldGitSyntax {
		rawCommitMessage, _ = runCommand("git", "log", "-1", "--pretty=format:%s%n%b")
	} else {
		rawCommitMessage, _ = runCommand("git", "log", "-1", "--pretty=format:%B")
	}
	commitMessage := strings.TrimSpace(rawCommitMessage)

	workingCopyStats, _ := runCommand("git", "diff", "HEAD", "--stat")

	rev.Type = Git
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

	rev.WorkingCopy = parseDiffStat(workingCopyStats)

	return rev

}

func (p *GitParser) Init() {
	executable := flag.Lookup("executable").Value.String()
	out := flag.Lookup("out").Value.String()

	if out == "<STDOUT>" {
		out = "REVISION.json"
	}

	hook := executable + " -out=\"" + out + "\"; # scm-status hook\r\n"

	hookDir := path.Join(p.Dir, ".git", "hooks")

	filenames := []string{
		path.Join(hookDir, "post-checkout"),
		path.Join(hookDir, "post-merge"),
		path.Join(hookDir, "post-commit"),
	}

	for _, filename := range filenames {
		fp, _ := os.OpenFile(filename, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0775)
		fp.WriteString(hook)
		fp.Close()
	}
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
	last_pos, last_neg := -1, -1 // last_neut := -1, -1, -1

	for i, j := len(reqVersion)-1, len(reqVersion)-1; i >= 0; i, j = i-1, j-1 {
		if testVersion[i] > reqVersion[i] {
			series[j] = 1
			last_pos = j
		} else if testVersion[i] < reqVersion[i] {
			series[j] = -1
			last_neg = j
		} else {
			series[j] = 0
			// last_neut = j
		}
	}

	valid_version := (last_pos < last_neg && last_pos >= 0) || (last_pos == -1 && last_neg == -1)

	return valid_version
}

func extractBranches(branch_raw string) (primary_branch string, all_branches []string) {
	branch_raw = strings.Replace(branch_raw, "\r\n", "\n", -1)
	temp := strings.Split(branch_raw, "\n")
	for _, b := range temp {
		b = strings.TrimSpace(b)

		if strings.HasPrefix(b, "* ") {
			b = strings.Replace(b, "* ", "", -1)
			primary_branch = b
		}

		if len(b) > 0 {
			all_branches = append(all_branches, b)
		}
	}

	return
}

func parseDiffStat(in string) (wc WorkingCopy) {

	in = strings.Replace(in, "\r\n", "\n", -1)
	lines := strings.Split(in, "\n")

	wc.Files = []WorkingFile{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		} else if gitFileDiffRegexp.MatchString(line) {
			match := gitFileDiffRegexp.FindAllStringSubmatch(line, -1)

			changes, _ := strconv.ParseInt(match[0][2], 10, 32)

			wc.Files = append(wc.Files, WorkingFile{
				Name:    match[0][1],
				Changes: int(changes),
			})

		} else if gitSummaryDiffRegexp.MatchString(line) {
			match := gitSummaryDiffRegexp.FindAllStringSubmatch(line, -1)

			additions, _ := strconv.ParseInt(match[0][2], 10, 32)
			deletions, _ := strconv.ParseInt(match[0][3], 10, 32)

			wc.Added = int(additions)
			wc.Deleted = int(deletions)
		}
	}

	return wc
}
