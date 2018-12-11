package icons // import "application/controllers/resource/icons"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

import (
	"fmt"
)

var (
	errNotFound = fmt.Errorf("Not found")
)

// ErrNotFound Not found
func (ici *impl) ErrNotFound() error { return errNotFound }
