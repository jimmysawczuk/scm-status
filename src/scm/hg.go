package scm

import (
	"errors"
	"flag"
	"os"
	"strings"
	"time"
)

type HgParser struct {
	dir  string
	Type ScmType
}

func NewHgParser(fq_dir string) (*HgParser, error) {
	err := os.Chdir(fq_dir)
	if err != nil {
		return nil, errors.New("Not a valid directory")
	}

	err = os.Chdir(".hg")
	if err != nil {
		return nil, errors.New("Not an hg repository")
	}

	h := new(HgParser)

	h.dir = fq_dir
	h.Type = Hg

	os.Chdir(fq_dir)

	return h, nil
}

func (p *HgParser) Parse() RevisionInfo {
	var rev RevisionInfo

	rev.Type = Hg

	info_joined, _ := runCommand("hg", "log -r . --template {rev}\\n{node|short}\\n{node}\\n{date|rfc822date}\\n"+"{branches}\\n{tags}\\n{author|person}\\n{author|email}", "")

	info := strings.Split(info_joined, "\n")

	message, _ := runCommand("hg", "log -r . --template {desc}", "")

	rev.Dec = info[0]
	rev.Hex = Hex{
		Short: info[1],
		Full:  info[2],
	}

	if commit_date, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", info[3]); err == nil {
		rev.CommitDate = CommitDate(commit_date)
	}

	rev.Author = Author{
		Name:  info[6],
		Email: info[7],
	}

	rev.Message = message

	if info[4] == "" {
		rev.Branch = "default"
	} else {
		rev.Branch = info[4]
	}

	rev.Tags = strings.Split(info[5], " ")

	return rev

}

func (p *HgParser) Setup() {

	executable := flag.Lookup("executable").Value.String()
	out := flag.Lookup("out").Value.String()

	if out == "<STDOUT>" {
		out = "REVISION.json"
	}

	hook := executable + " -out=\"" + out + "\"; # scm-status hook\r\n"

	full_hooks := "\r\n\r\n" + "[hooks]\r\n"
	full_hooks += "post-update = " + hook
	full_hooks += "post-commit = " + hook

	filename := strings.Join([]string{p.Dir(), ".hg", "hgrc"}, path_separator)
	fp, _ := os.OpenFile(filename, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0664)

	_, _ = fp.WriteString(full_hooks)

	fp.Close()
}

func (p *HgParser) Dir() string {
	return p.dir
}
