package main

import (
	"flag"
	"fmt"
	"os"

	"./bin/scm"
)

func main() {
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

	if len(args) == 1 && args[len(args)-1] == "setup" {
		setup(args[0])
	} else if len(args) == 0 {
		parse_revision(".")
	} else {
		parse_revision(args[0])
	}

}

func getScmParser(dir string) (parser scm.ScmParser, err error) {
	err = os.Chdir(dir)

	if err != nil {
		fmt.Errorf("Can't access %s; do you have execute permissions?\n", dir)
		return nil, err
	}

	parser, err = scm.GetParser(dir)

	return
}

func parse_revision(dir string) {

	parser, err := getScmParser(dir)

	if err != nil {
		fmt.Errorf("%s", err)
		return
	}

	if parser != nil {
		scm.ParseAndWrite(parser)
	}
}

func setup(dir string) {
	parser, err := getScmParser(".")

	if err != nil {
		fmt.Errorf("%s", err)
		return
	}

	if parser != nil {
		parser.Setup()
	}
}
