package main

import (
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

	fmt.Println(args)

	if len(args) == 1 && args[len(args)-1] == "setup" {
		fmt.Println("running setup")
	} else {
		fmt.Printf("running normally on %v\n", args)
	}

}
