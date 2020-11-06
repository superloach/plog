package plog

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

// Plog is a plugin connection. Use Client, Server, IO, or New to make one of these, where appropriate.
type Plog struct {
	mes      Messenger
	mesReady chan bool

	fns  map[string]fn
	fnMu sync.Mutex

	rets  map[int]*ret
	retMu sync.Mutex

	calls  map[int]bool
	callMu sync.Mutex

	openFn  func() error
	closeFn func()
}

func empty() *Plog {
	return &Plog{
		mesReady: make(chan bool),
		fns:      make(map[string]fn),
		rets:     make(map[int]*ret),
		calls:    make(map[int]bool),
		openFn: func() error {
			return nil
		},
		closeFn: func() {},
	}
}

// New creates a new Plog with the given Messenger.
func New(mes Messenger) *Plog {
	p := empty()

	p.openFn = func() error {
		p.mes = mes
		close(p.mesReady)
		return nil
	}

	return p
}

// IO creates a new Plog connected to in and out, using an ioMessenger.
func IO(in io.Reader, out io.Writer) *Plog {
	return New(ioMessenger{
		Decoder: json.NewDecoder(in),
		Encoder: json.NewEncoder(out),
	})
}

// StdIO makes a Plog which serves on the stdin/stdout of the binary, and can be run with an Exec.
func StdIO() *Plog {
	return IO(os.Stdin, os.Stdout)
}
