package scm

import(
	"time"
	"strings"
//	"fmt"
)

type GitParser struct {
	Dir string
}

func NewGitParser() *GitParser {
	g := new(GitParser)

	return g
}

func (p *GitParser) Parse() RevisionInfo {
	var rev RevisionInfo

	branch_raw := runCommand("git", "branch --contains HEAD", "")
	branch_raw = strings.TrimSpace(branch_raw)
	branch := strings.Replace(branch_raw, "* ", "", -1)

	tags_joined := runCommand("git", "tag --contains HEAD", "")
	tags_joined = strings.TrimSpace(tags_joined)
	tags := make([]string, 0)
	if len(tags_joined) > 0 {
		tags = strings.Split(tags_joined, "\n")	
	}
	
	meta_joined := runCommand("git", "log -1 --pretty=format:%h%n%h%n%H%n%ci%n%an%n%ae", "")
	meta := strings.Split(meta_joined, "\n")

	commit_message := runCommand("git", "log -1 --pretty=format:%B", "")

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

	return rev

}