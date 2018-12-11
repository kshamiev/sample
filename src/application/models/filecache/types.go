package filecache // import "application/models/filecache"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"context"
	"os"
	"regexp"
	"sync"
	"time"
)

var (
	singleton        *impl                                   // Singleton object
	extractExtension = regexp.MustCompile(`(?i)^.+\.(.+)?$`) // Выделение расширения имени файла
)

// Interface is an interface of package
type Interface interface {
	// Load Функция загружает объект в память и предоставляет интерфейс доступа к объекту в режиме только для чтения
	// Если объекта в памяти нет, то он запрашивается у функции CreateFn (функция по умолчанию - чтение объекта с диска)
	Load(name string) (ret Data, err error)

	// LoadWithLifetime Функция загружает объект в память и предоставляет интерфейс доступа к объекту в режиме только для чтения
	// Если объекта в памяти нет, то он запрашивается у функции CreateFn (функция по умолчанию - чтение объекта с диска)
	LoadWithLifetime(name string, lifetime time.Duration) (ret Data, err error)

	// List Возвращает список объектов в памяти
	List() (ret []string)

	// Size Возвращает суммарный размер всех объектов в памяти
	Size() (ret uint64)

	// IsExist Проверка существования файла в кеше
	IsExist(name string) bool

	// Virtual Добавление виртуального файла в кеш, файл создаётся путём вызова функции createFn
	// После создания файла, он запоминается в памяти на время указанное в lifetime
	Virtual(name string, lifetime time.Duration, createFn CreateFn) (err error)

	// Debug Set debug mode
	Debug(d bool) Interface

	// Errors Ошибки известного состояни, которые могут вернуть функции пакета
	Errors() *Error
}

// impl is an implementation of package
type impl struct {
	debug bool               // =true - debug mode is on
	mlist *list.List         // Список *memData
	mhash *memHash           // hash быстрого доступа по имени файла к элементам list из объектов *memData
	wdctx context.Context    // Интерфейс контекста для watchdog
	wdcfn context.CancelFunc // Функция прерывания работы watchdog
}

// Хэш с защитой
type memHash struct {
	sync.RWMutex
	Map map[string]*list.Element
}

// MemoryObject Структура данных передаваемых на хранение в памяти
type MemoryObject struct {
	Lifetime    time.Duration // Время нахождения объекта в памяти до удаления или пересоздания. =0 - без ограничения
	FullName    string        // Полное наименование объекта
	Body        []byte        // Тело объекта в памяти
	ContentType string        // Content-Type объекта в памяти
	Info        os.FileInfo   // Информация об объекте в памяти на момент его создания
	CreatorFn   CreateFn      // Функция создания объекта, если не указана, то используется функция чтения с диска
}

// CreateFn Функция создания нового объекта для размещения в памяти
type CreateFn func(name string) (ret *MemoryObject, err error)
