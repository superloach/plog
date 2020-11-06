package plug

type ret struct {
	Name string
	C    chan *msg
}

func (p *Plug) addRet(name string) (int, chan *msg) {
	p.retMu.Lock()
	defer p.retMu.Unlock()

	r := &ret{
		Name: name,
		C:    make(chan *msg),
	}

	for id, er := range p.rets {
		if er == nil {
			return id, p.putRet(id, r)
		}
	}

	id := len(p.rets)
	return id, p.putRet(id, r)
}

func (p *Plug) putRet(id int, r *ret) chan *msg {
	p.rets[id] = r
	return r.C
}

func (p *Plug) dropRet(id int) {
	p.rets[id] = nil
}

func (r *ret) Ret(m *msg) {
	if m.Name != r.Name {
		return
	}

	r.C <- m
}
