package plug

import (
	"fmt"
)

func (p *Plug) call(m *msg) error {
	debug("call %v", m)

	p.fnMu.Lock()
	fn, ok := p.fns[m.Name]
	if !ok {
		p.fnMu.Unlock()
		return fmt.Errorf("no fn %q", m.Name)
	}
	p.fnMu.Unlock()

	debug("call %q %q", m.Name, m.Args)

	retd, err := fn.callJSON(m.Args)
	if err != nil {
		return fmt.Errorf("call json %q: %w", m.Args, err)
	}

	debug("%q returned %q", m.Name, retd)

	ret := &msg{
		Name:   m.Name,
		Call:   m.Call,
		Return: retd,
	}

	p.encMu.Lock()
	err = p.enc.Encode(ret)
	p.encMu.Unlock()

	debug("encoded ret %v", ret)

	if err != nil {
		return fmt.Errorf("encode %v: %w", ret, err)
	}

	return nil
}
