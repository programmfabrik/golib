package golib

import (
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
			defer func() {
				r2 := recover()
				if r2 != nil {
					logrus.Errorf("Panic inside recover: %v", r2)
					logrus.Debugf(string(debug.Stack()))
				}
			}()
			err2, ok := r.(error)
			if ok {
				f(err2)
			} else {
				f(err)
			}
		}()
	}
}
