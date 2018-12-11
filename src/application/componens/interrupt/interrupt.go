package interrupt // import "application/componens/interrupt"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"os"

	"application/configuration"
	"application/modules/interrupt"
	"application/workflow"

	"gopkg.in/webnice/job.v1/job"
)

func init() { workflow.Register(New()) }

// New Create object and return interface
func New() Interface {
	return new(impl)
}

// After Возвращает массив зависимостей - имена пакетов компонентов, которые должны быть запущены до этого компонента
func (cpn *impl) After() []string {
	return []string{
		"application/componens/configuration",
		"application/componens/logging",
	}
}

// Init Функция инициализации компонента
func (cpn *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	cpn.Cfg = configuration.Get()
	return
}

// Start Выполнение компонента
func (cpn *impl) Start(cmd string) (done bool, exitCode uint8, err error) {
	cpn.Itp = interrupt.New(func(sig os.Signal) {
		log.Alertf(`Received OS interrupt signal`+": %s", sig.String())
		job.Cancel()
	})
	cpn.Itp.Start()
	if cpn.Cfg.Debug() {
		log.Debugf("Interception of system interruptions started successfully")
	}

	return
}

// Stop Функция завершения работы компонента
func (cpn *impl) Stop() (exitCode uint8, err error) {
	if cpn.Itp == nil {
		return
	}
	cpn.Itp.Stop()
	if cpn.Cfg.Debug() {
		log.Debugf("Interception of system interruptions has been shutdown")
	}

	return
}
