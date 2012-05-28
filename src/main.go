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
	flag.String("out", "REVISION", "File in which to (over)write the revision info into")
	flag.String("executable", "/usr/local/bin/scm-status", "Path at which this program can be executed (for hooks)")
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

	scm, err := getScmParser(dir)

	if err != nil {
		fmt.Errorf("%s", err)
		return
	}

	if scm != nil {
		result := scm.Parse()

		json, _ := result.ToJSON()

		filename := flag.Lookup("out").Value.String()

		fp, err := os.OpenFile(scm.Dir()+"/"+filename, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)

		if err == nil {
			fp.Write(json)

			fp.Close()

		} else {

			fmt.Errorf("%s\n", err)

		}
	}
}

func setup(dir string) {
	scm, err := getScmParser(".")

	if err != nil {
		fmt.Errorf("%s", err)
		return
	}

	if scm != nil {
		scm.Setup()
	}
}
