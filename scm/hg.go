package scm

import (
	"github.com/pkg/errors"

	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var hgHookTmpl = template.Must(
	template.New("hghooks").Parse(`{{ .ExecutablePath }} {{ with .OutputConfig.Filename }}{{ printf "-out=%q" . }}{{ end }} {{ with .OutputConfig.Pretty }}{{ printf "-pretty=%t" . }}{{ end }}; # installed by scm-status (github.com/jimmysawczuk/scm-status)`),
)

type hgParser struct {
	dir string
}

func (hg *hgParser) Build(dir string) error {
	_, err := os.Stat(filepath.Join(dir, ".hg"))
	if err != nil {
		switch errors.Cause(err).(type) {
		case *os.PathError:
			return errNotRepository
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

func (hg *hgParser) InstallHooks(config HooksConfig) error {

	buf := &bytes.Buffer{}
	hgHookTmpl.Execute(buf, config)
	hook := buf.String()

	fullHooks := "\r\n\r\n" + "[hooks]\r\n"
	fullHooks += "post-update = " + hook + "\r\n"
	fullHooks += "post-commit = " + hook + "\r\n"

	filename := filepath.Join(hg.dir, ".hg", "hgrc")
	fp, _ := os.OpenFile(filename, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0664)
	_, _ = fp.WriteString(fullHooks)
	fp.Close()

	return nil
}
