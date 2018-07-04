package migrations

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/workflow"
)

// impl is an implementation of package
type impl struct {
	conf configuration.Interface
}

func init() {
	workflow.Register(new(impl))
}

// Init Plug-in initialization function
func (cgn *impl) Init(appVersion string, appBuild string) (exitCode int, err error) {
	exitCode = workflow.ErrNone
	cgn.conf = configuration.Get()
	return
}

// Command Processing the command
func (cgn *impl) Command(cmd string) (exitCode int, err error, done bool) {
	var command string

	exitCode = workflow.ErrNone
	if cmd == `version` {
		return
	}
	log.Info(`Application database migrations apply started`)
	// Поиск утилиты применения миграций
	if command = cgn.migrationsUtility(); command == "" {
		return
	}
	if err = cgn.migrationsMysql(command); err != nil {
		log.Errorf("MySQL migration error: %s", err)
	}
	//	if err = cgn.migrationsClickhouse(command); err != nil {
	//		log.Errorf("ClickHouse migration error: %s", err)
	//	}
	log.Info(`Application database migrations apply completed`)

	return
}

// Done The function is called when the application terminates
func (cgn *impl) Shutdown() {
}
