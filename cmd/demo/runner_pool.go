package main

import (
	"reflect"

	"github.com/omcrgnt/res"
)

type runnerPool struct {
	reg res.Registry
}

func (p runnerPool) Walk(fn func(reflect.Type, any) bool) {
	p.reg.WalkEntries(func(e res.Entry) bool {
		return fn(e.Type, e.Value)
	})
}
