package gzip // import "application/middleware/gzip"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/status"
import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// Gzip Middleware to compress the server response if the client understands compressed content
func Gzip(hndl http.Handler) http.Handler {
	var fn = func(wr http.ResponseWriter, rq *http.Request) {
		var gzw = &impl{
			ResponseWriter: wr,
			isCompressing:  isCompressing(rq.Header),
			Writer:         wr, // by default
		}
		defer gzw.Close()
		hndl.ServeHTTP(gzw, rq)
	}
	return http.HandlerFunc(fn)
}

// Check header
func isCompressing(hdr http.Header) (ret bool) {
	var chk = hdr.Get(header.AcceptEncoding)
	switch {
	case strings.Contains(chk, "gzip"):
		ret = true
		//case strings.Contains(chk, "deflate"):
		//case strings.Contains(chk, "br"):
		//case strings.Contains(chk, "lzma"):
		//case strings.Contains(chk, "sdch"):
	}
	return
}

// WriteHeader Write header code
func (gzp *impl) WriteHeader(code int) {
	var err error

	if gzp.isHeaderWritten {
		return
	}
	if gzp.ResponseWriter.Header().Get(header.ContentEncoding) != "" {
		return
	}
	gzp.isHeaderWritten = true
	defer gzp.ResponseWriter.WriteHeader(code)

	if !gzp.isCompressing {
		return
	}
	gzp.ResponseWriter.Header().Del(header.ContentLength)
	if gzp.Writer, err = gzip.NewWriterLevel(gzp.ResponseWriter, gzip.BestCompression); err != nil {
		gzp.Writer = gzp.ResponseWriter
		return
	}
	gzp.ResponseWriter.Header().Set(header.ContentEncoding, "gzip")

	return
}

// Write Implementation of an interface io.Writer
func (gzp *impl) Write(p []byte) (n int, err error) {
	if !gzp.isHeaderWritten {
		gzp.WriteHeader(status.Ok)
	}
	n, err = gzp.Writer.Write(p)
	return
}

// Flush buffer
func (gzp *impl) Flush() {
	if flusher, ok := gzp.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
	return
}

// Hijack implements the Hijacker.Hijack method.
// Our response is both a ResponseWriter and a Hijacker
func (gzp *impl) Hijack() (conn net.Conn, buf *bufio.ReadWriter, err error) {
	if hijecker, ok := gzp.Writer.(http.Hijacker); ok {
		conn, buf, err = hijecker.Hijack()
	} else {
		err = fmt.Errorf("http.Hijacker is not implemented on this writer")
	}
	return
}

// CloseNotify method
func (gzp *impl) CloseNotify() (ret <-chan bool) {
	if cloceNotifier, ok := gzp.Writer.(http.CloseNotifier); ok {
		ret = cloceNotifier.CloseNotify()
	} else {
		ret = make(chan bool, 1)
	}
	return
}

// Write Implementation of an interface io.WriteCloser
func (gzp *impl) Close() (err error) {
	if writeCloser, ok := gzp.Writer.(io.WriteCloser); ok {
		err = writeCloser.Close()
	} else {
		err = fmt.Errorf("io.WriteCloser is not implemented on this writer")
	}
	return
}

// Push implementation of Push
func (gzp *impl) Push(target string, opts *http.PushOptions) (err error) {
	if pusher, ok := gzp.Writer.(http.Pusher); ok {
		err = pusher.Push(target, opts)
	} else {
		err = fmt.Errorf("http.Pusher is not implemented on this writer")
	}
	return
}
