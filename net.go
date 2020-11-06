package plog

import (
	"fmt"
	"net"
)

// Net returns an Opener that uses net.Dial and IO.
func Net(network, address string) Opener {
	return func() (Messenger, error) {
		conn, err := net.Dial(network, address)
		if err != nil {
			return nil, fmt.Errorf("dial %s/%s: %w", network, address, err)
		}

		io, err := IO(conn, conn)()
		if err != nil {
			return nil, fmt.Errorf("io: %w", err)
		}

		return netMes{
			Conn: conn,
			IO:   io,
		}, nil
	}
}

type netMes struct {
	Conn net.Conn
	IO   Messenger
}

func (n netMes) Recv() (*Msg, error) {
	return n.IO.Recv()
}

func (n netMes) Send(msg *Msg) error {
	return n.IO.Send(msg)
}

func (n netMes) Close() error {
	err := n.Conn.Close()
	if err != nil {
		return fmt.Errorf("conn close: %w", err)
	}

	return n.IO.Close()
}
