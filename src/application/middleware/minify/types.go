package minify // import "application/middleware/minify"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"net/http"

	"github.com/tdewolff/minify"
)

// WriteCloser Minify write closer interface
type WriteCloser interface {
	io.WriteCloser
}

// minifyResponseWriter wraps an http.ResponseWriter and makes sure that errors from the minifier are passed down through Close (can be blocking).
// All writes to the response writer are intercepted and minified on the fly.
// http.ResponseWriter loses all functionality such as Pusher, Hijacker, Flusher, ...
type minifyResponseWriter struct {
	http.ResponseWriter
	writer    WriteCloser
	m         *minify.M
	mediatype string
}
