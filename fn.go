package plog

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type fn struct {
	f reflect.Value
	t reflect.Type
}

func fnOf(f interface{}) fn {
	v := reflect.ValueOf(f)
	return fn{v, v.Type()}
}

// Register adds the given function to the Plog, to be called by the Plog on the other end. It returns the Plog to make chaining easy.
func (p *Plog) Register(name string, f interface{}) *Plog {
	p.fnMu.Lock()
	p.fns[name] = fnOf(f)
	p.fnMu.Unlock()

	return p
}

func (f fn) callJSON(data []byte) ([]byte, error) {
	argis := make([]interface{}, 0, f.t.NumIn())
	for i := 0; i < cap(argis); i++ {
		argis = append(argis, reflect.New(f.t.In(i)).Interface())
	}

	err := json.Unmarshal(data, &argis)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %q %v: %w", data, argis, err)
	}

	for i := 0; i < len(argis); i++ {
		argis[i] = reflect.ValueOf(argis[i]).Elem().Interface()
	}

	return json.Marshal(f.call(argis...))
}

func (f fn) call(argis ...interface{}) []interface{} {
	args := make([]reflect.Value, 0, len(argis))
	for i := 0; i < cap(args); i++ {
		args = append(args, reflect.ValueOf(argis[i]))
	}

	retis := make([]interface{}, 0, f.t.NumOut())
	for i := 0; i < cap(retis); i++ {
		retis = append(retis, reflect.Zero(f.t.Out(i)).Interface())
	}

	rets := f.f.Call(args)

	for i := 0; i < len(rets); i++ {
		reflect.ValueOf(&retis[i]).Elem().Set(rets[i])
	}

	return retis
}
