package rendering // import "application/modules/rendering"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"io"

	"application/modules/rendering/fasttemplate"
	"application/modules/rendering/options"
	"application/modules/rendering/pongo2"
	"application/modules/rendering/standard"
)

const (
	// RenderStandardTemplate Standard Template engine
	RenderStandardTemplate RenderEngine = 1

	// RenderPongo2Template Pongo2 Template engine
	RenderPongo2Template RenderEngine = 2

	// RenderFastTemplate Fast Template engine
	RenderFastTemplate RenderEngine = 3
)

type (
	// Option Rendering option
	Option options.Option

	// RenderEngine type define
	RenderEngine int8
)

// Interface is an interface
type Interface interface {
	// NewStandardTemplateRender Create new html/template implementation render
	NewStandardTemplateRender() standard.Interface

	// NewPongo2TemplateRender Create new pongo2 implementation render
	NewPongo2TemplateRender() pongo2.Interface

	// NewFastTemplateRender Create new fasttemplate implementation render
	NewFastTemplateRender() fasttemplate.Interface

	// SetDefaultEngine Set default engine
	SetDefaultEngine(RenderEngine) Interface

	// RenderHTML Парсинг множества шаблонов файлов с указанными переменными
	RenderHTML(io.Writer, interface{}, ...string) error

	// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов *bytes.Buffer
	RenderHTMLData(io.Writer, interface{}, ...*bytes.Buffer) error
}

// impl is an implementation
type impl struct {
	Option                 options.Option
	FastTemplateRender     fasttemplate.Interface
	StandardTemplateRender standard.Interface
	Pongo2TemplateRender   pongo2.Interface
	DefaultEngine          RenderEngine
}
