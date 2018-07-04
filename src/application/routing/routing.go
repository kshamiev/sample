package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/route"
import "gopkg.in/webnice/web.v1/middleware/pprof"
import "gopkg.in/webnice/web.v1/middleware/recovery"
import (
	"application/configuration"
	"application/controllers"
	"application/middleware/gzip"
)

// New Create new object
func New(cnf *configuration.WEBServerConfiguration, rou route.Interface) Interface {
	var rt = new(impl)
	rt.cfg = configuration.Get()
	rt.rou = rou
	rt.SrvCFG = cnf
	return rt
}

// Error Last error
func (rt *impl) Error() error { return rt.err }

// Routing Настройка роутинга
func (rt *impl) Routing() (ret Interface) {
	ret = rt

	// Middleware recovery after panic
	rt.rou.Use(recovery.Recover)

	// Gzip only production mode
	if !rt.cfg.Debug() {
		rt.rou.Use(gzip.Gzip)
	}

	// Статические файлы и шаблоны страниц
	rt.Assets()

	// Настройка роутинга к API
	rt.RoutingAPI()

	// Main custom contorllers
	rt.rou.Handlers().InternalServerError(
		controllers.InternalServerErrorController.InternalServerError,
	)

	// Включение профилирования
	rt.rou.Mount("/debug", pprof.Pprof())

	return
}
