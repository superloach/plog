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
		panic(err)
	}
}

// Serve runs the Plog's event loop and returns if an error occurs.
func (p *Plog) Serve() error {
	err := p.openFn()
	if err != nil {
		return fmt.Errorf("openfn: %w", err)
	}

	defer p.closeFn()

	errs := make(chan error)

	for {
		select {
		case err := <-errs:
			return err
		default:
		}

		m := &msg{}

		p.decMu.Lock()
		err := p.dec.Decode(&m)
		p.decMu.Unlock()

		debug("decoded m %v", m)

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}

		if m.Return != nil {
			debug("is a ret")

			p.retMu.Lock()
			for id, r := range p.rets {
				if r == nil {
					continue
				}

				debug("spawn ret %d", id)

				go r.Ret(m)
			}
			p.retMu.Unlock()

			continue
		}

		if m.Args != nil {
			debug("is a call")

			go func() {
				err := p.call(m)
				if err != nil {
					errs <- err
				}
			}()

			continue
		}

		debug("what")
	}
}
