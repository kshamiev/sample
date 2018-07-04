package configuration // import "application/configuration"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"os/user"

	"application/configuration/semver"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Interface is an interface of configuration
type Interface interface {
	// Version Return version of the application
	Version() *semver.Version

	// Debug Return debug mode status of the application
	Debug() bool

	// Init initialization configuration
	Init() error

	// InitJWT Загрузка, проверка, инициализация или генерация ключей для токенов JWT
	//InitJWT() error

	// Configuration Return full configuration object of the application
	Configuration() *Application

	//

	// Command Return the current command with which the application is running
	Command() string

	// AppName Return application name
	AppName() string

	// WorkingDirectory Return work path of the application
	WorkingDirectory() string

	// PidFile Return name and path to the PID file
	PidFile() string

	// TempPath Return path to the temporary files folder
	TempPath() string

	// LogConfiguration Путь и имя файла конфигурации системы логирования
	LogConfiguration() string

	// LogPath Return path to the log folder
	LogPath() string

	// Log configuration object
	Log() *ApplicationLog

	// StateFile Return name and path to the file of state of application
	StateFile() string

	// SocketFile Return name and path to the file of unix socket communication
	SocketFile() string

	//

	// GotoWorkingDirectory Change current working directory to the application's working directory
	GotoWorkingDirectory() error
}

// impl is an implementation of repository
type impl struct {
	version          *semver.Version // Версия приложения
	debug            bool            // =true - приложение запущено в режиме отладки
	args             args            // Разбор значений коммандной строки
	argsCommand      string          // Команда с которой запущено приложение
	user             *user.User      // Объект с информацией о текущем пользователе
	appFilePath      string          // Путь и имя файла загруженной конфигурации
	appConfiguration *Application    // Объект конфигурации приложения
	appLog           *ApplicationLog // Объект конфигурации системы логирования
}

// args Описание принимаемых приложением аргументов коммандной строки
type args struct {
	App      *kingpin.Application // Application object
	Version  *kingpin.CmdClause   // Команда version - приложение отображает версию и завершается с кодом ошибки 0
	Daemon   *kingpin.CmdClause   // Команда daemon - приложение запускает сервер, открывает порты и переходи в режим сервиса
	Cli      *kingpin.CmdClause   // Команда cli - коммандный режим, приложение ведёт себя как клиент сервиса, передавая сервису команды из коммандной строки (по умолчанию)
	FilePath string               // Флаг --conf - путь и имя конфигурационного файла
}
