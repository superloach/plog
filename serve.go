package plog

import (
	"errors"
	"fmt"
	"io"
)

// MustServe simply calls Serve and panics if there's an error.
func (p *Plog) MustServe() {
	err := p.Serve()
	if err != nil {
		panic(fmt.Errorf("serve: %w", err))
	}
}

// Serve runs the Plog's event loop and returns if an error occurs.
func (p *Plog) Serve() error {
	mes, err := p.opener()
	if err != nil {
		return fmt.Errorf("openfn: %w", err)
	}

	p.mes = mes
	close(p.ready)

	defer p.Close()

	for {
		msg, err := p.Recv()
		debug("recv msg %v", msg)

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}

		if msg.Ret != nil {
			debug("is a ret")

			p.retMu.Lock()
			for id, r := range p.rets {
				if r == nil {
					continue
				}

				debug("spawn ret %d", id)

				go r.Ret(msg)
			}
			p.retMu.Unlock()

			continue
		}

		if msg.Args != nil {
			debug("is a call")

			go func() {
				err := p.localCall(msg)
				if err != nil {
					debug("call %v: ERROR %s", msg, err)
				}
			}()

			continue
		}

		debug("what")
	}
}

func (p *Plog) localCall(msg *Msg) error {
	debug("call %v", msg)

	p.fnMu.Lock()
	fn, ok := p.fns[msg.Name]
	if !ok {
		p.fnMu.Unlock()
		return fmt.Errorf("no fn %q", msg.Name)
	}
	p.fnMu.Unlock()

	debug("call %q %q", msg.Name, msg.Args)

	retd, err := fn.callJSON(msg.Args)
	if err != nil {
		return fmt.Errorf("call json %q: %w", msg.Args, err)
	}

	debug("%q returned %q", msg.Name, retd)

	msg = &Msg{
		Name: msg.Name,
		Call: msg.Call,
		Ret:  retd,
	}

	err = p.Send(msg)
	if err != nil {
		return fmt.Errorf("send %v: %w", msg, err)
	}

	return nil
}
