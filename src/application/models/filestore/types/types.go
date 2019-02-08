package types // import "application/models/filestore/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
)

// File Interface of file
type File interface {
	io.Reader
	io.ReaderAt
	io.Closer
	io.Seeker
	io.Writer
	io.WriterAt
}
