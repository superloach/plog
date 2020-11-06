package plog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// StdIO returns an Opener which connects to stdin/stdout, to be connected to by Exec.
func StdIO() Opener {
	return IO(os.Stdin, os.Stdout)
}

// IO returns an Opener which connects to r and w with an ioMes.
func IO(r io.Reader, w io.Writer) Opener {
	return func() (Messenger, error) {
		return ioMes{
			Decoder: json.NewDecoder(r),
			Encoder: json.NewEncoder(w),
		}, nil
	}
}

type ioMes struct {
	*json.Decoder
	*json.Encoder
}

func (i ioMes) Recv() (*Msg, error) {
	msg := &Msg{}

	err := i.Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return msg, nil
}

func (i ioMes) Send(msg *Msg) error {
	err := i.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode %v: %w", msg, err)
	}

	return nil
}

func (i ioMes) Close() error {
	return nil
}
