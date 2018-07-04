package configuration // import "application/configuration"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"os"
)

// Environment Overriding configuration variables from os environment values
func (cnf *impl) Environment() {
	const (
		_EnvName    = `APPLICATION_NAME`
		_EnvHome    = `APPLICATION_HOME`
		_EnvPid     = `APPLICATION_PID`
		_EnvTemp    = `APPLICATION_TEMP`
		_EnvLogConf = `APPLICATION_LOG_CONFIGURATION`
		_EnvLogPath = `APPLICATION_LOG_PATH`
		_EnvState   = `APPLICATION_STATE`
		_EnvSocket  = `APPLICATION_SOCKET`
	)
	var tmp string
	if tmp = os.Getenv(_EnvName); tmp != "" {
		cnf.appConfiguration.ApplicationName = tmp
	}
	if tmp = os.Getenv(_EnvHome); tmp != "" {
		cnf.appConfiguration.WorkingDirectory = tmp
	}
	if tmp = os.Getenv(_EnvPid); tmp != "" {
		cnf.appConfiguration.PidFile = tmp
	}
	if tmp = os.Getenv(_EnvTemp); tmp != "" {
		cnf.appConfiguration.TempPath = tmp
	}
	if tmp = os.Getenv(_EnvLogConf); tmp != "" {
		cnf.appConfiguration.LogConfiguration = tmp
	}
	if tmp = os.Getenv(_EnvLogPath); tmp != "" {
		cnf.appConfiguration.LogPath = tmp
	}
	if tmp = os.Getenv(_EnvState); tmp != "" {
		cnf.appConfiguration.StateFile = tmp
	}
	if tmp = os.Getenv(_EnvSocket); tmp != "" {
		cnf.appConfiguration.SocketFile = tmp
	}
	return
}
