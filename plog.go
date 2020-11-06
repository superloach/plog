package plog

import (
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

// StdIO makes a Plog which serves on the stdin/stdout of the binary, and can be connected to by Exec.
func StdIO() *Plog {
	return New(IOMessenger(os.Stdin, os.Stdout))
}
