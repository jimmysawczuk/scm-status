package main

import (
	"bytes"
	"flag"
	"fmt"

	"./bin/scm"
)

func main() {
	handle()
}

func handle() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 1 && args[len(args)-1] == "setup" {
		fmt.Println("running setup")
	} else {
		// fmt.Printf("running normally on %v\n", args)
		parse_revision(args[0])

	}

}

func parse_revision(dir string) {

	scm, err := scm.GetParser(dir)

	if err != nil {
		fmt.Println(err)
	}

	if scm != nil {
		result := scm.Parse()

		json_bytes, _ := result.ToJSON()
		json := bytes.NewBuffer(json_bytes).String()

		fmt.Println(json)

		// now we write this to a file!
	}
}
