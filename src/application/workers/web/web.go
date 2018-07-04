package web // import "application/workers/web"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"errors"
	"time"

	"application/configuration"
	"application/webserver"

	"gopkg.in/webnice/job.v1/job"
	"gopkg.in/webnice/job.v1/types"
)

// Injecting a process object to goroutine management tools
func Init(cnf *configuration.WEBServerConfiguration) {
	job.Get().RegisterWorker(newWorker(cnf))
}

// Конструктор объекта
func newWorker(cnf *configuration.WEBServerConfiguration) types.WorkerInterface {
	return &impl{Cnf: cnf}
}

// Info Функция конфигурации процесса,
// - в функцию передаётся уникальный идентификатор присвоенный процессу
// - функция должна вернуть конфигурацию или nil, если возвращается nil, то применяется конфигурация по умолчанию
func (web *impl) Info(id string) (ret *types.Configuration) {
	web.ID, ret = id, &types.Configuration{
		Autostart:     true,
		Restart:       true,
		Fatality:      true,
		PriorityStart: types.HighPriopity,
		PriorityStop:  types.LowPriopity,
		CancelTimeout: time.Second * 30,
	}
	return
}

// Prepare Функция выполнения действий подготавливающих воркер к работе
// Завершение с ошибкой означает, что процесс не удалось подготовить к запуску
func (web *impl) Prepare() (err error) {
	web.Srv, err = webserver.New(web.Cnf)
	if err != nil {
		err = errors.New("Error create new web server: " + err.Error())
		return
	}
	return
}

// Cancel Функция прерывания работы
func (web *impl) Cancel() (err error) { web.Srv.Stop(); return }

// Worker Функция-реализация процесса, данная функция будет запущена в горутине
// до тех пор пока функция не завершился воркер считается работающим
func (web *impl) Worker() (err error) { web.Srv.Serve(); err = web.Srv.Error(); return }
