package plog

import (
	"encoding/json"
	"fmt"
)

// Msg is the structure used for data passed between Plogs. The internal structure is subject to change between minor versions.
type Msg struct {
	Name string          `json:"name"`
	Call int             `json:"call"`
	Args json.RawMessage `json:"args,omitempty"`
	Ret  json.RawMessage `json:"ret,omitempty"`
}

// Messenger is an interface for sending and receiving Msg data. Send and Recv should both block.
type Messenger interface {
	Send(*Msg) error
	Recv() (*Msg, error)
}

type ioMessenger struct {
	*json.Decoder
	*json.Encoder
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
