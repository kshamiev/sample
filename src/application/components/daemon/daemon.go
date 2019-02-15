package daemon // import "application/components/daemon"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"

	"application/configuration"
	"application/workers"
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
		"application/componens/bootstrap",
	}
}

// Init Функция инициализации компонента
func (cpn *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	cpn.Cfg = configuration.Get()
	if cpn.Jbo, err = workers.Init(); err != nil {
		exitCode = workflow.ErrCantCreateWorkers
		return
	}

	return
}

// Start Выполнение компонента
func (cpn *impl) Start(cmd string) (done bool, exitCode uint8, err error) {
	const commandDaemon = `daemon`
	var fatality bool

	if cmd != commandDaemon {
		return
	}
	done = true
	cpn.Jbo.
		// Регистрация функции получения ошибок
		RegisterErrorFunc(func(id string, err error) {
			if fatality {
				return
			}
			log.Errorf(" - worker process %q error: %s", id, err)
		}).
		// Регистрация функции получения изменений состояни процессов
		RegisterChangeStateFunc(func(id string, running bool) {
			if !cpn.Cfg.Debug() || fatality {
				return
			}
			if running {
				log.Noticef(" - worker process %q started", id)
			} else {
				log.Noticef(" - worker process %q stopped", id)
			}
		})
	if cpn.Cfg.Debug() {
		log.Info(`Application workers has initialized`)
	}
	// Запуск, запускаются все воркеры с флагом Autostart=true
	if err = cpn.Jbo.Do(); err != nil {
		defer log.Done()
		fatality = true
		err = fmt.Errorf("Workers error: %s", err)
		//log.Criticalf("Workers error: %s", err)
		return
	}
	// На этом этапе все воркеры приложения запущены
	if cpn.Cfg.Debug() {
		log.Info(`Application workers started successfully`)
	}
	// Ожидание завершения всех воркеров
	cpn.Jbo.Wait()

	return
}

// Stop Функция завершения работы компонента
func (cpn *impl) Stop() (exitCode uint8, err error) {
	return
}
