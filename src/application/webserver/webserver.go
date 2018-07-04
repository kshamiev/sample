package webserver // import "application/webserver"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1"
import (
	"fmt"
	stddebug "runtime/debug"

	"application/configuration"
	"application/routing"

	// Регистрация функции верификации данных
	_ "gopkg.in/webnice/kit.v1/modules/verify"
)

// Interface is an interface of package
type Interface interface {
	// Serve Запуск сервера
	Serve()

	// Stop Остановка сервера
	Stop()

	// Error Last error
	Error() error
}

// impl is an implementation of package
type impl struct {
	err   error
	web   web.Interface
	cnf   *configuration.WEBServerConfiguration
	Route routing.Interface
}

// New Create new object
func New(cnf *configuration.WEBServerConfiguration) (ret Interface, err error) {
	var wso *impl
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Panic recovery:\n%v\n%s", e.(error), string(stddebug.Stack()))
			ret = nil
		}
	}()
	wso = new(impl)
	wso.cnf = cnf
	wso.web = web.New()
	wso.Route = routing.New(cnf, wso.web.Route()).Logger().Routing()
	wso.err = wso.Route.Error()
	ret, err = wso, wso.err
	return
}

// Error Last error
func (wso *impl) Error() error { return wso.err }

// Serve Запуск сервера
func (wso *impl) Serve() {
	wso.web.ListenAndServeWithConfig(&wso.cnf.Server)
	if wso.err = wso.web.Error(); wso.err != nil {
		return
	}
	wso.web.Wait()
	wso.err = wso.web.Error()
}

// Stop Остановка сервера
func (wso *impl) Stop() { wso.web.Stop(); wso.err = wso.web.Error() }
