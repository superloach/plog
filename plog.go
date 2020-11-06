package plog

import (
	"sync"
)

type Opener = func() (Messenger, error)

// Plog is a plugin connection. Use StdIO, Exec, New, or Will to make one of these, where appropriate.
type Plog struct {
	opener Opener
	mes    Messenger
	ready  chan bool

	fns  map[string]fn
	fnMu sync.Mutex

	rets  map[int]*ret
	retMu sync.Mutex

	calls  map[int]bool
	callMu sync.Mutex
}

func (p *Plog) WaitReady() {
	<-p.ready
}

// New creates a new Plog with the given Opener.
func New(o Opener) *Plog {
	return &Plog{
		opener: o,
		ready:  make(chan bool),
		fns:    make(map[string]fn),
		rets:   make(map[int]*ret),
		calls:  make(map[int]bool),
	}
}
