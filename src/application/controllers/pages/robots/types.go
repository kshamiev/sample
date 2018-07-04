package robots // import "application/controllers/pages/robots"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"

	"application/configuration"
)

const (
	_TemplateName = `robots.txt`
)

// Interface is an interface of package
type Interface interface {
	// RobotsTxt robots.txt
	RobotsTxt(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of package
type impl struct {
	cfg          configuration.Interface // Конфигурация приложения
	url          string                  // Публичный адрес сервера
	DocumentRoot string                  // Путь к папке document root
}

// Переменные шаблонизатора
type templateVars struct {
	RequestDomain string
	RequestURL    string
}
