package filecache // import "application/models/filecache"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"os"
	"time"
)

// Структура хранимых в памяти объектов файлов
type memData struct {
	createAt    time.Time     // Дата и время создания объекта в памяти
	lifetime    time.Duration // Время нахождения объекта в памяти до удаления или пересоздания. =0 - без ограничения
	creatorFn   CreateFn      // Функция создания объекта, если не указана, то используется функция чтения с диска
	loaded      bool          // =true - тело объекта загружено, =false - необходимо выполнить функцию CreatorFn для загрузки тела
	fullName    string        // Полное наименование объекта
	contentType string        // Content-Type объекта в памяти
	body        []byte        // Тело объекта в памяти
	info        os.FileInfo   // Информация об объекте в памяти на момент его создания
}

// Data Интерфейс работы с данными объекта в памяти
type Data interface {
	io.WriterTo
	os.FileInfo

	// ContentType Возвращается Content-Type данных объекта в памяти
	ContentType() string

	// CreateAt Возвращает дату и время создания объекта в памяти
	CreateAt() time.Time

	// FullName Полное наименование объекта
	FullName() string

	// Lifetime Возвращает время жизни объекта в памяти
	// =  0 - без ограничения времени
	// = -1 - перечитывание файла с диска или пересоздание виртуального файла при каждом обращении
	Lifetime() time.Duration

	// LifetimeSet Устанавливает время жизни объекта в памяти
	// =  0 - без ограничения времени
	// = -1 - перечитывание файла с диска или пересоздание виртуального файла при каждом обращении
	LifetimeSet(t time.Duration) Data

	// Length Размер объекта данных в памяти
	Length() uint64

	// Reader Создаёт и возвращает интерфейс независимого объекта io.Reader
	Reader() DataReader
}

// DataReader Интерфейс независимого объекта для чтения данных разными способами
type DataReader interface {
	io.ReadSeeker
	io.ReaderAt
	io.ByteScanner
	io.RuneScanner
	io.WriterTo

	// Len returns the number of bytes of the unread portion of the slice
	Len() int

	// Reset resets the Reader to be reading from b
	Reset(b []byte)

	// Size returns the original length of the underlying byte slice.
	// Size is the number of bytes available for reading via ReadAt.
	// The returned value is always the same and is not affected by calls
	// to any other method
	Size() int64
}
