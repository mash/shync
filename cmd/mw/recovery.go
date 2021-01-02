package mw

import (
	"runtime/debug"

	"github.com/mash/shync"
	"github.com/mash/shync/log"
)

func Recover(next fn) fn {
	return func(c shync.Config) error {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				log.Error("PANIC", err, "stack", string(stack))
			}
		}()
		return next(c)
	}
}
