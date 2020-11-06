package plog

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
