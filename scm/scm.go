package scm

import (
	"github.com/pkg/errors"

	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ScmType string

const (
	Git ScmType = "git"
	Hg  ScmType = "hg"
)

type RevisionInfo struct {
	Type       ScmType    `json:"type"`
	Dec        string     `json:"dec"`
	Hex        Hex        `json:"hex"`
	Author     Author     `json:"author"`
	CommitDate CommitDate `json:"commit_date"`
	Message    string     `json:"message"`
	Tags       []string   `json:"tags"`
	Branch     string     `json:"branch"`

	WorkingCopy WorkingCopy `json:"working_copy"`

	Extra map[string]interface{} `json:"extra,omitempty"`
}

type Hex struct {
	Short string `json:"short"`
	Full  string `json:"full"`
}

type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CommitDate time.Time

type WorkingCopy struct {
	Added   int `json:"added"`
	Deleted int `json:"deleted"`

	Files []WorkingFile `json:"files"`
}

type WorkingFile struct {
	Name    string `json:"name"`
	Changes int    `json:"changes"`
}

type ScmParser interface {
	Parse() RevisionInfo
	Init()
}

// var pathSeparator string = string(os.PathSeparator)

func ParseAndWrite(scm ScmParser) {
	result := scm.Parse()

	filename := flag.Lookup("out").Value.String()
	pretty := false
	if flag.Lookup("pretty").Value.String() == "true" {
		pretty = true
	}

	if filename == "<STDOUT>" {
		result.WriteToStdout(pretty)
	} else {
		result.Write(filename, pretty)
	}

}

func resolveDir(dir string) (string, error) {
	os.Chdir(dir)

	fqDir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrapf(err, "can't access directory (need execute permissions): %s", dir)
	}

	return fqDir, nil
}

// GetParser returns a new ScmParser that's suited to read the status from the passed directory.
func GetParser(dir string) (ScmParser, error) {

	dir, err := resolveDir(dir)
	if err != nil {
		return nil, nil
	}

	g, err := NewGitParser(dir)
	if g != nil {
		return g, nil
	}

	h, err := NewHgParser(dir)
	if h != nil {
		return h, nil
	}

	return nil, nil
}

func runCommand(exe string, parts ...string) (string, error) {
	cmd := exec.Command(exe, parts...)

	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "runCommand: %s %s", exe, strings.Join(parts, " "))
	}

	return string(output), nil
}

func (ri RevisionInfo) ToJSON(pretty bool) (res []byte, err error) {

	if pretty {
		res, err = json.MarshalIndent(ri, "", "  ")
	} else {
		res, err = json.Marshal(ri)
	}

	return

}

func (ri RevisionInfo) Write(filepath string, pretty bool) {
	json, _ := ri.ToJSON(pretty)

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)

	if err == nil {
		fp.Write(json)

		fp.Close()

	} else {
		fmt.Errorf("%s\n", err)
	}
}

func (ri RevisionInfo) WriteToStdout(pretty bool) {
	json, _ := ri.ToJSON(pretty)

	fmt.Println(bytes.NewBuffer(json).String())
}

func (d CommitDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d)

	str := fmt.Sprintf(`{"date":%q,"timestamp":%d,"iso8601":%q}`, t.Format(time.UnixDate), t.Unix(), t.Format("2006-01-02T15:04:05-07:00"))

	return bytes.NewBufferString(str).Bytes(), nil
}
