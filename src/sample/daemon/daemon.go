package daemon

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/workers"
	"application/workflow"

	"gopkg.in/webnice/job.v1/job"
)

// impl is an implementation of package
type impl struct {
	jbo job.Interface
}

func init() {
	workflow.Register(new(impl))
}

// Init Plug-in initialization function
func (cgn *impl) Init(appVersion string, appBuild string) (exitCode int, err error) {
	exitCode = workflow.ErrNone
	if cgn.jbo, err = workers.Init(); err != nil {
		exitCode = workflow.ErrCantCreateWorkers
		return
	}

	return
}

// Command Processing the command
func (cgn *impl) Command(cmd string) (exitCode int, err error, done bool) {

	exitCode = workflow.ErrNone
	if cmd != `daemon` {
		return
	}
	done = true
	cgn.jbo.
		// Регистрация функции получения ошибок
		RegisterErrorFunc(func(id string, err error) {
			log.Errorf("Worker process %q error: %s", id, err)
		}).
		// Регистрация функции получения изменений состояни процессов
		RegisterChangeStateFunc(func(id string, running bool) {
			if !configuration.Get().Debug() {
				return
			}
			if running {
				log.Noticef("Worker process %q started", id)
			} else {
				log.Noticef("Worker process %q stopped", id)
			}
		})
	log.Info(`Application workers has initialized successfully`)
	// Запуск, запускаются все процессы с флагом Autostart=true
	if err = cgn.jbo.Do(); err != nil {
		log.Criticalf("Workers error: %s", err)
		return
	}
	// На этом этапе все компоненты приложения запущены
	log.Info(`Application started successfully`)
	// Ожидание завершения всех процессов
	cgn.jbo.Wait()

	return
}

// Done The function is called when the application terminates
func (cgn *impl) Shutdown() {
	if configuration.Get().Debug() {
		log.Debugf("Beginning shutdown")
	}
}
