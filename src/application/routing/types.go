package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/route"
import (
	"application/configuration"
)

// Interface is an interface of package
type Interface interface {
	// Logger Настройка логгера
	Logger() Interface

	// Routing Настройка роутинга
	Routing() Interface

	// Error Last error
	Error() error
}

// impl is an implementation of package
type impl struct {
	rou    route.Interface
	err    error
	SrvCFG *configuration.WEBServerConfiguration
	cfg    configuration.Interface
}
