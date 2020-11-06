package plug

import (
	"fmt"
	"strings"
)

func (p *Plug) call(m *msg, errs chan error) {
	idx := strings.LastIndexByte(m.Name, '_')
	if idx < 0 {
		errs <- fmt.Errorf("idx %d < 0", idx)
		return
	}

	name := m.Name[:idx]

	debug("name is %q", name)

	p.fnMu.Lock()
	fn, ok := p.fns[name]
	if !ok {
		p.fnMu.Unlock()
		errs <- fmt.Errorf("no fn %q", name)
		return
	}
	p.fnMu.Unlock()

	debug("call %q %q", name, m.Call)

	retd, err := fn.callJSON(m.Call)
	if err != nil {
		errs <- fmt.Errorf("call json %q: %w", m.Call, err)
		return
	}

	debug("%q returned %q", name, retd)

	ret := &msg{
		Name:   m.Name,
		Return: retd,
	}

	p.encMu.Lock()
	err = p.enc.Encode(ret)
	p.encMu.Unlock()

	debug("encoded ret %v", ret)

	if err != nil {
		errs <- fmt.Errorf("encode %v: %w", ret, err)
		return
	}
}
