package web // import "application/workers/web"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/webserver"

	"gopkg.in/webnice/job.v1/types"
)

// Interface is an interface of package
type Interface interface {
	// Info Функция конфигурации процесса,
	// - в функцию передаётся уникальный идентификатор присвоенный процессу
	// - функция должна вернуть конфигурацию или nil, если возвращается nil, то применяется конфигурация по умолчанию
	Info(id string) *types.Configuration

	// Cancel Функция прерывания работы
	Cancel() error

	// Prepare Функция выполнения действий подготавливающих воркер к работе
	// Завершение с ошибкой означает, что процесс не удалось подготовить к запуску
	Prepare() error

	// Worker Функция-реализация процесса, данная функция будет запущена в горутине
	// до тех пор пока функция не завершился воркер считается работающим
	Worker() error
}

// impl is an implementation of package
type impl struct {
	ID  string                                // Уникальный идентификатор процесса
	Cnf *configuration.WEBServerConfiguration // Конфигурация веб сервера
	Srv webserver.Interface                   // Интерфейс веб сервера
}
