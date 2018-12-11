package types // import "application/models/pages/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	//"io"
	//"bytes"
	//
	//"gopkg.in/webnice/web.v1/method"
)

// RequestResponse Интерфейс запроса и ответа контроллера
type RequestResponse interface {
	// Request HTTP request received by a server or to be sent by a client
	Request() *http.Request

	// ResponseWriter interface is used by an HTTP handler to construct an HTTP response
	ResponseWriter() http.ResponseWriter
}

// Responder Интерфейс к объекту - обработчику Response
type Responder interface {
}

// Templater Интерфейс доступа к объекту сгруппированных по URN шаблонов
type Templater interface {
}

// Response Объект обработчик запроса
type Response struct {
	// Handlers Обработчики запроса
	Handlers []Handler

	// TemplatesURN Список URN, шаблоны которых требуются контроллеру для работы
	TemplatesURN []string
}

// Body Информация и содержимое файла шаблона
type Body struct {
	FullName   string // FullName   Полный путь и имя файла
	Name       string // Name       Имя файла без пути
	NameBasis  string // NameBasis  Имя файла без расширения
	Path       string // Path       Относительный путь к файлу он же URN
	IsInclude  bool   // IsInclude  =true если шаблон является *.inc.html
	IsTemplate bool   // IsTemplate =true если шаблон является index.tpl.html или layout.tpl.html
	Source     string // Source Если файл является ссылкой, source указывает на файл контента
}

// Templates Объект сгруппированных по URN шаблонов
type Templates []*Body

// Template Интерфейс загруженных в память шаблонов
//type Template interface {
//	// Len Количество шаблонов содержащихся в объекте
//	Len() int
//
//	// List Список URN всех шаблонов
//	List() []string
//
//	// First Первый URN в списке, как правило это URN к которому пришёл запрос
//	First() string
//
//	// HasURN (URN) True если существует шаблон с указанным URN
//	HasURN(urn string) bool
//
//	// Map (urn) Данные шаблона указанного URN, представленные в виде map.
//	// Ключём является имя файла без расширения. Например для файла layout.tpl.html ключ будет layout
//	Map(urn string) map[string]*Body
//
//	// Index (urn) Возвращает главный шаблон для указанного URN
//	// если шаблон является многофайловым (папка), то ищется шаблон с именем index,
//	// если шаблон является однофайловым, то возвращается единственный существующий шаблон
//	// если шаблона нет, то возвращается nil
//	Index(string) *Body
//
//	// Layout Возвращает шаблон-макет (layout) для указанного URN, если шаблона нет, то возвращается nil
//	Layout(urn string) *Body
//
//	// KeyBody (urn, key) Возвращает структуру тела шаблона с указанным URN и ключём, если шаблона нет или он пустой то вернётся nil
//	KeyBody(string, string) *Body
//
//	// KeyData (urn, key) Возвращает данные шаблона с указанным URN и ключём, если шаблона нет или он пустой то вернётся nil
//	KeyData(string, string) *bytes.Buffer
//
//	// Keys (urn) Получение списка всех ключей шаблонов для указанного URN
//	Keys(string) []string
//
//	// HasKey (urn, key) True если для указанного URN существует шаблон с указанным ключём
//	HasKey(urn string, key string) bool
//
//	// Data (urn) Все данные шаблонов отсортированные в порядке: [все inc, layout, index] для указанного URN
//	// Шаблоны с пустыми данными пропускаются
//	Data(urn string) []*bytes.Buffer
//
//	// Reload Метод перечитывает все шаблоны объекта и возвращает интерфейс
//	Reload() Template
//}

// Map Карта всех URN совместно с контроллерами
//type Map struct {
//	// URN адрес
//	URN string
//
//	// Method Метод HTTP запроса
//	Method method.Method
//
//	// Template Шаблоны адреса в порядке запрошенном контроллером
//	Template *Template
//
//	// HandlerFunc Контроллеры адреса в порядке вызова
//	HandlerFunc []*http.HandlerFunc
//}
