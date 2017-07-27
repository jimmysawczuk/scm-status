package scm

import (
	"errors"
)

var errNotRepository = errors.New("Not a repository")
var errInvalidDirectory = errors.New("Invalid directory")
