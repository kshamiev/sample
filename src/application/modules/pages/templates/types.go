package templates // include "application/modules/pages/templates"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/method"
import (
	"bytes"
	"container/list"
	"net/http"
	"regexp"
	"time"
)

const (
	_KeyIndex  = `index`
	_KeyLayout = `layout`
)

var (
	singleton    *impl
	rexExtension = regexp.MustCompile(`\.[inc|tpl]+\.html$`)
	rexIsInclude = regexp.MustCompile(`\.inc\.html$`)
	rexSourceLnk = regexp.MustCompile(`^[\t\n\f\r ]*<!--[\t\n\f\r ]+#source:[\t\n\f\r ]+([^(?:>)]+)[\t\n\f\r ]+-->`)
	rexSpaceLast = regexp.MustCompile(`[\t\n\f\r ]+$`)
	rexSpaceFrst = regexp.MustCompile(`^[\t\n\f\r ]+`)
	rexSlashLast = regexp.MustCompile(`/+$`)
	rexSlashFrst = regexp.MustCompile(`^/+`)
)

// Interface is an interface of module
type Interface interface {
	// Error Последняя ошибка
	Error() error

	// Init Инициализация шаблонизатора и всех подключенных контроллеров шаблонов
	Init() Interface

	// Routing Выгрузка настроек роутинга
	Routing() http.Handler

	// Template Интерфейс ко всем загруженным шаблонам файлов
	Template() Template
}

// impl is an implementation of monule
type impl struct {
	Controllers *list.List // (Controller) Список зарегистрированных контроллеров
	Err         error      // Последняя ошибка
	Tpls        *tpls      // Все загруженные в память шаблоны
	Urn         *list.List // (*UrnMap) Подготовленная для загрузки в роутинг URN с контроллерами
	PagesRoot   string     // Путь к папкам шаблонов
	Debug       bool       // =true - дебаг режим
}

// Controller interface
type Controller interface {
	// Init Метод инициализации контроллера
	// Вызывается из modules/pages/templates
	Init() ([]*Rider, error)
}

// Rider Требования контроллера и регистрационные данные для роутинга
type Rider struct {
	// Templates Список URN соответствующий файлам шаблонов, которые требует контроллер при каждом своём вызове
	Templates []string
	// Hundler Список описаний URN которые контроллер будет обрабатывать и соответствующие функции
	Hundler Handler
}

// Handler Структура подписки контроллера на URN
type Handler struct {
	// Urn Основной URN
	Urn string
	// Method Метод HTTP запроса
	Method method.Method
	// Func Функция-обработчик запроса
	Func http.HandlerFunc
}

// Template Интерфейс к tpls - загруженным в память шаблонам html страницы
type Template interface {
	// Len Количество шаблонов содержащихся в объекте
	Len() int

	// ListUrn Список URN всех шаблонов
	ListUrn() []string

	// FirstUrn Первый URN в списке, как правило это URN к которому пришел запрос
	FirstUrn() string

	// HasUrn (urn) True если существует шаблон с указанным URN
	HasUrn(string) bool

	// Map (urn) Данные шаблона указанного URN, представленные в виде map.
	// Ключём является имя файла без расширения. Например для файла layout.tpl.html ключ будет layout
	Map(string) map[string]*TplBody

	// Index (urn) Возвращает главный шаблон для указанного URN
	// если шаблон является многофайловым (папка), то ищется шаблон с именем index,
	// если шаблон является однофайловым, то возвращается единственный существующий шаблон
	// если шаблона нет, то возвращается nil
	Index(string) *TplBody

	// Layout Возвращает шаблон-макет (layout) для указанного URN, если шаблона нет, то возвращается nil
	Layout(urn string) *TplBody

	// KeyBody (urn, key) Возвращает структуру тела шаблона с указанным URN и ключём, если шаблона нет или он пустой то вернётся nil
	KeyBody(string, string) *TplBody

	// KeyData (urn, key) Возвращает данные шаблона с указанным URN и ключём, если шаблона нет или он пустой то вернётся nil
	KeyData(string, string) *bytes.Buffer

	// Keys (urn) Получение списка всех ключей шаблонов для указанного URN
	Keys(string) []string

	// HasKey (urn, key) True если для указанного URN существует шаблон с указанным ключём
	HasKey(string, string) bool

	// AllData (urn) Все данные шаблонов отсортированные в порядке: [все inc, layout, index] для указанного URN
	// Шаблоны с пустыми данными пропускаются
	AllData(urn string) []*bytes.Buffer

	// Reload Перечитываем все шаблоны объекта и возвращаем интерфейс
	Reload() Template
}

// Объект данных файла шаблона
type tpl struct {
	// UrnAddress Основной URN
	UrnAddress string
	// MapData Данные шаблона или шаблонов если URN соответствует нескольким файлам.
	// Ключём является имя файла без расширения. Например для файла layout.tpl.html ключ будет layout
	MapData map[string]*TplBody
}

// Объект реализующий интерфейс Template
type tpls struct {
	// Tpl Шаблоны
	Tpl []*tpl
}

// TplBody Информация и содержимое файла шаблона
type TplBody struct {
	// Полный путь и имя файла
	FullPath string
	// Path Путь к файлу URN
	Path string
	// Name Имя файла без пути
	Name string
	// NameBasis Имя файла без расширения
	NameBasis string
	// IsInclude =true если шаблон является *.inc.html
	IsInclude bool
	// IsTemplate =true если шаблон является index.tpl.html или layout.tpl.html
	IsTemplate bool
	// Size Размер файла в байтах
	Size int64
	// ModTime Дата и время модификации файла в файловой системе
	ModTime time.Time
	// Data Содержимое файла. Если =nil то файл не загружался
	Data *bytes.Buffer
}

// UrnMap Карта всех URN совместно с контроллерами
type UrnMap struct {
	// Urn адрес
	Urn string
	// Method Метод HTTP запроса
	Method method.Method
	// Tpls Шаблоны адреса в порядке запрошенном контроллером
	Tpls *tpls
	// Hndl Контроллеры адреса в порядке вызова
	Hndl []*Handler
}
