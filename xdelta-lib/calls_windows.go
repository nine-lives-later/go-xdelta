package lib

import (
	"syscall"
)

var (
	xdeltaDLL = syscall.NewLazyDLL(`.\go-xdelta-lib.dll`)
)
