package scm

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type hgParser struct {
	dir string
}

func (hg *hgParser) Build(dir string) error {
	_, err := os.Stat(filepath.Join(dir, ".hg"))
	if err != nil {
		switch errors.Cause(err).(type) {
		case *os.PathError:
			return ErrNotRepository
		default:
			return errors.Wrap(err, "access .hg directory")
		}
	}

	hg.dir = dir

	return nil
}

func (hg *hgParser) Parse() (Snapshot, error) {
	rev := Snapshot{
		Type: "hg",
	}

	rawInfo, _ := runCommand(hg.dir, "hg", "log", "-r", ".", "--template", "{rev}\\n{node|short}\\n{node}\\n{date|rfc822date}\\n"+"{branches}\\n{tags}\\n{author|person}\\n{author|email}")
	info := strings.Split(rawInfo, "\n")

	message, _ := runCommand(hg.dir, "hg", "log", "-r", ".", "--template", "{desc}")

	rev.Dec = info[0]
	rev.Hex = Hex{
		Short: info[1],
		Full:  info[2],
	}

	if date, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", info[3]); err == nil {
		rev.Date = date
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

	return rev, nil
}
