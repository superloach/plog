package plog

type ret struct {
	ID   int
	Name string
	Call int
	C    chan *Msg
}

func (p *Plog) addRet(name string, call int) *ret {
	p.retMu.Lock()
	defer p.retMu.Unlock()

	r := &ret{
		Name: name,
		Call: call,
		C:    make(chan *Msg),
	}

	for id, er := range p.rets {
		if er == nil {
			return p.putRet(id, r)
		}
	}

	id := len(p.rets)
	return p.putRet(id, r)
}

func (p *Plog) putRet(id int, r *ret) *ret {
	r.ID = id
	p.rets[id] = r
	return r
}

func (p *Plog) dropRet(r *ret) {
	close(r.C)
	p.rets[r.ID] = nil
}

func (r *ret) Ret(msg *Msg) {
	if msg.Name != r.Name {
		return
	}

	if msg.Call != r.Call {
		return
	}

	r.C <- msg
}
