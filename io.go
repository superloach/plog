package plog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// StdIO returns an Opener which connects to stdin/stdout, to be connected to by Exec.
func StdIO() Opener {
	return IO(os.Stdin, os.Stdout)
}

// IO returns an Opener which connects to r and w with an ioMes.
func IO(r io.Reader, w io.Writer) Opener {
	return func() (Messenger, error) {
		return ioMes{
			RClose: closer(r),
			Dec:    json.NewDecoder(r),

			WClose: closer(w),
			Enc:    json.NewEncoder(w),
		}, nil
	}
}

type ioMes struct {
	RClose func() error
	Dec    *json.Decoder
	DecMu  sync.Mutex

	WClose func() error
	Enc    *json.Encoder
	EncMu  sync.Mutex
}

func (i ioMes) Recv() (*Msg, error) {
	i.DecMu.Lock()
	defer i.DecMu.Unlock()

	msg := &Msg{}

	err := i.Dec.Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return msg, nil
}

func (i ioMes) Send(msg *Msg) error {
	i.EncMu.Lock()
	defer i.EncMu.Unlock()

	err := i.Enc.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode %v: %w", msg, err)
	}

	return nil
}

func (i ioMes) Close() error {
	err := i.RClose()
	if err != nil {
		return fmt.Errorf("close r: %w", err)
	}

	err = i.WClose()
	if err != nil {
		return fmt.Errorf("close w: %w", err)
	}

	return nil
}

func closer(v interface{}) func() error {
	c, ok := v.(io.Closer)
	if ok {
		return c.Close
	}

	return func() error {
		return nil
	}
}
