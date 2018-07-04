package pages // import "application/modules/pages"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
)

var (
	errFileNotFound        = fmt.Errorf("File not found")
	errIncorrectFileLength = fmt.Errorf("Incorrect file length")
)

// ErrFileNotFound File not found
func (pgs *impl) ErrFileNotFound() error { return errFileNotFound }

// ErrIncorrectFileLength Incorrect file length
func (pgs *impl) ErrIncorrectFileLength() error { return errIncorrectFileLength }
