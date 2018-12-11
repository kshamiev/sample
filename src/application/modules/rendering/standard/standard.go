package standard // import "application/modules/rendering/standard"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	textTemplate "text/template"

	"application/modules/rendering/options"
)

// Interface is an interface
type Interface options.Renderrer

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
	var tpl *htmlTemplate.Template
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
	if tpl, err = htmlTemplate.ParseFiles(templates...); err != nil {
		return
	}
	if err = tpl.Execute(wr, values); err != nil {
		return
	}

	return
}

// RenderText Парсинг множества шаблонов файлов с указанными переменными
func (t *impl) RenderText(wr io.Writer, values interface{}, tpls ...string) (err error) {
	var tpl *textTemplate.Template
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
	if tpl, err = textTemplate.ParseFiles(templates...); err != nil {
		return
	}
	if err = tpl.Execute(wr, values); err != nil {
		return
	}

	return
}

// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов *bytes.Buffer
func (t *impl) RenderHTMLData(wr io.Writer, values interface{}, buffers ...io.Reader) (err error) {
	var tpl *htmlTemplate.Template
	var tmpl *htmlTemplate.Template
	var buf []byte
	var i int

	// Шаблонизатор паникует, сука...
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	for i = len(buffers) - 1; i >= 0; i-- {
		if tpl == nil {
			tpl = htmlTemplate.New(fmt.Sprintf("%d", i))
			tmpl = tpl
		} else {
			tmpl = tpl.New(fmt.Sprintf("%d", i))
		}
		if buf, err = ioutil.ReadAll(buffers[i]); err != nil {
			return
		}
		if _, err = tmpl.Parse(string(buf)); err != nil {
			return
		}
	}
	if err = tpl.Execute(wr, values); err != nil {
		return
	}

	return
}

// RenderTextData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов *bytes.Buffer
func (t *impl) RenderTextData(wr io.Writer, values interface{}, buffers ...io.Reader) (err error) {
	var tpl *textTemplate.Template
	var tmpl *textTemplate.Template
	var buf []byte
	var i int

	// Шаблонизатор паникует, сука...
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	for i = len(buffers) - 1; i >= 0; i-- {
		if tpl == nil {
			tpl = textTemplate.New(fmt.Sprintf("%d", i))
			tmpl = tpl
		} else {
			tmpl = tpl.New(fmt.Sprintf("%d", i))
		}
		if buf, err = ioutil.ReadAll(buffers[i]); err != nil {
			return
		}
		if _, err = tmpl.Parse(string(buf)); err != nil {
			return
		}
	}
	if err = tpl.Execute(wr, values); err != nil {
		return
	}

	return
}
