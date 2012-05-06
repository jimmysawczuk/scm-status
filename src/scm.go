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
}

func Write(scm ScmParser) {
	fmt.Println("got here")
}

func GetParser(dir string) (ScmParser, error) {

	err := os.Chdir(dir)

	if err != nil {
		return nil, fmt.Errorf("Can't access %s; do you have execute permissions?\n", dir)
	}

	err = os.Chdir(".git")
	if err == nil {
		g := NewGitParser()
		return g, nil
	}

	return nil, nil
}

func runCommand(exe string, args string, dir string) string {
	parts := strings.Split(args, " ")
	cmd := exec.Command(exe, parts...) // "git", "branch", "--contains", "HEAD")

	output, _ := cmd.Output()

	str := bytes.NewBuffer(output).String()

	return str
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
