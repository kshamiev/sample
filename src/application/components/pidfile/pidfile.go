package pidfile // import "application/components/pidfile"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"

	"application/configuration"
	"application/modules/pidfile"
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
		"application/components/configuration",
		"application/components/logging",
	}
}

// Init Функция инициализации компонента
func (cpn *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	cpn.Cfg = configuration.Get()
	if cpn.Cfg.Debug() && cpn.Cfg.PidFile() != "" {
		log.Debugf("create PID file: %q", cpn.Cfg.PidFile())
	}
	if cpn.Cfg.PidFile() != "" {
		cpn.Pid, err = pidfile.New(cpn.Cfg.PidFile()).Error()
	}

	return
}

// Start Выполнение компонента
func (cpn *impl) Start(cmd string) (done bool, exitCode uint8, err error) {
	const commandDaemon = `daemon`

	if cmd != commandDaemon || cpn.Pid == nil {
		if cpn.Cfg.Debug() {
			log.Debugf("PID file not created. Command %q, pid: %v", cmd, cpn.Pid)
		}
		return
	}
	// Создание PID файла, блокировка файла на запись и удаление (по возможности)
	if err = cpn.Pid.Lock(); err != nil {
		exitCode, err = workflow.ErrPidFile, fmt.Errorf("Create PID file error: %s", err)
		return
	}
	if cpn.Cfg.Debug() {
		log.Info(`Application PID file has created successfully`)
	}

	return
}

// Stop Функция завершения работы компонента
func (cpn *impl) Stop() (exitCode uint8, err error) {
	if cpn.Pid == nil {
		return
	}
	// Удаление PID файла при завершении приложения
	if err = cpn.Pid.Unlock(); err != nil {
		exitCode, err = workflow.ErrPidFileDelete, fmt.Errorf("Delete PID file error: %s", err)
		return
	}
	if cpn.Cfg.Debug() {
		log.Info(`Application PID file deleted successfully`)
	}

	return
}
