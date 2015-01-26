package scm

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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
	Setup()
	Dir() string
}

var path_separator string = string(os.PathSeparator)

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
		result.Write(scm.Dir()+path_separator+filename, pretty)
	}

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

func (wc WorkingCopy) MarshalJSON() ([]byte, error) {
	if len(wc.Files) == 0 {
		wc.Files = make([]WorkingFile, 0)
	}

	files, _ := json.Marshal(wc.Files)

	str := fmt.Sprintf(`{"added":%d,"deleted":%d,"files":%s}`, wc.Added, wc.Deleted, files)
	return bytes.NewBufferString(str).Bytes(), nil
}

func parseDiffStat(in string) (wc WorkingCopy) {

	in = strings.Replace(in, "\r\n", "\n", -1)
	lines := strings.Split(in, "\n")

	file_re := regexp.MustCompile(`^(.+?)\s+\|\s+(\d+) [\+|\-]*$`)
	summary_re := regexp.MustCompile(`^(?:(\d+) files? changed)?(?:, )?(?:(\d+) insertions?\(\+\))?(?:, )?(?:(\d+) deletions?\(\-\))?$`)

	wc.Files = []WorkingFile{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		} else if file_re.MatchString(line) {
			match := file_re.FindAllStringSubmatch(line, -1)

			changes, _ := strconv.ParseInt(match[0][2], 10, 32)

			wc.Files = append(wc.Files, WorkingFile{
				Name:    match[0][1],
				Changes: int(changes),
			})

		} else if summary_re.MatchString(line) {
			match := summary_re.FindAllStringSubmatch(line, -1)

			additions, _ := strconv.ParseInt(match[0][2], 10, 32)
			deletions, _ := strconv.ParseInt(match[0][3], 10, 32)

			wc.Added = int(additions)
			wc.Deleted = int(deletions)
		}
	}

	return wc
}
