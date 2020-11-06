package main

import (
	"strings"

	"github.com/superloach/plug"
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
	plug.Guest().
		Register("upper", upper).
		Wrap("getString", &getString).
		MustServe()
}
