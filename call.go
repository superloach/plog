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

	debug("got msg %v", msg)

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
