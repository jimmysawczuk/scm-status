package scm

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

type GitParser struct {
	dir  string
	Type ScmType
}

func NewGitParser(fq_dir string) (*GitParser, error) {
	err := os.Chdir(fq_dir)
	if err != nil {
		return nil, errors.New("Not a valid directory")
	}

	err = os.Chdir(".git")
	if err != nil {
		return nil, errors.New("Not a git repository")
	}

	g := new(GitParser)

	g.dir = fq_dir
	g.Type = Git

	os.Chdir(fq_dir)

	return g, nil
}

func (p *GitParser) Parse() RevisionInfo {
	var rev RevisionInfo

	version, _ := runCommand("git", "--version", "")
	version = strings.Replace(version, "git version ", "", -1)

	branch_raw, _ := runCommand("git", "branch --contains HEAD", "")
	branch_raw = strings.TrimSpace(branch_raw)
	branch := strings.Replace(branch_raw, "* ", "", -1)

	tags_joined, _ := runCommand("git", "tag --contains HEAD", "")
	tags_joined = strings.TrimSpace(tags_joined)
	tags := make([]string, 0)
	if len(tags_joined) > 0 {
		tags = strings.Split(tags_joined, "\n")
	}

	meta_joined, _ := runCommand("git", "log -1 --pretty=format:%h%n%h%n%H%n%ci%n%an%n%ae%n%p%n%P%n%s", "")
	meta := strings.Split(meta_joined, "\n")

	use_old_git := !meetsVersion(version, "1.7.2")
	commit_message_raw := ""
	if use_old_git {
		commit_message_raw, _ = runCommand("git", "log -1 --pretty=format:%s%n%b", "")
	} else {
		commit_message_raw, _ = runCommand("git", "log -1 --pretty=format:%B", "")
	}

	commit_message := strings.TrimSpace(commit_message_raw)

	rev.Type = Git
	rev.Message = commit_message
	rev.Tags = tags
	rev.Branch = branch

	rev.Dec = meta[0]
	rev.HexShort = meta[1]
	rev.HexFull = meta[2]
	rev.AuthorName = meta[4]
	rev.AuthorEmail = meta[5]

	rev.CommitDate, _ = time.Parse("2006-01-02 15:04:05 -0700", meta[3])

	rev.Extra = make(map[string]interface{})

	short_parents := strings.Split(strings.TrimSpace(meta[6]), " ")
	long_parents := strings.Split(strings.TrimSpace(meta[7]), " ")

	subject := strings.TrimSpace(meta[8])

	parents := make([]map[string]interface{}, 0)
	var parent map[string]interface{}
	for i, _ := range short_parents {
		parent = make(map[string]interface{})
		parent["short"] = short_parents[i]
		parent["full"] = long_parents[i]
		parents = append(parents, parent)
	}

	rev.Extra["parents"] = parents
	rev.Extra["subject"] = subject

	return rev

}

func (p *GitParser) Setup() {
	executable := flag.Lookup("executable").Value.String()
	out := flag.Lookup("out").Value.String()

	if out == "<STDOUT>" {
		out = "REVISION.json"
	}

	hook := executable + " -out=\"" + out + "\"; # scm-status hook\r\n"

	filenames := []string{
		p.Dir() + "/.git/hooks/post-checkout",
		p.Dir() + "/.git/hooks/post-merge",
		p.Dir() + "/.git/hooks/post-commit",
	}

	for _, filename := range filenames {
		fp, _ := os.OpenFile(filename, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0775)

		_, _ = fp.WriteString(hook)

		fp.Close()
	}
}

func (p *GitParser) Dir() string {
	return p.dir
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
