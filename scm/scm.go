package scm

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Parser is the interface that's used to query the current working directory for revision info. Currently git and hg are implemented.
type Parser interface {
	Build(string) error
	Parse() (Snapshot, error)
	// InstallHooks(HooksConfig) error
}

// OutputConfig holds configuration information for how a working copy snapshot should be output.
type OutputConfig struct {
	Filename string
	Pretty   bool
}

// ParseAndWrite uses the provided ScmParser to get a snapshot and outputs it as specified in config.
func ParseAndWrite(scm Parser, config OutputConfig) error {
	result, err := scm.Parse()
	if err != nil {
		return errors.Wrap(err, "scm: parse")
	}

	filename := ""
	if config.Filename != "" {
		filename = config.Filename
	}

	pretty := config.Pretty

	if filename == "" {
		result.WriteToStdout(pretty)
	} else {
		result.Write(filename, pretty)
	}

	return nil
}

// GetParser returns a new ScmParser that's suited to read the status from the passed directory.
func GetParser(dir string) (Parser, error) {
	dir, err := resolveDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve directory")
	}

	parsers := []Parser{
		&gitParser{},
		&hgParser{},
	}

	for _, parser := range parsers {
		if err := parser.Build(dir); err != nil {
			switch errors.Cause(err) {
			case ErrNotRepository:
				continue
			default:
				return nil, errors.Wrapf(err, "build parser (%T)", parser)
			}
		}
		return parser, nil
	}

	return nil, errors.New("couldn't find a valid repository")
}

func resolveDir(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", ErrInvalidDirectory
	}

	return absDir, nil
}

func runCommand(dir, exe string, parts ...string) (string, error) {
	cmd := exec.Command(exe, parts...)
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "runCommand: %s %s", exe, strings.Join(parts, " "))
	}

	return string(output), nil
}
