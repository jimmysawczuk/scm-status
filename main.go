package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jimmysawczuk/scm-status/scm"
	"github.com/pkg/errors"
)

var (
	version  string
	revision string = "dev"
	date     string = time.Now().Format(time.RFC3339)
)

var (
	out    = flag.String("out", "", "File in which to (over)write the revision info into (defaults to stdout)")
	pretty = flag.Bool("pretty", true, "Indent and format output")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s %s\nUsage:\n", path.Base(os.Args[0]), version)
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

	err := getRevision(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
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
