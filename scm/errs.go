package scm

import (
	"errors"
)

var (
	ErrNotRepository    = errors.New("not a repository")
	ErrInvalidDirectory = errors.New("invalid directory")
)
