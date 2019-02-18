package migrations // import "application/components/migrations"

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
func (mrs *impl) After() []string {
	return []string{
		"application/components/configuration",
		"application/components/logging",
	}
}

// Init Функция инициализации компонента
func (mrs *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	exitCode, mrs.Cfg = workflow.ErrNone, configuration.Get()
	return
}

// Start Выполнение компонента
func (mrs *impl) Start(cmd string) (done bool, exitCode uint8, err error) {
	const commandVersion = `version`
	var command string

	exitCode = workflow.ErrNone
	if cmd == commandVersion {
		return
	}
	log.Info(`Application database migrations apply started`)
	// Поиск утилиты применения миграций
	if command = mrs.migrationsUtility(); command == "" {
		return
	}
	if err = mrs.migrationsMysql(command); err != nil {
		log.Errorf("MySQL migration error: %s", err)
	}
	//	if err = cgn.migrationsClickhouse(command); err != nil {
	//		log.Errorf("ClickHouse migration error: %s", err)
	//	}
	log.Info(`Application database migrations apply completed`)

	return
}

// Stop Функция завершения работы компонента
func (mrs *impl) Stop() (exitCode uint8, err error) {
	return
}
