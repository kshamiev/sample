package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	webserverTypes "application/webserver/types"

	"gopkg.in/webnice/web.v1/route"
)

// Interface is an interface of package
type Interface interface {
	// Logger Настройка логирования запросов web сервера
	Logger() Interface

	// Routing Настройка роутинга
	Routing() Interface

	// Debug Включение или отключение режима отладки
	Debug(d bool) Interface

	// Stop Stopping operations, closing connections, or flushing the cache
	Stop() error

	// Error Last error
	Error() error
}

// impl is an implementation of package
type impl struct {
	err   error
	debug bool
	Wsc   *webserverTypes.Configuration
	Rou   route.Interface
}
