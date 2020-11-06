package plog

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Bind wraps a call to a function on the connected Plog into a function given by reference. Returns the Plog to make chaining easy.
//
//  // note the added error return
//  var fn func(int) (int, error)
//  p.Bind("fn", &fn)
//  sqr2, err := fn(2)
func (p *Plog) Bind(name string, fr interface{}) *Plog {
	debug("wrapping %q", name)

	wrapCallJSON(func(data []byte) ([]byte, error) {
		call := p.newCall()
		defer p.releaseCall(call)

		debug("enter wrapped call json %s_%d", name, call)

		r := p.addRet(name, call)
		defer p.dropRet(r)

		debug("added hook")

		msg := &Msg{
			Name: name,
			Call: call,
			Args: data,
		}

		err := p.Send(msg)
		if err != nil {
			debug("enc error %s", err)
			return nil, fmt.Errorf("encode %v: %w", msg, err)
		}

		debug("sent msg %v", msg)

		msg = <-r.C
		debug("wrap got msg %v", msg)

		return msg.Ret, nil
	}, fr)

	return p
}

func wrapCallJSON(cj func([]byte) ([]byte, error), fr interface{}) {
	f := reflect.ValueOf(fr).Elem()
	t := f.Type()

	newf := reflect.MakeFunc(t, func(args []reflect.Value) []reflect.Value {
		rets := make([]reflect.Value, 0, t.NumOut())
		for i := 0; i < cap(rets); i++ {
			rets = append(rets, reflect.New(t.Out(i)))
		}

		err := error(nil)

		defer func(errp *error) {
			for i, r := range rets {
				rets[i] = r.Elem()
			}
			rets[len(rets)-1] = reflect.ValueOf(errp).Elem()

			debug("%s", rets)
			debug("%q", rets)
			debug("%#v", rets)
		}(&err)

		argis := make([]interface{}, 0, len(args))
		for i := 0; i < cap(argis); i++ {
			argis = append(argis, args[i].Interface())
		}

		argd := []byte(nil)
		argd, err = json.Marshal(argis)
		if err != nil {
			return rets
		}

		debug("calljson %q", argd)

		retd := []byte(nil)
		retd, err = cj(argd)
		if err != nil {
			return rets
		}

		debug("calljson retd %q", retd)

		retis := make([]interface{}, 0, len(rets))
		for i := 0; i < cap(rets); i++ {
			iface := rets[i].Interface()
			retis = append(retis, &iface)
		}

		err = json.Unmarshal(retd, &retis)
		if err != nil {
			return rets
		}

		return rets
	})

	f.Set(newf)
}

func isPrim(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface:
		return false
	case reflect.Array, reflect.Ptr, reflect.Slice:
		return isPrim(t.Elem())
	case reflect.Map:
		return isPrim(t.Key()) && isPrim(t.Elem())
	default:
		panic("unknown prim-ness of type " + t.String())
	}
}
