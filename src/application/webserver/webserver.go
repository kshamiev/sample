package webserver // import "application/webserver"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	stddebug "runtime/debug"

	"application/routing"
	webserverTypes "application/webserver/types"

	// Регистрация функции верификации данных
	_ "gopkg.in/webnice/kit.v1/modules/verify"

	"gopkg.in/webnice/web.v1"
)

// New Create new object
func New() Interface {
	var wso = &impl{
		Wsi: web.New(),
	}

	return wso
}

// Init Инициализация веб сервера, настройка роутинга, подготовка к запуску
func (wso *impl) Init(wsc *webserverTypes.Configuration) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Panic recovery:\n%v\n%s", e.(error), string(stddebug.Stack()))
		}
	}()
	wso.Wsc, wso.Rti = wsc, routing.New(wsc, wso.Wsi.Route()).
		Debug(wso.debug).
		Logger().
		Routing()
	if err, wso.err = wso.Rti.Error(), wso.Rti.Error(); err != nil {
		return
	}
	// Прогрев кеша для DocumentRoot
	if err = wso.CacheWarmingUp(wso.Wsc.DocumentRoot); err != nil {
		return
	}
	// Прогрев кеша для шаблонов html страниц
	if err = wso.CacheWarmingUp(wso.Wsc.Pages); err != nil {
		return
	}

	return
}

// Serve Запуск сервера
func (wso *impl) Serve() Interface {
	wso.Wsi.
		ListenAndServeWithConfig(&wso.Wsc.Server)
	if wso.err = wso.Wsi.Error(); wso.err != nil {
		return wso
	}
	wso.err = wso.Wsi.Wait().Error()

	return wso
}

// Stop Остановка сервера
func (wso *impl) Stop() Interface {
	wso.Wsi.Stop()
	if wso.err = wso.Rti.Stop(); wso.err != nil {
		log.Warningf("Route stop with error: %s", wso.err)
	}
	return wso
}

// Debug Enable or disable debug mode
func (wso *impl) Debug(d bool) Interface { wso.debug = d; return wso }

// Error Last error
func (wso *impl) Error() error { return wso.err }
