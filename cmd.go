package main

import (
	"time"

	"github.com/jimmysawczuk/scm-status/scm"
	"github.com/pkg/errors"

	"flag"
	"fmt"
	"os"
	"path"
)

var version string
var revision string = "dev"
var date string = time.Now().Format(time.RFC3339)

var out = flag.String("out", "", "File in which to (over)write the revision info into (defaults to stdout)")
var executable = flag.String("executable", "", "Path to this executable (used for hooks; default: $GOPATH/bin/scm-status)")
var pretty = flag.Bool("pretty", true, "Indent and format output")
var installHooks = flag.Bool("install-hooks", false, "Install hooks")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s v%s\nUsage:\n", path.Base(os.Args[0]), version)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	if *installHooks {
		err := setup(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	err := getRevision(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return
}

func getRevision(dir string) error {
	parser, err := scm.GetParser(dir)
	if err != nil {
		return errors.Wrapf(err, "get revision (directory: %s)", dir)
	}

	err = scm.ParseAndWrite(parser, scm.OutputConfig{
		Filename: *out,
		Pretty:   *pretty,
	})
	if err != nil {
		return errors.Wrapf(err, "get revision (directory: %s)", dir)
	}

	return nil
}

func setup(dir string) error {
	parser, err := scm.GetParser(dir)
	if err != nil {
		return errors.Wrapf(err, "install hooks (directory: %s)", dir)
	}

	err = scm.InstallHooks(parser, scm.HooksConfig{
		ExecutablePath: *executable,
		OutputConfig: scm.OutputConfig{
			Filename: *out,
			Pretty:   *pretty,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "install hooks")
	}

	return nil
}
