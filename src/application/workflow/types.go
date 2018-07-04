package workflow // import "application/workflow"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import ()

var (
	singleton *impl
)

// Interface is an interface of package
type Interface interface {
	// Initialize all registered plugins
	Initialize(appVersion string, appBuild string) (exitCode int, err error)

	// Command Running the command with registered plugins
	Command(cmd string) (exitCode int, err error)

	// Shutdown of all plugins
	Shutdown()
}

// impl is an implementation of package
type impl struct {
	// Array of registered plugins
	plugins []PluginInterface
}

// PluginInterface is a interface of plug-in
type PluginInterface interface {
	// Init Функция инициализации плагина
	// должна вернуть код ошибки:
	// - errNone - нет ошибки
	// - любой другой код заставит приложение завершится с ошибкой
	Init(appVersion string, appBuild string) (exitCode int, err error)

	// Command Обработка команды
	// Функция возвращает:
	// - код ошибки. Если возвращен код errNone - нет ошибки
	// - флаг завершения.
	//   Если возвращается true, приложение завершится выполнив команду
	//   Если возвращается false, команда будет передана следующему плагину
	Command(cmd string) (exitCode int, err error, done bool)

	// Shutdown Функция вызывается при завершении работы приложения
	// Каждый плагин, если требуется, по этой функции должен завершить работу
	Shutdown()
}
