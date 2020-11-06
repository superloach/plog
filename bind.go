package plog

import (
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

		retis := make([]interface{}, 0, len(rets))
		for i := 0; i < cap(rets); i++ {
			iface := rets[i].Interface()
			retis = append(retis, &iface)
		}

		err = p.Call(name, argis, retis)
		if err != nil {
			err = fmt.Errorf("call raw %q %v: %w", name, argis, err)
			return rets
		}

		return rets
	})

	f.Set(newf)

	return p
}
