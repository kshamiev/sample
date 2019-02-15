package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/controllers"
	"application/middleware/gzip"
	webserverTypes "application/webserver/types"

	"gopkg.in/webnice/web.v1/middleware/contentTypeDefault"
	"gopkg.in/webnice/web.v1/middleware/pprof"
	"gopkg.in/webnice/web.v1/middleware/recovery"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/route"
)

// New Create new object and return interface
func New(wsc *webserverTypes.Configuration, rou route.Interface) Interface {
	var rt = &impl{
		Rou: rou,
		Wsc: wsc,
	}
	return rt
}

// Debug Enable or disable debug mode
func (rt *impl) Debug(d bool) Interface { rt.debug = d; return rt }

// Error Last error
func (rt *impl) Error() error { return rt.err }

// Routing Настройка роутинга
func (rt *impl) Routing() Interface {
	// Middleware of panic recovery
	rt.Rou.
		Use(recovery.Recover)
	// Gzip all content only in production mode
	if !rt.debug {
		rt.Rou.
			Use(gzip.Gzip)
	}
	// Статические файлы и шаблоны страниц
	rt.Assets()
	// Настройка роутинга к API
	rt.API()
	// Custom controller of web server, of exception for internal server error
	rt.Rou.
		Handlers().
		InternalServerError(controllers.InternalServerErrorController.InternalServerError)
	// Включение профилирования
	rt.Rou.Subroute("/debug", func(sr route.Interface) {
		sr.Use(contentTypeDefault.New(mime.TextHTMLCharsetUTF8).Handler)
		sr.Mount("/", pprof.Pprof())
	})

	return rt
}

// Stop Stopping operations, closing connections, or flushing the cache
func (rt *impl) Stop() (err error) {
	err = controllers.ResourceController.
		Sitemap().
		Close()

	return
}
