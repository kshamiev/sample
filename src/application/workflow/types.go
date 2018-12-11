package workflow // import "application/workflow"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

var (
	singleton *impl
)

// Interface is an interface of package
type Interface interface {
	// Debug Enable or disable debug mode
	Debug(d bool) Interface

	// Init initialize all registered components
	Init(appVersion string, appBuild string) (exitCode uint8, err error)

	// Start Running start command in the all registered components
	// if all of components ignored the command, runs the usageFunc
	Start(cmd string, usageFunc func()) (exitCode uint8, err error)

	// Stop of all registered components
	Stop() (exitCode uint8, err error)
}

// impl is an implementation of package
type impl struct {
	debug      bool                 // =true - debug mode
	Components []ComponentInterface // All registered components
}

// ComponentInterface is a interface of an application component
type ComponentInterface interface {
	// After Возвращает массив зависимостей (аналог After в systemd), массив состоит из названия пакетов, после которых
	// должен выполняться компонен. Если массив пустой, то компонент выполняется в порядке регистрации
	// Приоритет влияет на все стадии выполнения компонента (Init, Start, Stop)
	After() []string

	// Init Функция инициализации компонента, должна вернуть код ошибки:
	// - errNone - нет ошибки
	// - любой другой код заставит приложение завершится с ошибкой
	Init(appVersion string, appBuild string) (exitCode uint8, err error)

	// Start Выполнение компонента
	// Функция возвращает:
	// - (exitCode) код ошибки. Если возвращен код errNone - нет ошибки
	// - (err) Ошибка в виде интерфейса error
	// - (done) флаг завершения работы
	//   - Если возвращается true, выполнение компонентов прерывается и приложение завершается
	//   - Если возвращается false, выполняется следующий компонент
	Start(cmd string) (done bool, exitCode uint8, err error)

	// Stop Функция вызывается перед завершением приложения
	// Каждый компонент, при вызове этой функции должен остановить свои зависимости и
	// завершить работу в нормальном режиме
	Stop() (exitCode uint8, err error)
}
