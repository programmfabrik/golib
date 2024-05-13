package golib

import (
	"fmt"
	"runtime/debug"
)

// Recover can be used in defer to execute "f" on recover
func Recover(f func(err error)) {
	r := recover()
	if r == nil {
		return
	}
	err := fmt.Errorf("Panic: %v", r)
	Pln(err.Error())
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
