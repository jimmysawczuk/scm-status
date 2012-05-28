package scm

import (
	"bytes"
	"encoding/json"
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

func Write(scm ScmParser) {
	fmt.Println("got here")
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

	err = os.Chdir(dir)

	err = os.Chdir(".git")
	if err == nil {
		g := NewGitParser(dir)
		os.Chdir(dir)
		return g, nil
	}

	return nil, nil
}

func runCommand(exe string, args string, dir string) (string, error) {
	parts := strings.Split(args, " ")
	cmd := exec.Command(exe, parts...) // "git", "branch", "--contains", "HEAD")

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
