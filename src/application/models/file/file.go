package file // import "application/models/file"

//go:generate go run mime_generate.go
//go:generate go fmt mime.go

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	stdMime "mime"
)

func init() {
	// Добавление всех mime types
	mimeAddAll()
}

// New creates new object and return Interface
func New() Interface {
	var obj = new(impl)
	return obj
}

// Добавление всех mime types
func mimeAddAll() {
	const preffix = `.`
	var mt string

	for mt = range mimeTypeExtension {
		_ = stdMime.AddExtensionType(preffix+mimeTypeExtension[mt], mt) // nolint: gosec
	}
}
