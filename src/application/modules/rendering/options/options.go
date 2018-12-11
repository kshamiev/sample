package options // import "application/modules/rendering/options"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
)

// Option Templating settings
type Option struct {
	// Directory to load templates. Default is "assets/www"
	Directory string

	// Reload to reload templates everytime. Default is false - сaching templates in memory
	Reload bool
}

// Renderrer interface
type Renderrer interface {
	// RenderHTML Парсинг множества шаблонов файлов с указанными переменными
	RenderHTML(io.Writer, interface{}, ...string) error

	// RenderText Парсинг множества шаблонов файлов с указанными переменными
	RenderText(io.Writer, interface{}, ...string) error

	// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов io.Reader
	RenderHTMLData(io.Writer, interface{}, ...io.Reader) error

	// RenderHTMLData Парсинг множества шаблонов с указанием переменных. Все шаблоны указываются в виде объектов io.Reader
	RenderTextData(io.Writer, interface{}, ...io.Reader) error
}
