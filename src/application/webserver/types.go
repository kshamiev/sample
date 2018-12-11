package webserver // import "application/webserver"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/routing"
	webserverTypes "application/webserver/types"

	// Регистрация функции верификации данных
	_ "gopkg.in/webnice/kit.v1/modules/verify"

	"gopkg.in/webnice/web.v1"
)

// Interface is an interface of package
type Interface interface {
	// Init Инициализация веб сервера, настройка роутинга, подготовка к запуску
	Init(wsc *webserverTypes.Configuration) (err error)

	// Serve Запуск сервера
	Serve() Interface

	// Stop Остановка сервера
	Stop() Interface

	// Debug Enable or disable debug mode
	Debug(d bool) Interface

	// Error Last error
	Error() error
}

// impl is an implementation of package
type impl struct {
	err   error
	debug bool
	Wsc   *webserverTypes.Configuration
	Wsi   web.Interface
	Rti   routing.Interface
}
