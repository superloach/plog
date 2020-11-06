package plog

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

// Plog is a plugin connection. Use Host, Guest, or New to make one of these, where appropriate.
type Plog struct {
	dec   *json.Decoder
	decMu sync.Mutex

	enc   *json.Encoder
	encMu sync.Mutex

	fns  map[string]fn
	fnMu sync.Mutex

	rets  map[int]*ret
	retMu sync.Mutex

	calls  map[int]bool
	callMu sync.Mutex

	openFn  func() error
	closeFn func()

	ioReady chan bool
}

// New creates a new Plog connected to in and out. It is preferable to use Host or Guest instead.
func New(in io.Reader, out io.Writer) *Plog {
	p := empty()

	p.dec = json.NewDecoder(in)
	p.enc = json.NewEncoder(out)

	close(p.ioReady)

	return p
}

func empty() *Plog {
	return &Plog{
		fns:   make(map[string]fn),
		rets:  make(map[int]*ret),
		calls: make(map[int]bool),
		openFn: func() error {
			return nil
		},
		closeFn: func() {},
		ioReady: make(chan bool),
	}
}

// Guest makes a Plog which serves on the stdin/stdout of the binary, and can be run with a Host.
func Guest() *Plog {
	return New(os.Stdin, os.Stdout)
}
