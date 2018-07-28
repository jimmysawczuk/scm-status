package scm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot contains information about the current revision of the working copy
type Snapshot struct {

	// The type of repository (currently either git or hg)
	Type string `json:"type"`

	// The decimal revision (currently only relevant to hg; with git repositories, this field is the same as the short hex)
	Dec string `json:"dec"`

	// The hexadecimal revision (both the long and short)
	Hex Hex `json:"hex"`

	// The author of the commit
	Author Author `json:"author"`

	// The date and time of the commit
	Date time.Time `json:"date"`

	// The commit message.
	Message string `json:"message"`

	// Any tags attached to this commit.
	Tags []string `json:"tags"`

	// The working copy's current branch.
	Branch string `json:"branch"`

	// A summary of all uncommitted changes to the repository.
	UncommittedChanges UncommittedChanges `json:"uncommitted"`

	// Any extra information that may not be relevant to every SCM.
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// Hex holds a short and full hexadecimal revision hash for one revision
type Hex struct {
	Short string `json:"short"`
	Full  string `json:"full"`
}

// Author contains the name and e-mail of an author
type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UncommittedChanges contains information about the current uncommitted changes
type UncommittedChanges struct {
	Added   int `json:"added"`
	Deleted int `json:"deleted"`

	Files []UncommittedFile `json:"files"`
}

// UncommittedFile contains a diff summary for just one file
type UncommittedFile struct {
	Name    string `json:"name"`
	Changes int    `json:"changes"`
}

// ToJSON returns a JSON-string of the snapshot, with or without formatting
func (ri Snapshot) ToJSON(pretty bool) (res []byte, err error) {
	if pretty {
		return json.MarshalIndent(ri, "", "  ")
	}

	return json.Marshal(ri)
}

// Write writes a snapshot's JSON representation to a file
func (ri Snapshot) Write(filepath string, pretty bool) error {
	json, _ := ri.ToJSON(pretty)

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	fp.Write(json)
	fp.Close()
	return nil

}

// WriteToStdout writes a snapshot's JSON representation to stdout
func (ri Snapshot) WriteToStdout(pretty bool) {
	json, _ := ri.ToJSON(pretty)

	fmt.Println(bytes.NewBuffer(json).String())
}
