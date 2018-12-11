package sitemap // import "application/models/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"encoding/xml"
	"io"
	"math"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	keyDb           = `sitemap.bolt.db` // Название файла базы данных для хранения sitemap записей
	keyBucketURLSet = `urlset`          // Название ведра для URL
	keyBucketSystem = `system`          // Название ведра для системных ключей кеша
	keyURLSetAt     = `URLSetAt`        // Хранение даты и времени последнего изменения ведра URL
	keyURLCount     = `URLCount`        // Хранение количества записей key/value хранилища
	//keySlash        = `/`               // URN path separator

	// XMLns Указывает стандарт текущего протокола
	XMLns = `http://www.sitemaps.org/schemas/sitemap/0.9`

	// DefaultPriority Приоритет URI по умолчанию
	DefaultPriority = float32(0.5)

	// MaxRecords Максимальное число записей в sitemap и sitemap-index
	MaxRecords = uint64(50000)

	// MaxIndex Максимальный индекс файла sitemap и sitemap-index
	MaxIndex = (math.MaxUint64 - 1) / MaxRecords
)

// Interface is an interface of package
type Interface interface {
	io.Closer

	// Debug Set debug mode
	Debug(d bool) Interface

	// Add Добавление URN в sitemap
	Add(items []*Location) (Interface, error)

	// Del Удаление URN из sitemap
	Del(urn string) (Interface, error)

	// Count Возвращает количество записей URN в sitemap
	Count() (ret uint64)

	// GetByRange Выгружает URN в соответствии с порядковыми номерами ключей
	// Функция вернёт данные начиная с [from] элемента, массив не более чем [count] элементов
	GetByRange(from uint64, count uint64) (ret []*Location, err error)

	// SitemapXMLWriteTo Возвращает sitemap.xml, с данными начиная с позиции from в количестве count
	SitemapXMLWriteTo(wr io.Writer, url string, from uint64, count uint64) (err error)

	// SitemapXML Возвращает sitemap.xml с данными начиная с позиции from в количестве count
	SitemapXML(url string, from uint64, count uint64) (ret []byte, err error)

	// SitemapIndexXMLWriteTo Возвращает sitemap-index.xml
	SitemapIndexXMLWriteTo(wr io.Writer, url string, idx uint64) (err error)

	// SitemapIndexXML Возвращает sitemap-index.xml
	SitemapIndexXML(url string, idx uint64) (ret []byte, err error)

	// Errors Ошибки известного состояни, которые могут вернуть функции пакета
	Errors() *Error

	// Reset all data and create new database
	Reset() (Interface, error)

	// URLSetAt Функция возвращает дату и время последнего изменения в списке URL
	URLSetAt() time.Time
}

// impl is an implementation of package
type impl struct {
	debug  bool        // =true - debug mode is on
	db     *bolt.DB    // База данных
	dbSync *sync.Mutex // Защита от гонки
	dbPath string      // Путь к базе данных bolt
}

// Sitemap is a structure of <sitemap>
// https://www.sitemaps.org/ru/protocol.html
type Sitemap struct {
	XMLName xml.Name   `xml:"urlset"`     // Инкапсулирует этот файл и указывает стандарт текущего протокола
	XMLNS   string     `xml:"xmlns,attr"` // Адрес описания стандарта текущего протокола
	URL     []*URLPart `xml:"url"`        // Массив URL страниц
}

// URLPart is a structure of <url> in <sitemap>
type URLPart struct {
	Loc        string     `xml:"loc"`        // URL-адрес страницы. Длина этого значения не должна превышать 2048 символов
	LastMod    time.Time  `xml:"lastmod"`    // Дата последнего изменения ресурса в формате ISO 8601
	ChangeFreq ChangeFreq `xml:"changefreq"` // Вероятная частота изменения этой страницы
	Priority   float32    `xml:"priority"`   // Приоритетность URL ресурса относительно других URL в пределах сайта. От 0.0 до 1.0. Приоритет страницы по умолчанию: 0.5
}

// Index is a structure of <sitemap-index>
// https://www.sitemaps.org/ru/protocol.html
type Index struct {
	XMLName xml.Name     `xml:"sitemapindex"` // Инкапсулирует информацию о всех файлах Sitemap в этом файле
	XMLNS   string       `xml:"xmlns,attr"`   // Адрес описания стандарта текущего протокола
	Sitemap []*IndexPart `xml:"sitemap"`      // Инкапсулирует информацию об отдельно взятых файлах Sitemap
}

// IndexPart is a partition structure of <sitemap-index>
type IndexPart struct {
	Loc     string    `xml:"loc"`     // Указывает местоположение файла Sitemap. Этим местоположением может быть файл Sitemap, файл Atom, файл RSS или простой текстовый файл
	LastMod time.Time `xml:"lastmod"` // Указывает время изменения соответствующего файла Sitemap в формате ISO 8601
}

// Location Структура описывающая URL адрес ресурса
type Location struct {
	URN      string     `msgpack:"urn" json:"urn"` // Адрес ресурса в пределах сервера. Не должна превышать 2048 символов
	ModTime  time.Time  `msgpack:"mod" json:"mod"` // Дата и время последнего изменения ресурса
	Change   ChangeFreq `msgpack:"chf" json:"chf"` // Вероятная частота изменения ресурса
	Priority float32    `msgpack:"pri" json:"pri"` // Приоритетность URL ресурса относительно других URL в пределах сайта. От 0.0 до 1.0
}
