package plug

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Wrap wraps a call to a function on the connected Plug into a function value. Returns the Plug to make chaining easy.
func (p *Plug) Wrap(name string, fr interface{}) *Plug {
	debug("wrapping %q", name)

	wrapCallJSON(func(data []byte) ([]byte, error) {
		<-p.ioReady

		id := p.newCall()
		defer p.releaseCall(id)

		callName := fmt.Sprintf("%s_%d", name, id)

		debug("enter wrapped call json %q", callName)

		rid, got := p.addRet(callName)
		defer p.dropRet(rid)

		debug("added hook")

		m := &msg{
			Name: callName,
			Call: data,
		}

		p.encMu.Lock()
		debug("enc %#v", p.enc)
		err := p.enc.Encode(m)
		p.encMu.Unlock()

		debug("encoded m")

		if err != nil {
			debug("enc error %s", err)
			return nil, fmt.Errorf("encode %v: %w", m, err)
		}

		m = <-got
		debug("got m %v", m)

		return m.Return, nil
	}, fr)

	return p
}

func wrapCallJSON(cj func([]byte) ([]byte, error), fr interface{}) {
	f := reflect.ValueOf(fr).Elem()
	t := f.Type()

	ins := make([]reflect.Type, 0, t.NumIn())
	for i := 0; i < cap(ins); i++ {
		ins = append(ins, t.In(i))
	}

	outs := make([]reflect.Type, 0, t.NumOut())
	for i := 0; i < cap(outs); i++ {
		outs = append(outs, t.Out(i))
	}

	nt := reflect.FuncOf(ins, outs, t.IsVariadic())

	newf := reflect.MakeFunc(nt, func(args []reflect.Value) []reflect.Value {
		rets := make([]reflect.Value, 0, t.NumOut())
		for i := 0; i < cap(rets); i++ {
			rets = append(rets, reflect.New(t.Out(i)))
		}

		defer func() {
			for i, ret := range rets {
				rets[i] = ret.Elem()
			}
		}()

		argis := make([]interface{}, 0, len(args))
		for i := 0; i < cap(argis); i++ {
			argis = append(argis, args[i].Interface())
		}

		argd, err := json.Marshal(argis)
		if err != nil {
			rets[len(rets)-1] = reflect.ValueOf(err)
			return rets
		}

		debug("calljson %q", argd)

		retd, err := cj(argd)
		if err != nil {
			rets[len(rets)-1] = reflect.ValueOf(err)
			return rets
		}

		debug("calljson retd %q", retd)

		retis := make([]interface{}, 0, len(rets))
		for i := 0; i < cap(rets); i++ {
			iface := rets[i].Interface()
			retis = append(retis, iface)
		}

		err = json.Unmarshal(retd, &retis)
		if err != nil {
			rets[len(rets)-1] = reflect.ValueOf(err)
			return rets
		}

		return rets
	})

	f.Set(newf)
}
