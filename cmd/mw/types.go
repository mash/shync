package mw

import (
	"github.com/mash/shync"
)

type fn func(c shync.Config) error
