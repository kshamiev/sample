package configuration // import "application/configuration"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/kit.v1/modules/db"
import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"strings"

	"application/configuration/semver"

	"gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	singleton     *impl                         // configuration object
	rexSlashAtEnd = regexp.MustCompile(`(/+)$`) // regex to clean slash at end of string
)

// initialization at startup
func init() {
	var err error
	singleton = new(impl)
	initCLIArgs(singleton)
	if singleton.user, err = user.Current(); err != nil {
		log.Errorf("Unable to get current user: %s", err.Error())
	}
}

// Version Set version and return configuration interface
func Version(version *semver.Version) Interface { singleton.version = version; return singleton }

// Get Return singleton configuration interface
func Get() Interface { return singleton }

// Version Get application vertion
func (cnf *impl) Version() *semver.Version { return cnf.version }

// Debug mode status
func (cnf *impl) Debug() bool { return cnf.debug }

// Command Команда с которой запущено приложение
func (cnf *impl) Command() string { return cnf.argsCommand }

// AppName Return application name
func (cnf *impl) AppName() string { return cnf.appConfiguration.ApplicationName }

// WorkingDirectory Return work path of the application
func (cnf *impl) WorkingDirectory() string { return cnf.appConfiguration.WorkingDirectory }

// PidFile Return name and path to the PID file
func (cnf *impl) PidFile() string {
	if cnf.appConfiguration == nil {
		return ""
	}
	return cnf.appConfiguration.PidFile
}

// TempPath Return path to the temporary files folder
func (cnf *impl) TempPath() string { return cnf.appConfiguration.TempPath }

// LogConfiguration Return name and path to the file of logging configuration
func (cnf *impl) LogConfiguration() string { return cnf.appConfiguration.LogConfiguration }

// LogPath Return path to the log folder
func (cnf *impl) LogPath() string { return cnf.appConfiguration.LogPath }

// Log configuration object
func (cnf *impl) Log() *ApplicationLog { return cnf.appLog }

// StateFile Return name and path to the file of state of application
func (cnf *impl) StateFile() string { return cnf.appConfiguration.StateFile }

// SocketFile Return name and path to the file of unix socket communication
func (cnf *impl) SocketFile() string { return cnf.appConfiguration.SocketFile }

// MustConfigurationInitialization Проверка обязательного наличия инициализированной конфигурации
func (cnf *impl) MustConfigurationInitialization() {
	if cnf.appConfiguration == nil {
		log.Fatal("Configuration is not initialized. Please first run Init()\n")
	}
}

// GotoWorkingDirectory Change current working directory to the application's working directory
func (cnf *impl) GotoWorkingDirectory() (err error) {
	if cnf.appConfiguration == nil {
		err = fmt.Errorf("Configuration is not initialized. Please first run Init()")
		return
	}
	if cnf.appConfiguration.WorkingDirectory == "" {
		err = fmt.Errorf("Working directory is empty")
		return
	}
	err = os.Chdir(cnf.appConfiguration.WorkingDirectory)
	return
}

// Configuration return configuration object
func (cnf *impl) Configuration() *Application {
	return cnf.appConfiguration
}

// Init Initialization configuration
func (cnf *impl) Init() (err error) {
	var configPaths []string
	var configData []byte

	cnf.argsCommand = strings.ToLower(kingpin.Parse())

	// Version
	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case `version`, `help`, `--help`, `--help-long`, `--help-man`:
			return
		}
	}

	// Поиск конфигурационного файла
	configPaths = cnf.InitPaths()
	if cnf.appFilePath, err = cnf.SeekConfigFile(configPaths); err != nil {
		return
	}

	// Read all file data
	if configData, err = ioutil.ReadFile(cnf.appFilePath); err != nil {
		err = fmt.Errorf("Can't read configuration from file: %s", err.Error())
		return
	}

	// Unmarshal yaml configuration data to structure
	cnf.appConfiguration = new(Application)
	if err = yaml.Unmarshal(configData, cnf.appConfiguration); err != nil {
		err = fmt.Errorf("Can't unmarshal data from yaml: %s", err.Error())
		return
	}

	// Таймаут CLI команды
	//cnf.appConfiguration.CliTimeout = cnf.argsTimeout

	// Дефолтовые значения
	cnf.MakeDefaults()

	// Переопределение переменных значениями из os environment
	cnf.Environment()

	// Очистка/преобразование путей
	cnf.CleanAllPath()

	// Установка конфигурации подключения к базе данных по умолчанию
	db.Default(&cnf.appConfiguration.Database)

	// Cli
	//initCLIHelp(cnf.appConfiguration.ApplicationName)
	//cnf.argsCommand = strings.ToLower(kingpin.Parse())

	// Конфигурация системы логирования
	if cnf.appConfiguration.LogConfiguration != "" {
		err = cnf.InitLog(cnf.appConfiguration.LogConfiguration)
		if err != nil {
			return
		}
	}

	return
}
