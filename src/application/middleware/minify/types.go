package minify // import "application/middleware/minify"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"

	"github.com/tdewolff/minify"
)

type MinifyWriter interface {
	Write(b []byte) (int, error)
	Close() error
}

// minifyResponseWriter wraps an http.ResponseWriter and makes sure that errors from the minifier are passed down through Close (can be blocking).
// All writes to the response writer are intercepted and minified on the fly.
// http.ResponseWriter loses all functionality such as Pusher, Hijacker, Flusher, ...
type minifyResponseWriter struct {
	http.ResponseWriter
	writer    MinifyWriter
	m         *minify.M
	mediatype string
}
