package rendering // import "application/modules/rendering"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
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
	options.Renderrer

	// NewStandardTemplateRender Create new html/template implementation render
	NewStandardTemplateRender() standard.Interface

	// NewPongo2TemplateRender Create new pongo2 implementation render
	NewPongo2TemplateRender() pongo2.Interface

	// NewFastTemplateRender Create new fasttemplate implementation render
	NewFastTemplateRender() fasttemplate.Interface

	// SetDefaultEngine Set default engine
	SetDefaultEngine(RenderEngine) Interface
}

// impl is an implementation
type impl struct {
	Option                 options.Option
	FastTemplateRender     fasttemplate.Interface
	StandardTemplateRender standard.Interface
	Pongo2TemplateRender   pongo2.Interface
	DefaultEngine          RenderEngine
}
