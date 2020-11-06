//+build debug

package plog

import (
	"fmt"
	"os"
)

func debug(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "["+os.Args[0]+"] "+f+"\n", args...)
}
