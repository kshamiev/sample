package pidfile

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/modules/pidfile"
	"application/workflow"
)

// impl is an implementation of package
type impl struct {
	conf configuration.Interface
	pid  pidfile.Interface
}

func init() {
	workflow.Register(new(impl))
}

// Init Plug-in initialization function
func (cgn *impl) Init(appVersion string, appBuild string) (exitCode int, err error) {

	exitCode = workflow.ErrNone
	cgn.conf = configuration.Get()
	if cgn.conf.PidFile() != "" {
		cgn.pid = pidfile.New(cgn.conf.PidFile())
	}

	return
}

// Command Processing the command
func (cgn *impl) Command(cmd string) (exitCode int, err error, done bool) {

	exitCode = workflow.ErrNone
	if cmd == `daemon` && cgn.pid != nil {
		// Создание PID файла, блокировка файла на запись и удаление (по возможности)
		if err = cgn.pid.Lock(); err != nil {
			log.Warningf("Error create PID file: %s", err)
			exitCode = workflow.ErrPidFile
			return
		}
		log.Info(`Application PID file has created successfully`)
	}

	return
}

// Done The function is called when the application terminates
func (cgn *impl) Shutdown() {
	if cgn.pid == nil {
		return
	}
	// Удаление PID файла при завершении приложения
	if err := cgn.pid.Unlock(); err != nil {
		log.Warningf("Error delete PID file: %s", err)
		return
	}
}
