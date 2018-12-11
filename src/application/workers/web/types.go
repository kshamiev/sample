package web // import "application/workers/web"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/webserver"
	webserverTypes "application/webserver/types"

	"gopkg.in/webnice/job.v1/types"
)

// Interface is an interface of package
type Interface types.WorkerInterface

// impl is an implementation of package
type impl struct {
	ID  string                        // Уникальный идентификатор процесса
	Wsc *webserverTypes.Configuration // Конфигурация веб сервера
	Wsi webserver.Interface           // Интерфейс веб сервера
}
