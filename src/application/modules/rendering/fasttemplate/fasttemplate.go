package fasttemplate // import "application/modules/rendering/fasttemplate"

import (
	"fmt"
	"io"

	"application/modules/rendering/options"
	// "github.com/valyala/fasttemplate"
)

// Interface is an interface
type Interface options.Renderrer

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

// RenderText Парсинг множества шаблонов файлов с указанными переменными
func (ft *impl) RenderText(wr io.Writer, values interface{}, tpls ...string) (err error) {
	err = fmt.Errorf("fasttemplate is not implemented :(")
	return
}

// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов io.Reader
func (ft *impl) RenderHTMLData(wr io.Writer, values interface{}, buffers ...io.Reader) (err error) {
	err = fmt.Errorf("fasttemplate is not implemented :(")
	return
}

// RenderTextData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов io.Reader
func (ft *impl) RenderTextData(wr io.Writer, values interface{}, buffers ...io.Reader) (err error) {
	err = fmt.Errorf("fasttemplate is not implemented :(")
	return
}
