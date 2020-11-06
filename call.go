package plog

import (
	"encoding/json"
	"fmt"
)

// Call makes a function call to the connected Plug manually. Using Bind is preferable.
//
//  // note the reference in rets
//  sqr2 := 0
//  err := p.Call("fn",
//  	[]interface{}{2},
//  	[]interface{}{&sqr2},
//  )
func (p *Plog) Call(name string, args, rets []interface{}) error {
	call := p.newCall()
	defer p.releaseCall(call)

	debug("enter call %s_%d", name, call)

	r := p.addRet(name, call)
	defer p.dropRet(r)

	debug("added hook")

	data, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("marshal %v: %w", args, err)
	}

	msg := &Msg{
		Name: name,
		Call: call,
		Args: data,
	}

	err = p.Send(msg)
	if err != nil {
		return fmt.Errorf("send %v: %w", msg, err)
	}

	debug("encoded %v", msg)

	msg = <-r.C
	debug("call got msg %v", msg)

	err = json.Unmarshal(msg.Ret, &rets)
	if err != nil {
		return fmt.Errorf("unmarshal %q: %w", msg.Ret, err)
	}

	return nil
}

func (p *Plog) newCall() int {
	for id, ok := range p.calls {
		if !ok {
			return p.takeCall(id)
		}
	}

	id := len(p.calls)
	return p.takeCall(id)
}

func (p *Plog) takeCall(id int) int {
	p.calls[id] = true
	return id
}

func (p *Plog) releaseCall(id int) {
	p.calls[id] = false
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
