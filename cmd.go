package main

import (
	"github.com/jimmysawczuk/scm-status/scm"
	"github.com/pkg/errors"

	"flag"
	"fmt"
	"os"
	"path"
)

var version = "1.2.0"

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s v%s\nUsage:\n", path.Base(os.Args[0]), version)
		flag.PrintDefaults()
	}

	setupFlags()
	handle()
}

func setupFlags() {
	flag.String("out", "<STDOUT>", "File in which to (over)write the revision info into (defaults to stdout)")
	flag.String("executable", "/usr/local/bin/scm-status", "Path at which this program can be executed (for hooks)")
	flag.Bool("old-git", false, "Obsolete, now ignored")
	flag.Bool("pretty", true, "Set to false to output compressed JSON")
}

func handle() {
	flag.Parse()

	args := flag.Args()

	switch {
	case len(args) == 1 && args[len(args)-1] == "setup":
		setup(args[0])

	case len(args) == 0:
		parseRevision(".")

	default:
		parseRevision(args[0])
	}
}

func getScmParser(dir string) (scm.ScmParser, error) {
	if err := os.Chdir(dir); err != nil {
		return nil, errors.Wrapf(err, "can't access directory (need execute permissions): %s", dir)
	}

	parser, err := scm.GetParser(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "build parser (directory: %s)", dir)
	}

	return parser, nil
}

func parseRevision(dir string) error {

	parser, err := getScmParser(dir)
	if err != nil {
		return err
	}

	scm.ParseAndWrite(parser)
	return nil
}

func setup(dir string) {
	parser, err := getScmParser(".")

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	if parser != nil {
		parser.Init()
	}
}
