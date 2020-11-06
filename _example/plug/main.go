package main

import (
	"strings"

	"github.com/superloach/plog"
)

var getString func() (string, error)

func upper() string {
	s, err := getString()
	if err != nil {
		panic(err)
	}

	return strings.ToUpper(s)
}

func main() {
	plog.Guest().
		Register("upper", upper).
		Wrap("getString", &getString).
		MustServe()
}
