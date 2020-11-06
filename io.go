package plog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// IO creates a new Plog connected to in and out, using IOMessenger.
func IO(in io.Reader, out io.Writer) *Plog {
	return New(IOMessenger(in, out))
}

// StdIO makes a Plog which serves on the stdin/stdout of the binary, and can be run with an Exec.
func StdIO() *Plog {
	return IO(os.Stdin, os.Stdout)
}

type ioMessenger struct {
	*json.Decoder
	*json.Encoder
}

// IOMessenger creates a Messenger connected to r and w.
func IOMessenger(r io.Reader, w io.Writer) Messenger {
	return ioMessenger{
		Decoder: json.NewDecoder(r),
		Encoder: json.NewEncoder(w),
	}
}

func (i ioMessenger) Recv() (*Msg, error) {
	msg := &Msg{}

	err := i.Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return msg, nil
}

func (i ioMessenger) Send(msg *Msg) error {
	err := i.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode %v: %w", msg, err)
	}

	return nil
}
