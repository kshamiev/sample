package assets // import "application/controllers/resource/assets"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/controllers/resource/pool"
	"application/models/filecache"

	"gopkg.in/webnice/web.v1/route"
)

const (
	preffixPath               = `/assets`                       // Префик статических файлов
	ifModifiedSinceTimeFormat = `Mon, 02 Jan 2006 15:04:05 GMT` // Формат даты и времени для заголовка IfModifiedSince
)

// Interface is an interface of package
type Interface interface {
	// Debug Set debug mode
	Debug(d bool) Interface

	// DocumentRoot Устанавливает путь к корню веб сервера
	DocumentRoot(path string) Interface

	// Установка роутинга к статическим файлам
	SetRouting(rou route.Interface) Interface
}

// impl is an implementation of package
type impl struct {
	debug    bool                // =true - debug mode is on
	rootPath string              // Путь к корню веб сервера
	Mfc      filecache.Interface // Интерфейс файлового кеша в памяти
	Pool     pool.Interface      // Интерфейс пула переменных для переиспользования памяти
}
