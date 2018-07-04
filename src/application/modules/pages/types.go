package pages // import "application/modules/pages"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"os"
	"regexp"
)

const (
	_IndexFile = `index.html`
)

var (
	_RootFilesPattern = regexp.MustCompile(`(?i)^.+\.(htm|html|js|css|eot|svg|woff2|ttf|woff)?$`)
	_ExtractExtension = regexp.MustCompile(`(?i)^.+\.(.+)?$`)
)

// Interface is an interface of package
type Interface interface {
	// Index Выгрузка отпределённых типов файлов из document root
	Index(rfn string) (*File, error)

	// Assets Выгрузка всех файлов из папки assets
	Assets(rfn string) (*File, error)

	// ERRORS

	// ErrFileNotFound File not found
	ErrFileNotFound() error
	// ErrIncorrectFileLength Incorrect file length
	ErrIncorrectFileLength() error
}

// impl is an implementation of package
type impl struct {
	DocumentRoot string // Папка с файлами для выгрузки в web
}

// File Файл со всеми атрибутами
type File struct {
	Body        *bytes.Buffer
	Info        os.FileInfo
	ContentType string
}
