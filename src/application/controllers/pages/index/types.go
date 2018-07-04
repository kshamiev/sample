package index // import "application/controllers/pages/index"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"

	"application/configuration"
	"application/modules/pages"
)

const (
	_IfModifiedSinceTimeFormat = `Mon, 02 Jan 2006 15:04:05 GMT`
)

// Interface is an interface of package
type Interface interface {
	// Index page
	Index(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of package
type impl struct {
	cfg configuration.Interface // Конфигурация приложения
	pgs pages.Interface         // Интерфейс модуля pages
}
