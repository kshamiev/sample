package pongo2 // import "application/modules/rendering/pongo2"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"io/ioutil"
	"path/filepath"
	"sync"

	"application/modules/rendering/options"

	"github.com/flosch/pongo2"
)

// Interface is an interface
type Interface options.Renderrer

// impl is an implementation of repository
type impl struct {
	Option    options.Option
	templates map[string]*pongo2.Template
	lock      sync.RWMutex
}

// New creates new implementation render
func New(args ...options.Option) Interface {
	var o = new(impl)
	for _, arg := range args {
		o.Option = arg
	}
	o.templates = make(map[string]*pongo2.Template)
	return o
}

// getTemplate Load template by name
func (t *impl) getTemplate(name string) (tpl *pongo2.Template, err error) {
	var ok bool

	if t.Option.Reload {
		return pongo2.FromFile(filepath.Join(t.Option.Directory, name))
	}
	t.lock.RLock()
	defer t.lock.RUnlock()

	if tpl, ok = t.templates[name]; !ok {
		tpl, err = t.buildTemplatesCache(name)
	}

	return
}

// getContext Template context
func (t *impl) getContext(templateData interface{}) pongo2.Context {
	if templateData == nil {
		return nil
	}
	contextData, isMap := templateData.(map[string]interface{})
	if isMap {
		return contextData
	}
	return nil
}

// buildTemplatesCache Cache template in memory map
func (t *impl) buildTemplatesCache(name string) (tpl *pongo2.Template, err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	tpl, err = pongo2.FromFile(filepath.Join(t.Option.Directory, name))
	if err != nil {
		return
	}
	t.templates[name] = tpl

	return
}

// RenderHTML Парсинг множества шаблонов файлов с указанными переменными
func (t *impl) RenderHTML(wr io.Writer, values interface{}, tpls ...string) (err error) {
	var template *pongo2.Template
	var tpl string
	for _, tpl = range tpls {
		if template, err = t.getTemplate(tpl); err != nil {
			return
		}
		if err = template.ExecuteWriter(t.getContext(values), wr); err != nil {
			return
		}
	}
	return
}

// RenderText Парсинг множества шаблонов файлов с указанными переменными
func (t *impl) RenderText(wr io.Writer, values interface{}, tpls ...string) error {
	return t.RenderHTML(wr, values, tpls...)
}

// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов io.Reader
func (t *impl) RenderHTMLData(wr io.Writer, values interface{}, buffers ...io.Reader) (err error) {
	var template *pongo2.Template
	var buf []byte
	var i int

	for i = range buffers {
		if buf, err = ioutil.ReadAll(buffers[i]); err != nil {
			return
		}
		if template, err = pongo2.FromString(string(buf)); err != nil {
			return
		}
		if err = template.ExecuteWriter(t.getContext(values), wr); err != nil {
			return
		}
	}

	return
}

// RenderTextData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов io.Reader
func (t *impl) RenderTextData(wr io.Writer, values interface{}, buffers ...io.Reader) error {
	return t.RenderHTMLData(wr, values, buffers...)
}
