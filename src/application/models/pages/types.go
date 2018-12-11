package pages // import "application/models/pages"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"regexp"
	"sync"

	"application/models/filecache"
	pagesTypes "application/models/pages/types"

	"gopkg.in/webnice/web.v1/route"
)

const (
	keyIndex  = `index`
	keyLayout = `layout`
)

var (
	singleton         *impl // Singleton object of package
	rexFilenameFilter = regexp.MustCompile(`(?mi)^(.*)\.[inc|tpl]+\.html$`)
	rexIsInclude      = regexp.MustCompile(`(?mi)\.inc\.html$`)
	rexSourceLnk      = regexp.MustCompile(`(?mi)^[\s]*<!--[\s]+#source:[\s]+([^(?:>)]+)[\s]+-->`)
)

// Interface is an interface of package
type Interface interface {
	// Debug Set debug mode
	Debug(d bool) Interface

	// TemplatePages Путь к корню файлов шаблонов
	TemplatePages(p string) Interface

	// ServerURL Устанавливает основной адрес веб сервера
	ServerURL(url string) Interface

	// SetRouting Настройка роутинга
	SetRouting(rou route.Interface) uint64

	// Init Инициализация объекта
	Init() (err error)

	// RegisterController Регистрация контроллера
	RegisterController(ctl pagesTypes.Controller) Interface

	// Template Интерфейс ко всем загруженным шаблонам файлов
	//Template() Template
}

// impl is an implementation of package
type impl struct {
	debug         bool                            // =true - debug mode is on
	templatePages string                          // Путь к файлам шаблонов
	serverURL     string                          // Основной URL сервера
	serverScheme  string                          // Выделенный из URL протокол сервера
	serverDomain  string                          // Выделенный из URL домен сервера
	controllers   *list.List                      // Список зарегистрированных контроллеров (pagesTypes.Controller)
	Mfc           filecache.Interface             // Интерфейс файлового кеша в памяти
	urn           map[string]pagesTypes.Responder // Интерфейс к карте всех обрабатываемых контроллерами URN и всех шаблонов
	urnSync       *sync.RWMutex                   // Защита urn map от конкурентного досутпа
	tpl           map[string]pagesTypes.Templater // Интерфейс всех шаблонов сгруппированных по URN
	tplSync       *sync.RWMutex                   // Защита tpl map от конкурентного досутпа
}
