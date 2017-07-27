package scm

import (
	"github.com/pkg/errors"

	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Parser is the interface that's used to query the current working directory for revision info. Currently git and hg are implemented.
type Parser interface {
	Build(string) error
	Parse() (Snapshot, error)
	InstallHooks(HooksConfig) error
}

// HooksConfig holds configuration information for how hooks should be installed in a working copy to generate snapshots automatically.
type HooksConfig struct {
	OutputConfig
	ExecutablePath string
}

// OutputConfig holds configuration information for how a working copy snapshot should be output.
type OutputConfig struct {
	Filename string
	Pretty   bool
}

// InstallHooks uses the provided ScmParser and installs hooks as specified in config.
func InstallHooks(scm Parser, config HooksConfig) error {
	executable := filepath.Join(os.Getenv("GOPATH"), "bin/scm-status")
	if config.ExecutablePath != "" {
		executable, _ = filepath.Abs(config.ExecutablePath)
	}

	filename := ""
	if config.OutputConfig.Filename != "" {
		filename = config.OutputConfig.Filename
	}

	pretty := config.OutputConfig.Pretty

	err := scm.InstallHooks(HooksConfig{
		ExecutablePath: executable,
		OutputConfig: OutputConfig{
			Filename: filename,
			Pretty:   pretty,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "install hooks")
	}
	return nil
}

// ParseAndWrite uses the provided ScmParser to get a snapshot and outputs it as specified in config.
func ParseAndWrite(scm Parser, config OutputConfig) error {
	result, err := scm.Parse()
	if err != nil {
		return errors.Wrap(err, "parseandwrite")
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
			case errNotRepository:
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
		return "", errInvalidDirectory
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
