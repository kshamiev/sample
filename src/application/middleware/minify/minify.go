package minify // import "application/middleware/minify"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/header"
import (
	"net/http"
	"regexp"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

// Minify Middleware that implements minification on-the-fly for CSS, HTML, JSON, SVG and XML
func Minify(hndl http.Handler) http.Handler {
	var mnf = minify.New()
	mnf.AddFunc("text/css", css.Minify)
	mnf.AddFunc("text/html", html.Minify)
	mnf.AddFunc("text/javascript", js.Minify)
	mnf.AddFunc("image/svg+xml", svg.Minify)
	mnf.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	mnf.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	var fn = func(wr http.ResponseWriter, rq *http.Request) {
		var mware = &minifyResponseWriter{wr, nil, mnf, ""}
		defer mware.Close()
		hndl.ServeHTTP(mware, rq)
	}

	return http.HandlerFunc(fn)
}

// WriteHeader intercepts any header writes and removes the Content-Length header.
func (w *minifyResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

// Write intercepts any writes to the response writer.
// The first write will extract the Content-Type as the mediatype. Otherwise it falls back to the RequestURI extension.
func (w *minifyResponseWriter) Write(b []byte) (int, error) {
	var parts []string
	if w.writer == nil {
		parts = strings.Split(w.ResponseWriter.Header().Get(header.ContentType), ";")
		if len(parts) > 0 && parts[0] != "" {
			w.mediatype = parts[0]
		}
		w.Header().Set(header.ContentType, w.mediatype)
		w.writer = w.m.Writer(w.mediatype, w.ResponseWriter)
	}
	return w.writer.Write(b)
}

// Close must be called when writing has finished. It returns the error from the minifier.
func (w *minifyResponseWriter) Close() (err error) {
	if w.writer != nil {
		err = w.writer.Close()
	}
	return
}
