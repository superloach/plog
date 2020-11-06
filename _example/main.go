package main

import (
	"fmt"
	"os"

	"github.com/superloach/plog"
	"github.com/superloach/plog/_example/common"
)

const theString = "this is the string"

type Str string

var (
	upper func() (Str, error)
	foo   func() (common.Struct, error)
)

func getString() Str {
	return theString
}

func main() {
	p := plog.Exec(os.Args[1], os.Args[2:]...).
		Expose("getString", getString).
		Bind("upper", &upper).
		Bind("foo", &foo)

	go p.MustServe()

	upperString, err := upper()
	if err != nil {
		panic(err)
	}

	fmt.Printf("(wrap) %q -> %q\n", theString, upperString)

	err = p.Call("upper",
		[]interface{}{},
		[]interface{}{&upperString},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("(call) %q -> %q\n", theString, upperString)

	s, err := foo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("(foo) %v\n", s)
}
