package gzip // import "application/middleware/gzip"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"net/http"
)

// impl is an implementation of package
type impl struct {
	http.ResponseWriter
	isCompressing   bool
	isHeaderWritten bool
	Writer          io.Writer
}
