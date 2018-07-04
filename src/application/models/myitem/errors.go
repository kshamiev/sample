package myitem

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
)

var (
	errNotFound = fmt.Errorf("Not found")
)

// ErrNotFound Not found
func (mim *impl) ErrNotFound() error { return errNotFound }
