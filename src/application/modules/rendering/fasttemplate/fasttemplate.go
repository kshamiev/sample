package fasttemplate // import "application/modules/rendering/fasttemplate"

import (
	"bytes"
	"fmt"
	"io"

	"application/modules/rendering/options"
	// "github.com/valyala/fasttemplate"
)

// Interface is an interface of repository
type Interface interface {
	// RenderHTML Парсинг множества шаблонов файлов с указанными переменными
	RenderHTML(io.Writer, interface{}, ...string) error

	// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов *bytes.Buffer
	RenderHTMLData(io.Writer, interface{}, ...*bytes.Buffer) error
}

// impl is an implementation of repository
type impl struct {
	Option options.Option
}

// New creates new implementation render
func New(args ...options.Option) Interface {
	var ft = new(impl)
	for _, arg := range args {
		ft.Option = arg
	}
	return ft
}

// RenderHTML Парсинг множества шаблонов файлов с указанными переменными
func (ft *impl) RenderHTML(wr io.Writer, values interface{}, tpls ...string) (err error) {
	err = fmt.Errorf("fasttemplate is not implemented :(")
	return
}

// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов *bytes.Buffer
func (ft *impl) RenderHTMLData(wr io.Writer, values interface{}, buffers ...*bytes.Buffer) (err error) {
	err = fmt.Errorf("fasttemplate is not implemented :(")
	return
}
