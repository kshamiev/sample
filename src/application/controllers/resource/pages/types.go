package pages // import "application/controllers/resource/pages"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"regexp"

	"application/controllers/resource/pool"
	"application/models/file"
	"application/models/filecache"
	modelsPages "application/models/pages"

	"gopkg.in/webnice/web.v1/route"
)

const (
	keyIndexHTML              = `index.html`                    // Имя файла в DocumentRoot для главной страницы
	keyRouteTxt               = `/route.txt`                    // Имя файла в DocumentRoot для описания дополнительного роутинга для SPA/PWA
	ifModifiedSinceTimeFormat = `Mon, 02 Jan 2006 15:04:05 GMT` // Формат даты и времени для заголовка IfModifiedSince
)

var (
	// URN адреса главной страницы ведущие к файлу keyIndexHTML или к Pages() URN '/'
	// К данным адресам добавятся адреса указанные в файле route.txt
	indexURN = []string{
		`/`,
		`/index.htm`,
		`/index.html`,
	}
	excludeStaticRootFiles = []string{
		`icons.tpl.html`,
		`icons.tpl.png`,
		`icons.tpl.svg`,
		`index.html`,
		`robots.tpl.txt`,
		`route.txt`,
	}

	rexRouteTxtLine = regexp.MustCompile(`(?m)^(\s*)(/.+)$`)
	rexRootFiles    = regexp.MustCompile(`(?msi)^([^/]+)$`)
)

// Interface is an interface of package
type Interface interface {
	// Debug Set debug mode
	Debug(d bool) Interface

	// DocumentRoot Устанавливает путь к корню веб сервера
	DocumentRoot(path string) Interface

	// TemplatePages Устанавливает путь к папке файлов шаблонов
	TemplatePages(path string) Interface

	// ServerURL Устанавливает основной адрес веб сервера
	ServerURL(url string) Interface

	// Установка роутинга к статическим файлам
	SetRouting(rou route.Interface) Interface
}

// impl is an implementation of package
type impl struct {
	debug         bool                  // =true - debug mode is on
	rootPath      string                // Путь к корню веб сервера
	templatePages string                // Путь к файлам шаблонов
	serverURL     string                // Основной URL сервера
	serverScheme  string                // Выделенный из URL протокол сервера
	serverDomain  string                // Выделенный из URL домен сервера
	indexURN      []string              // Все URN зарегистрированные для выдачи index страницы
	Mfi           file.Interface        // Интерфейс работы с файлами
	Mfc           filecache.Interface   // Интерфейс файлового кеша в памяти
	Pool          pool.Interface        // Интерфейс пула переменных для переиспользования памяти
	Pgm           modelsPages.Interface // Интерфейс работы с шаблонами страниц
}
