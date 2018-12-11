package web // import "application/workers/web"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"time"

	"application/configuration"
	"application/webserver"
	webserverTypes "application/webserver/types"

	"gopkg.in/webnice/job.v1/job"
	"gopkg.in/webnice/job.v1/types"
)

// Init Injecting a process object to goroutine management tools
func Init(cnf *webserverTypes.Configuration) {
	job.Get().RegisterWorker(New(cnf))
}

// New Конструктор объекта
func New(wsc *webserverTypes.Configuration) types.WorkerInterface {
	return &impl{Wsc: wsc}
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
		CancelTimeout: time.Minute,
	}

	return
}

// Prepare Функция выполнения действий подготавливающих воркер к работе
// Завершение с ошибкой означает, что процесс не удалось подготовить к запуску
func (web *impl) Prepare() (err error) {
	web.Wsi = webserver.New().
		Debug(configuration.Get().Debug())
	if err = web.Wsi.Init(web.Wsc); err != nil {
		err = fmt.Errorf("Create web server (ID: %q) error: %s", web.ID, err)
		return
	}

	return
}

// Cancel Функция прерывания работы
func (web *impl) Cancel() (err error) {
	err = web.Wsi.
		Stop().
		Error()

	return
}

// Worker Функция-реализация процесса, данная функция будет запущена в горутине
// до тех пор пока функция не завершился воркер считается работающим
func (web *impl) Worker() (err error) {
	err = web.Wsi.
		Serve().
		Error()

	return
}
