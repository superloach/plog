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
	err := p.openFn()
	if err != nil {
		return fmt.Errorf("openfn: %w", err)
	}

	defer p.closeFn()

	for {
		msg, err := p.mes.Recv()
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
