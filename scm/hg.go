package scm

import (
	"errors"
	"flag"
	"os"
	"path"
	"strings"
	"time"
)

type HgParser struct {
	Dir  string
	Type ScmType
}

func NewHgParser(fqDir string) (*HgParser, error) {
	err := os.Chdir(fqDir)
	if err != nil {
		return nil, errors.New("Not a valid directory")
	}

	err = os.Chdir(".hg")
	if err != nil {
		return nil, errors.New("Not an hg repository")
	}

	h := &HgParser{
		Dir:  fqDir,
		Type: Hg,
	}

	os.Chdir(fqDir)

	return h, nil
}

func (p *HgParser) Parse() RevisionInfo {
	var rev RevisionInfo

	rev.Type = Hg

	rawInfo, _ := runCommand("hg", "log", "-r", ".", "--template", "{rev}\\n{node|short}\\n{node}\\n{date|rfc822date}\\n"+"{branches}\\n{tags}\\n{author|person}\\n{author|email}")
	info := strings.Split(rawInfo, "\n")

	message, _ := runCommand("hg", "log", "-r", ".", "--template", "{desc}")

	rev.Dec = info[0]
	rev.Hex = Hex{
		Short: info[1],
		Full:  info[2],
	}

	if date, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", info[3]); err == nil {
		rev.CommitDate = CommitDate(date)
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

func (p *HgParser) Init() {

	executable := flag.Lookup("executable").Value.String()
	out := flag.Lookup("out").Value.String()

	if out == "<STDOUT>" {
		out = "REVISION.json"
	}

	hook := executable + " -out=\"" + out + "\"; # scm-status hook\r\n"

	fullHooks := "\r\n\r\n" + "[hooks]\r\n"
	fullHooks += "post-update = " + hook
	fullHooks += "post-commit = " + hook

	filename := path.Join(p.Dir, ".hg", "hgrc")
	fp, _ := os.OpenFile(filename, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0664)

	_, _ = fp.WriteString(fullHooks)

	fp.Close()
}
