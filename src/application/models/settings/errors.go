package settings // import "application/models/settings"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "fmt"

var (
	errKeyOrValueNotFound = fmt.Errorf("Key or value not found")
	errKeyIsNotUnique     = fmt.Errorf("Key is not unique")
)

// ErrKeyOrValueNotFound Key or value not found
func (st *impl) ErrKeyOrValueNotFound() error { return errKeyOrValueNotFound }

// ErrKeyIsNotUnique Key is not unique
func (st *impl) ErrKeyIsNotUnique() error { return errKeyIsNotUnique }
