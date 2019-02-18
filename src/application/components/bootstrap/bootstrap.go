package bootstrap // import "application/components/bootstrap"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/workflow"
)

func init() { workflow.Register(New()) }

// New Create object and return interface
func New() Interface {
	return new(impl)
}

// After Возвращает массив зависимостей - имена пакетов компонентов, которые должны быть запущены до этого компонента
func (bst *impl) After() []string {
	return []string{
		"application/components/environment",
		"application/components/configuration",
		"application/components/logging",
		"application/components/interrupt",
		"application/components/pidfile",
	}
}

// Init Функция инициализации компонента
func (bst *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	bst.Cfg = configuration.Get()
	return
}

// Start Выполнение компонента
func (bst *impl) Start(cmd string) (done bool, exitCode uint8, err error) {
	// На этом этапе все основные зависимости приложения загружены и проинициализированы
	if bst.Cfg.Debug() {
		log.Info(`Application started successfully`)
	}

	return
}

// Stop Функция завершения работы компонента
func (bst *impl) Stop() (exitCode uint8, err error) {
	return
}
