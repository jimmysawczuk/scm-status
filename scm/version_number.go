package scm

import (
	"strconv"
	"strings"
)

type versionNumber string

func (ver versionNumber) split() []int {
	a := strings.Split(strings.TrimSpace(string(ver)), ".")
	b := make([]int, len(a))
	for i, v := range a {
		b[i], _ = strconv.Atoi(v)
	}

	return b
}

func (ver versionNumber) meets(compare versionNumber) bool {
	a := ver.split()
	b := compare.split()

	if len(b) > len(a) {
		a = append(a, make([]int, len(b)-len(a))...)
	} else if len(a) > len(b) {
		b = append(b, make([]int, len(a)-len(b))...)
	}

	for i := 0; i < len(a); i++ {
		if a[i] > b[i] {
			return true
		} else if a[i] < b[i] {
			return false
		}
	}

	return true
}
