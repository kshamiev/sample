package robots // import "application/controllers/resource/robots"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/controllers/resource/pool"
	"application/controllers/resource/sitemap"
	"application/models/filecache"

	"gopkg.in/webnice/web.v1/route"
)

const (
	keyRobots                 = `/robots.txt`                   // URL robots.txt
	keyRobotsTemplate         = `robots.tpl.txt`                // Шаблон robots.txt
	ifModifiedSinceTimeFormat = `Mon, 02 Jan 2006 15:04:05 GMT` // Формат даты и времени для заголовка IfModifiedSince
)

// Interface is an interface of package
type Interface interface {
	// Debug Set debug mode
	Debug(d bool) Interface

	// DocumentRoot Устанавливает путь к корню веб сервера
	DocumentRoot(path string) Interface

	// ServerURL Устанавливает основной адрес веб сервера
	ServerURL(url string) Interface

	// Sitemap Устанавливает интерфейс sitemap
	Sitemap(smi sitemap.Interface) Interface

	// Установка роутинга к статическим файлам
	SetRouting(rou route.Interface) Interface
}

// impl is an implementation of package
type impl struct {
	debug        bool                // =true - debug mode is on
	rootPath     string              // Путь к корню веб сервера
	serverURL    string              // Основной URL сервера
	serverScheme string              // Выделенный из URL протокол сервера
	serverDomain string              // Выделенный из URL домен сервера
	Mfc          filecache.Interface // Интерфейс файлового кеша в памяти
	Pool         pool.Interface      // Интерфейс пула переменных для переиспользования памяти
	Smi          sitemap.Interface   // Интерфейс sitemap
}

// Переменные шаблонизатора
type templateVars struct {
	RequestScheme string      // Протокол запроса (http или https)
	RequestDomain string      // Домен запроса, может быть разным если сервер размещён на нескольких доменах
	ServerURL     string      // URL сервера описанный в конфигурации веб сервера
	ServerScheme  string      // Протокол сервера, описанный в конфигурации веб сервера (http или https)
	ServerDomain  string      // Домен сервера, описанный в конфигурации веб сервера
	Sitemap       []tvSitemap // Блок sitemap.xml
}

// Переменные шаблонизатора sitemap блок
type tvSitemap struct {
	URN string // URL адрес sitemap.xml или sitemap-index.xml
}
