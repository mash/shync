package mw

import (
	"github.com/mash/shync"
	"github.com/mash/shync/log"
)

func StatusLog(next fn) fn {
	return func(c shync.Config) error {
		err := next(c)
		if err != nil {
			log.Error("msg", err)
		}
		return nil
	}
}
