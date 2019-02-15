package configuration // import "application/components/configuration"

import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"

	"application/configuration"
	"application/workflow"
)

func init() { workflow.Register(New()) }

// New Create object and return interface
func New() Interface {
	return new(impl)
}

// After Возвращает массив зависимостей - имена пакетов компонентов, которые должны быть запущены до этого компонента
func (cpn *impl) After() []string {
	return []string{
		"application/componens/environment",
	}
}

// Init Функция инициализации компонента
func (cpn *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	cpn.Cfg = configuration.Version(configuration.NewVersion(appVersion, appBuild))
	// Initializing the application configuration
	if err = cpn.Cfg.Init(); err != nil {
		exitCode = workflow.ErrInitConfiguration
		return
	}
	// Set debug mode
	workflow.Get().Debug(cpn.Cfg.Debug())

	return
}

// Start Выполнение компонента
func (cpn *impl) Start(cmd string) (done bool, exitCode uint8, err error) {
	const commandVersion = `version`

	// Print version and exit
	if cmd == commandVersion {
		fmt.Println(cpn.Cfg.Version().String())
		done = true
		return
	}
	// Go to the working directory
	if err = cpn.Cfg.GotoWorkingDirectory(); err != nil {
		err = fmt.Errorf("%q: %s", cpn.Cfg.WorkingDirectory(), err)
		exitCode, done = workflow.ErrCantChangeWorkDirectory, true
		return
	}
	// Only in the debug mode, the current configuration is displayed
	if cpn.Cfg.Debug() {
		log.Debugf("Application configuration:\n%s", debug.DumperString(cpn.Cfg.Configuration()))
	}

	return
}

// Stop Функция завершения работы компонента
func (cpn *impl) Stop() (exitCode uint8, err error) {
	return
}
