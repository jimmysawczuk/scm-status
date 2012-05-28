package scm

import (
	"flag"
	//	"fmt"
	"os"
	"strings"
	"time"
)

type GitParser struct {
	dir  string
	Type ScmType
}

func NewGitParser(fq_dir string) *GitParser {
	g := new(GitParser)

	g.dir = fq_dir
	g.Type = Git

	return g
}

func (p *GitParser) Parse() RevisionInfo {
	var rev RevisionInfo

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

	commit_message_raw, _ := runCommand("git", "log -1 --pretty=format:%B", "")
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

	hook := executable + "; # scm-status hook\r\n"

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
