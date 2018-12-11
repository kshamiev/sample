package types // import "application/webserver/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"gopkg.in/webnice/web.v1"
)

// Configuration WEB server configuration structure
type Configuration struct {
	Server            web.Configuration `yaml:"Server"`            // Конфигурация WEB сервера
	DocumentRoot      string            `yaml:"DocumentRoot"`      // Корень http сервера
	Pages             string            `yaml:"Pages"`             // Расположение html шаблонов страниц, чей код генерируется на стороне сервера
	ErrorCodeTemplate map[int]string    `yaml:"ErrorCodeTemplate"` // Шаблоны для Content-type text/html соответствующие кодам http ответа, файл шаблона ищется в Template
}
