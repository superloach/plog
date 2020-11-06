package main

import (
	"strings"

	"github.com/superloach/plog"
	"github.com/superloach/plog/_example/common"
)

var getString func() (string, error)

func foo() common.Struct {
	return common.Struct{
		Str:    "owo",
		Int:    123,
		Uint32: 0b1001011001101001,
		Bools: []bool{
			true,
			false,
			false,
			true,
		},
	}
}

func upper() string {
	s, err := getString()
	if err != nil {
		panic(err)
	}

	return strings.ToUpper(s)
}

func main() {
	plog.New(plog.StdIO()).
		Expose("upper", upper).
		Expose("foo", foo).
		Bind("getString", &getString).
		MustServe()
}
