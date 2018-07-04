package assets // import "application/controllers/pages/assets"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"

	"application/configuration"
	"application/modules/pages"
)

const (
	ifModifiedSinceTimeFormat = `Mon, 02 Jan 2006 15:04:05 GMT`
)

// Interface is an interface of package
type Interface interface {
	// Assets Statisc files
	Assets(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of package
type impl struct {
	cfg configuration.Interface // Конфигурация приложения
	pgs pages.Interface         // Интерфейс модуля pages
}
