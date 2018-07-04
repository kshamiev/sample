package configuration // import "application/configuration"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"application/configuration/paths"
)

// SeekConfigFile Seek configuration file
func (cnf *impl) SeekConfigFile(configPaths []string) (configPath string, err error) {
	var i int
	for i = range configPaths {
		if _, err = os.Stat(configPaths[i]); !os.IsNotExist(err) {
			configPath = configPaths[i]
			break
		}

	}
	if configPath != "" {
		configPath, err = filepath.Abs(configPath)
	} else {
		var e = `Can't find configuration file:`
		for i = range configPaths {
			e += fmt.Sprintf("\n - '%s'", configPaths[i])
		}
		err = fmt.Errorf("%s", e)
	}
	return
}

// InitPaths Formulating the ways to find the configuration file (os dependent)
func (cnf *impl) InitPaths() (configPaths []string) {
	if cnf.args.FilePath != "" {
		configPaths = append(configPaths, cnf.args.FilePath)
		return
	}
	// Default configuration file search paths
	configPaths = append(configPaths, paths.SearchСonfigurationPaths()...)
	return
}

// AbsPath Преобразование пути к файлу/папке в абсолютные
// Замена символа '~' на путь к домашней директории пользователя
func (cnf *impl) AbsPath(str string) (ret string) {
	if len(str) > 0 {
		switch str[0] {
		case '~':
			ret = path.Join(cnf.user.HomeDir, str[1:])
		case '/':
			ret = str
		default:
			ret = path.Join(cnf.appConfiguration.WorkingDirectory, str)
		}
	}
	return
}
