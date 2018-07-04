package standard // import "application/modules/rendering/standard"

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"application/modules/rendering/options"
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
	var o = new(impl)
	for _, arg := range args {
		o.Option = arg
	}
	return o
}

// RenderHTML Парсинг множества шаблонов файлов с указанными переменными
func (t *impl) RenderHTML(wr io.Writer, values interface{}, tpls ...string) (err error) {
	var tpl *template.Template
	var i int
	var templates []string

	// Шаблонизатор паникует, сука...
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	for i = range tpls {
		templates = append(templates, filepath.Join(t.Option.Directory, tpls[i]))
	}
	if tpl, err = template.ParseFiles(templates...); err != nil {
		return
	}
	if err = tpl.Execute(wr, values); err != nil {
		return
	}

	return
}

// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов *bytes.Buffer
func (r *impl) RenderHTMLData(wr io.Writer, values interface{}, buffers ...*bytes.Buffer) (err error) {
	var tpl *template.Template
	var tmpl *template.Template
	var i int

	// Шаблонизатор паникует, сука...
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	for i = len(buffers) - 1; i >= 0; i-- {
		if tpl == nil {
			tpl = template.New(fmt.Sprintf("%d", i))
			tmpl = tpl
		} else {
			tmpl = tpl.New(fmt.Sprintf("%d", i))
		}
		if _, err = tmpl.Parse(buffers[i].String()); err != nil {
			return
		}
	}
	if err = tpl.Execute(wr, values); err != nil {
		return
	}

	return
}
