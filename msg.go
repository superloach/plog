package plog

import (
	"encoding/json"
)

// Msg is the structure used for data passed between Plogs. JSON is to be used for encoding/decoding. The internal structure is subject to change between versions.
type Msg struct {
	Name string          `json:"name"`
	Call int             `json:"call"`
	Args json.RawMessage `json:"args,omitempty"`
	Ret  json.RawMessage `json:"ret,omitempty"`
}

// Messenger is an interface for sending and receiving Msg data. Send and Recv should both block. JSON is the preferred data format.
type Messenger interface {
	Send(*Msg) error
	Recv() (*Msg, error)
	Close() error
}

func (p *Plog) Send(msg *Msg) error {
	p.WaitReady()
	return p.mes.Send(msg)
}

func (p *Plog) Recv() (*Msg, error) {
	p.WaitReady()
	return p.mes.Recv()
}

func (p *Plog) Close() error {
	p.WaitReady()
	return p.mes.Close()
}
