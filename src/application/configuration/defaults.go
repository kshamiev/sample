package configuration // import "application/configuration"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

// CleanAllPath Приведение всех путей в конфиге к абсолютному виду
// Создание папок
func (cnf *impl) CleanAllPath() {
	var i int
	var err error
	var makeDirectory []string

	// Root var
	cnf.appConfiguration.WorkingDirectory = cnf.AbsPath(cnf.appConfiguration.WorkingDirectory)
	cnf.appConfiguration.PidFile = cnf.AbsPath(cnf.appConfiguration.PidFile)
	makeDirectory = append(makeDirectory, cnf.AbsPath(path.Dir(cnf.appConfiguration.PidFile)))
	cnf.appConfiguration.TempPath = cnf.AbsPath(cnf.appConfiguration.TempPath)
	makeDirectory = append(makeDirectory, cnf.appConfiguration.TempPath)
	cnf.appConfiguration.CachePath = cnf.AbsPath(cnf.appConfiguration.CachePath)
	makeDirectory = append(makeDirectory, cnf.appConfiguration.CachePath)
	// State file
	makeDirectory = append(makeDirectory, cnf.AbsPath(path.Dir(cnf.appConfiguration.StateFile)))
	cnf.appConfiguration.StateFile = cnf.AbsPath(cnf.appConfiguration.StateFile)
	// Socket file
	makeDirectory = append(makeDirectory, cnf.AbsPath(path.Dir(cnf.appConfiguration.SocketFile)))
	cnf.appConfiguration.SocketFile = cnf.AbsPath(cnf.appConfiguration.SocketFile)

	// Loging
	if cnf.appConfiguration.LogPath == "" {
		cnf.appConfiguration.LogPath = "log/"
	}
	cnf.appConfiguration.LogPath = cnf.AbsPath(cnf.appConfiguration.LogPath)
	makeDirectory = append(makeDirectory, cnf.appConfiguration.LogPath)
	if cnf.appConfiguration.LogConfiguration != "" {
		cnf.appConfiguration.LogConfiguration = cnf.AbsPath(cnf.appConfiguration.LogConfiguration)
	}

	//// Keys var
	//cnf.appConfiguration.Keys.JWT.PrivateKeyFile = cnf.AbsPath(cnf.appConfiguration.Keys.JWT.PrivateKeyFile)
	//cnf.appConfiguration.Keys.JWT.PublicKeyFile = cnf.AbsPath(cnf.appConfiguration.Keys.JWT.PublicKeyFile)
	//makeDirectory = append(makeDirectory, path.Dir(cnf.appConfiguration.Keys.JWT.PrivateKeyFile))
	//makeDirectory = append(makeDirectory, path.Dir(cnf.appConfiguration.Keys.JWT.PublicKeyFile))

	// Database
	if cnf.appConfiguration.Database.Migrations != "" {
		cnf.appConfiguration.Database.Migrations = cnf.AbsPath(cnf.appConfiguration.Database.Migrations)
		makeDirectory = append(makeDirectory, cnf.appConfiguration.Database.Migrations)
	}

	//// SMTP
	//if cnf.appConfiguration.Smtp.Templates != "" {
	//	cnf.appConfiguration.Smtp.Templates = cnf.AbsPath(cnf.appConfiguration.Smtp.Templates)
	//	makeDirectory = append(makeDirectory, cnf.appConfiguration.Smtp.Templates)
	//}

	// WEB Server var
	for i = range cnf.appConfiguration.WEBServers {
		// Address
		if cnf.appConfiguration.WEBServers[i].Server.Address != "" {
			cnf.appConfiguration.WEBServers[i].Server.Address = strings.TrimRight(cnf.appConfiguration.WEBServers[i].Server.Address, "/")
		}
		// Paths
		cnf.appConfiguration.WEBServers[i].Server.Socket = cnf.AbsPath(cnf.appConfiguration.WEBServers[i].Server.Socket)
		cnf.appConfiguration.WEBServers[i].DocumentRoot = cnf.AbsPath(cnf.appConfiguration.WEBServers[i].DocumentRoot)
		cnf.appConfiguration.WEBServers[i].Pages = cnf.AbsPath(cnf.appConfiguration.WEBServers[i].Pages)
		makeDirectory = append(makeDirectory, cnf.appConfiguration.WEBServers[i].Pages)
		makeDirectory = append(makeDirectory, cnf.appConfiguration.WEBServers[i].DocumentRoot)
	}

	// Storage
	// Хранилище файлов
	if cnf.appConfiguration.Storage != "" {
		cnf.appConfiguration.Storage = rexSlashAtEnd.ReplaceAllString(cnf.appConfiguration.Storage, ``)
		cnf.appConfiguration.Storage = cnf.AbsPath(cnf.appConfiguration.Storage)
		makeDirectory = append(makeDirectory, cnf.appConfiguration.Storage)
	}

	// Создание папок
	for i := range makeDirectory {
		if makeDirectory[i] == "" {
			continue
		}
		if err = os.MkdirAll(makeDirectory[i], os.FileMode(0770)); err != nil && err != os.ErrExist {
			log.Printf("Error create folder '%s': %s", makeDirectory[i], err.Error())
		}
	}
}

// MakeDefaults Set configuration default value
func (cnf *impl) MakeDefaults() {
	var tmp = strings.Split(os.Args[0], string(os.PathSeparator))

	// Root var
	if cnf.appConfiguration.ApplicationName == "" {
		if len(tmp) > 0 {
			cnf.appConfiguration.ApplicationName = tmp[0]
		}
	}
	if cnf.appConfiguration.WorkingDirectory == "" {
		cnf.appConfiguration.WorkingDirectory = `~`
	}
	if cnf.appConfiguration.TempPath == "" {
		cnf.appConfiguration.TempPath = os.TempDir()
	}
	if cnf.appConfiguration.CachePath == "" {
		if len(tmp) > 0 {
			cnf.appConfiguration.CachePath = fmt.Sprintf("/var/cache/%s", tmp[0])
		}
	}
	// State var
	if cnf.appConfiguration.StateFile == "" {
		if len(tmp) > 0 {
			cnf.appConfiguration.StateFile = fmt.Sprintf("/var/spool/%s/%s.state", tmp[0], tmp[0])
		}
	}
	if cnf.appConfiguration.SocketFile == "" {
		if len(tmp) > 0 {
			cnf.appConfiguration.SocketFile = fmt.Sprintf("/var/run/%s.sock", tmp[0])
		}
	}
	cnf.makeDefaultsDatabase()
	cnf.makeDefaultsWebServer()

	// Storage
	// Хранилище файлов
	if cnf.appConfiguration.Storage == "" {
		if len(tmp) > 0 {
			cnf.appConfiguration.Storage = fmt.Sprintf("/var/lib/%s", tmp[0])
		}
	}
}

func (cnf *impl) makeDefaultsDatabase() {
	const (
		socket = `socket`
		tcp    = `tcp`
	)

	// Database var
	if cnf.appConfiguration.Database.Host == "" {
		cnf.appConfiguration.Database.Host = "localhost"
	}
	if cnf.appConfiguration.Database.Driver == "" {
		cnf.appConfiguration.Database.Driver = `mysql`
	}
	if cnf.appConfiguration.Database.Port == 0 {
		cnf.appConfiguration.Database.Port = 3306
	}
	switch strings.ToLower(cnf.appConfiguration.Database.Type) {
	case socket:
		cnf.appConfiguration.Database.Type = socket
	default:
		cnf.appConfiguration.Database.Type = tcp
	}
	if cnf.appConfiguration.Database.Name == "" {
		cnf.appConfiguration.Database.Name = `database`
	}
	if cnf.appConfiguration.Database.Login == "" {
		cnf.appConfiguration.Database.Login = `root`
	}
	if cnf.appConfiguration.Database.Charset == "" {
		cnf.appConfiguration.Database.Charset = `utf8`
	}
}

func (cnf *impl) makeDefaultsWebServer() {
	const (
		socket = `socket`
		tcp    = `tcp`
	)
	var i int

	// WEB Server var
	for i = range cnf.appConfiguration.WEBServers {
		if cnf.appConfiguration.WEBServers[i].Server.Host == "" {
			cnf.appConfiguration.WEBServers[i].Server.Host = `0.0.0.0`
		}
		if cnf.appConfiguration.WEBServers[i].Server.Port == 0 {
			cnf.appConfiguration.WEBServers[i].Server.Port = 80
		}
		switch strings.ToLower(cnf.appConfiguration.WEBServers[i].Server.Mode) {
		case socket:
			cnf.appConfiguration.WEBServers[i].Server.Mode = socket
		default:
			cnf.appConfiguration.WEBServers[i].Server.Mode = tcp
		}
	}
}
