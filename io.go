package plog

import (
	"encoding/json"
	"fmt"
	"io"
)

type ioMes struct {
	*json.Decoder
	*json.Encoder
}

// IOMessenger creates a Messenger connected to r and w with an ioMes.
func IOMessenger(r io.Reader, w io.Writer) Messenger {
	return ioMes{
		Decoder: json.NewDecoder(r),
		Encoder: json.NewEncoder(w),
	}
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

type chanMes chan *Msg

// ChanMessenger creates a Messenger connected to c with a chanMes.
func ChanMessenger(c chan *Msg) Messenger {
	return chanMes(c)
}

func (c chanMes) Recv() (*Msg, error) {
	return <-c, nil
}

func (c chanMes) Send(msg *Msg) error {
	c <- msg
	return nil
}
