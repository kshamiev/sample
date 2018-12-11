package sitemap // import "application/controllers/resource/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"regexp"

	"application/controllers/resource/pool"
	"application/models/filecache"
	modelsSitemap "application/models/sitemap"

	"gopkg.in/webnice/web.v1/route"
)

const (
	urnPattern                = `/sitemap:suffix`               // Шаблон роутинга
	ifModifiedSinceTimeFormat = `Mon, 02 Jan 2006 15:04:05 GMT` // Формат даты и времени для заголовка IfModifiedSince
)

var (
	rexSitemap   = regexp.MustCompile(`(?mis)^/(sitemap|sitemap-index)-*(\d*)-*(\d*).xml$`) // Regexp разбора URN запроса
	rexURLCheck  = regexp.MustCompile(`(?mis)^(http|https)://`)                             // Проверка на URL: адрес
	rexSlashLast = regexp.MustCompile(`/*$`)                                                // Слэш в конце
	//rexSlashFirst = regexp.MustCompile(`^/*`)                                                // Слэш в начале
)

// Interface is an interface of package
type Interface interface {
	io.Closer

	// Debug Set debug mode
	Debug(d bool) Interface

	// DocumentRoot Устанавливает путь к корню веб сервера
	DocumentRoot(path string) Interface

	// ServerURL Устанавливает основной адрес веб сервера
	ServerURL(url string) (Interface, error)

	// Установка роутинга к статическим файлам
	SetRouting(rou route.Interface) Interface

	// Add Добавление URI ресурса в sitemap.xml
	Add(items []*modelsSitemap.Location) (err error)

	// Del Удаление URI ресурса из sitemap.xml
	Del(urn string) (err error)

	// Count Количество занисей в sitemap
	Count() uint64

	// Links Получение текущих линков sitemap.xml или sitemap-index.xml
	Links() []string

	// Errors Ошибки известного состояни, которые могут вернуть функции пакета
	Errors() *Error
}

// impl is an implementation of package
type impl struct {
	debug     bool                    // =true - debug mode is on
	rootPath  string                  // Путь к корню веб сервера
	serverURL string                  // Основной URL сервера
	Mfc       filecache.Interface     // Интерфейс файлового кеша в памяти
	Pool      pool.Interface          // Интерфейс пула переменных для переиспользования памяти
	Smm       modelsSitemap.Interface // Интерфейс модели sitemap
}
