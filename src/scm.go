package scm

import (
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
	Type        ScmType
	Dec         string
	HexShort    string
	HexFull     string
	AuthorEmail string
	AuthorName  string
	CommitDate  time.Time
	Message     string
	Tags        []string
	Branch      string

	Extra map[string]interface{}
}

type ScmParser interface {
	Parse() RevisionInfo
	Setup()
	Dir() string
}

func ParseAndWrite(scm ScmParser) {
	result := scm.Parse()

	filename := flag.Lookup("out").Value.String()

	result.Write(scm.Dir() + "/" + filename)
}

func resolveDir(dir string) (fq_dir string, err error) {
	os.Chdir(dir)

	fq_dir, err = os.Getwd()

	if err != nil {
		fmt.Errorf("Can't resolve working directory %s; do you have execute permissions?\n", dir)
		return "", err
	}

	return fq_dir, nil
}

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

func runCommand(exe string, args string, dir string) (string, error) {
	parts := strings.Split(args, " ")
	cmd := exec.Command(exe, parts...)

	output, err := cmd.Output()

	if err != nil {
		fmt.Errorf("%s", err)
		return "", err
	}

	str := bytes.NewBuffer(output).String()

	return str, nil
}

func (ri RevisionInfo) toMap() map[string]interface{} {
	payload := make(map[string]interface{})

	payload["type"] = ri.Type
	payload["dec"] = ri.Dec

	hex := make(map[string]interface{})
	hex["short"] = ri.HexShort
	hex["full"] = ri.HexFull
	payload["hex"] = hex

	author := make(map[string]interface{})
	author["name"] = ri.AuthorName
	author["email"] = ri.AuthorEmail
	payload["author"] = author

	payload["branch"] = ri.Branch
	payload["tags"] = ri.Tags

	payload["message"] = ri.Message

	payload["commit_date"] = ri.CommitDate.Format(time.UnixDate)
	payload["commit_timestamp"] = ri.CommitDate.Unix()

	for idx, val := range ri.Extra {
		payload[idx] = val
	}

	return payload
}

func (ri RevisionInfo) ToJSON() ([]byte, error) {
	ri_map := ri.toMap()

	return json.Marshal(ri_map)
}

func (ri RevisionInfo) Write(filepath string) {
	json, _ := ri.ToJSON()

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)

	if err == nil {
		fp.Write(json)

		fp.Close()

	} else {
		fmt.Errorf("%s\n", err)
	}
}
