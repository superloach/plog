package main

import (
	"fmt"
	"os"

	"github.com/superloach/plug"
)

const theString = "this is the string"

var upper func() (string, error)

func getString() string {
	return theString
}

func main() {
	p := plug.Host(os.Args[1], os.Args[2:]...).
		Register("getString", getString).
		Wrap("upper", &upper)

	go p.MustServe()

	upperString, err := upper()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q -> %q\n", theString, upperString)
}
