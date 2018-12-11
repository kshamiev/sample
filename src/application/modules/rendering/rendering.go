package rendering // import "application/modules/rendering"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io"

	"application/modules/rendering/fasttemplate"
	"application/modules/rendering/options"
	"application/modules/rendering/pongo2"
	"application/modules/rendering/standard"
)

// New creates new object
func New(args ...Option) Interface {
	var o = new(impl)
	for _, arg := range args {
		o.Option = options.Option(arg)
	}
	o.SetDefaultEngine(RenderPongo2Template)
	return o
}

// SetDefaultEngine Set default engine
func (r *impl) SetDefaultEngine(e RenderEngine) Interface {
	r.DefaultEngine = e
	switch r.DefaultEngine {
	case RenderFastTemplate:
		r.NewFastTemplateRender()
	case RenderStandardTemplate:
		r.NewStandardTemplateRender()
	case RenderPongo2Template:
		r.NewPongo2TemplateRender()
	}
	return r
}

// NewStandardTemplateRender Create new html/template implementation render
func (r *impl) NewStandardTemplateRender() standard.Interface {
	if r.StandardTemplateRender == nil {
		r.StandardTemplateRender = standard.New(r.Option)
	}
	return r.StandardTemplateRender
}

// NewPongo2TemplateRender Create new pongo2 implementation render
func (r *impl) NewPongo2TemplateRender() pongo2.Interface {
	if r.Pongo2TemplateRender == nil {
		r.Pongo2TemplateRender = pongo2.New(r.Option)
	}
	return r.Pongo2TemplateRender
}

// NewFastTemplateRender Create new fasttemplate implementation render
func (r *impl) NewFastTemplateRender() fasttemplate.Interface {
	if r.FastTemplateRender == nil {
		r.FastTemplateRender = fasttemplate.New(r.Option)
	}
	return r.FastTemplateRender
}

// RenderHTML Парсинг множества шаблонов файлов с указанными переменными
func (r *impl) RenderHTML(wr io.Writer, values interface{}, tpl ...string) (err error) {
	switch r.DefaultEngine {
	case RenderStandardTemplate:
		err = r.StandardTemplateRender.RenderHTML(wr, values, tpl...)
	case RenderPongo2Template:
		err = r.Pongo2TemplateRender.RenderHTML(wr, values, tpl...)
	case RenderFastTemplate:
		err = r.FastTemplateRender.RenderHTML(wr, values, tpl...)
	}
	if err != nil {
		err = fmt.Errorf("RenderHTML(%v) error: %s", r.DefaultEngine, err)
	}
	return
}

// RenderText Парсинг множества шаблонов файлов с указанными переменными
func (r *impl) RenderText(wr io.Writer, values interface{}, tpl ...string) (err error) {
	switch r.DefaultEngine {
	case RenderStandardTemplate:
		err = r.StandardTemplateRender.RenderText(wr, values, tpl...)
	case RenderPongo2Template:
		err = r.Pongo2TemplateRender.RenderText(wr, values, tpl...)
	case RenderFastTemplate:
		err = r.FastTemplateRender.RenderText(wr, values, tpl...)
	}
	if err != nil {
		err = fmt.Errorf("RenderText(%v) error: %s", r.DefaultEngine, err)
	}
	return
}

// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов io.Reader
func (r *impl) RenderHTMLData(wr io.Writer, values interface{}, buffers ...io.Reader) (err error) {
	switch r.DefaultEngine {
	case RenderStandardTemplate:
		err = r.StandardTemplateRender.RenderHTMLData(wr, values, buffers...)
	case RenderPongo2Template:
		err = r.Pongo2TemplateRender.RenderHTMLData(wr, values, buffers...)
	case RenderFastTemplate:
		err = r.FastTemplateRender.RenderHTMLData(wr, values, buffers...)
	}
	if err != nil {
		err = fmt.Errorf("RenderHTMLData(%v) error: %s", r.DefaultEngine, err)
	}
	return
}

// RenderTextData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов io.Reader
func (r *impl) RenderTextData(wr io.Writer, values interface{}, buffers ...io.Reader) (err error) {
	switch r.DefaultEngine {
	case RenderStandardTemplate:
		err = r.StandardTemplateRender.RenderTextData(wr, values, buffers...)
	case RenderPongo2Template:
		err = r.Pongo2TemplateRender.RenderTextData(wr, values, buffers...)
	case RenderFastTemplate:
		err = r.FastTemplateRender.RenderTextData(wr, values, buffers...)
	}
	if err != nil {
		err = fmt.Errorf("RenderTextData(%v) error: %s", r.DefaultEngine, err)
	}
	return
}
