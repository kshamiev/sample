package logging // import "application/components/logging"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/workflow"

	r "gopkg.in/webnice/log.v2/receiver"
	s "gopkg.in/webnice/log.v2/sender"
)

func init() { workflow.Register(New()) }

// New Create object and return interface
func New() Interface {
	return new(impl)
}

// After Возвращает массив зависимостей - имена пакетов компонентов, которые должны быть запущены до этого компонента
func (cpn *impl) After() []string {
	return []string{
		"application/components/configuration",
	}
}

// Init Функция инициализации компонента
func (cpn *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	var receiver s.Receiver

	cpn.Cfg = configuration.Get()
	if cpn.Cfg.Log() == nil {
		return
	}
	if cpn.Cfg.Log().GraylogEnable {
		receiver = r.GelfReceiver.
			SetAddress(cpn.Cfg.Log().GraylogProto, cpn.Cfg.Log().GraylogAddress, cpn.Cfg.Log().GraylogPort).
			SetCompression("gzip").
			Receiver
		s.Gist().AddSender(receiver)
	}
	if cpn.Cfg.Debug() {
		s.Gist().AddSender(r.Default.Receiver)
	}
	log.Gist().StandardLogSet()
	if cpn.Cfg.Debug() {
		log.Info(`Logging system has initialized successfully`)
	}

	return
}

// Start Выполнение компонента
func (cpn *impl) Start(cmd string) (done bool, exitCode uint8, err error) {
	return
}

// Stop Функция завершения работы компонента
func (cpn *impl) Stop() (exitCode uint8, err error) {
	if cpn.Cfg.Debug() {
		log.Info(`Logging system has been shutdown`)
	}
	log.Gist().StandardLogUnset()
	log.Done()
	return
}
