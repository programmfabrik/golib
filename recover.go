package golib

import (
	"runtime/debug"

	"github.com/pkg/errors"
)

// Recover can be used in defer to execute "f" on recover
func Recover(f func(err error)) {
	r := recover()
	if r == nil {
		return
	}
	err := errors.Errorf("Panic: %v", r)
	println(err.Error())
	debug.PrintStack()
	if f != nil {
		func() {
			// if f throws a panic we will not catch it here so the caller knows
			err2, ok := r.(error)
			if ok {
				f(err2)
			} else {
				f(err)
			}
		}()
	}
}
