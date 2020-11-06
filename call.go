package plog

import (
	"encoding/json"
	"fmt"
)

func (p *Plog) Call(name string, args, rets []interface{}) error {
	<-p.ioReady

	call := p.newCall()
	defer p.releaseCall(call)

	debug("enter call %s_%d", name, call)

	rid, got := p.addRet(name, call)
	defer p.dropRet(rid)

	debug("added hook")

	data, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("marshal %v: %w", args, err)
	}

	m := &msg{
		Name: name,
		Call: call,
		Args: data,
	}

	p.encMu.Lock()
	err = p.enc.Encode(m)
	p.encMu.Unlock()

	debug("encoded %v", m)

	if err != nil {
		debug("enc error %s", err)
		return fmt.Errorf("encode %v: %w", m, err)
	}

	m = <-got
	debug("got m %v", m)

	err = json.Unmarshal(m.Return, &rets)
	if err != nil {
		return fmt.Errorf("unmarshal %q: %w", m.Return, err)
	}

	return nil
}

func (p *Plog) call(m *msg) error {
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
