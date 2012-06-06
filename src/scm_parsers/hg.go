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
	h.Type = Git

	return h, nil
}

func (p *HgParser) Parse() RevisionInfo {
	var rev RevisionInfo

	rev.Type = Hg

	info_joined, _ := runCommand("hg", "log -r . --template {rev}\\n{node|short}\\n{node}\\n{date|isodate}\\n"+"{branches}\\n{tags}\\n{author|person}\\n{author|email}", "")

	info := strings.Split(info_joined, "\n")

	message, _ := runCommand("hg", "log -r . --template {desc}", "")

	rev.Dec = info[0]
	rev.HexShort = info[1]
	rev.HexFull = info[2]
	rev.CommitDate, _ = time.Parse("2006-01-02 15:04 -0700", info[3])

	rev.AuthorName = info[6]
	rev.AuthorEmail = info[7]

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

	hook := executable + "; # scm-status hook\r\n"

	full_hooks := "\r\n\r\n" + "[hooks]\r\n"
	full_hooks += "post-update = " + hook
	full_hooks += "post-commit = " + hook

	filename := p.Dir() + "/.hg/hgrc"
	fp, _ := os.OpenFile(filename, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0775)

	_, _ = fp.WriteString(full_hooks)

	fp.Close()
}

func (p *HgParser) Dir() string {
	return p.dir
}
