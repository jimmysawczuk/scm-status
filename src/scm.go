package scm

import (
	"fmt"
)

type ScmManager interface {
	Parse() map[string]interface{}
}

func Write(scm ScmManager) {
	fmt.Println("got here")
}
