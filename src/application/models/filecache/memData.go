package filecache // import "application/models/filecache"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"io"
	"os"
	"time"
)

// ContentType Возвращается Content-Type данных объекта
func (md *memData) ContentType() string { return md.contentType }

// CreateAt Возвращает дату и время создания объекта в памяти
func (md *memData) CreateAt() time.Time { return md.createAt }

// FullName Полное наименование объекта
func (md *memData) FullName() string { return md.fullName }

// Lifetime Возвращает время жизни объекта в памяти
// =  0 - без ограничения времени
// = -1 - перечитывание файла с диска или пересоздание виртуального файла при каждом обращении
func (md *memData) Lifetime() time.Duration { return md.lifetime }

// LifetimeSet Устанавливает время жизни объекта в памяти
// =  0 - без ограничения времени
// = -1 - перечитывание файла с диска или пересоздание виртуального файла при каждом обращении
func (md *memData) LifetimeSet(t time.Duration) Data { md.lifetime = t; return md }

// Length Размер объекта данных в памяти
func (md *memData) Length() uint64 { return uint64(len(md.body)) }

// WriteTo Is an interface implementation of io.WriteTo
func (md *memData) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, bytes.NewBuffer(md.body))
}

// Name Is an interface implementation of os.FileInfo
func (md *memData) Name() string { return md.info.Name() }

// Size Is an interface implementation of os.FileInfo
func (md *memData) Size() int64 { return md.info.Size() }

// Mode Is an interface implementation of os.FileInfo
func (md *memData) Mode() os.FileMode { return md.info.Mode() }

// ModTime Is an interface implementation of os.FileInfo
func (md *memData) ModTime() time.Time { return md.info.ModTime() }

// IsDir Is an interface implementation of os.FileInfo
func (md *memData) IsDir() bool { return md.info.IsDir() }

// Sys Is an interface implementation of os.FileInfo
func (md *memData) Sys() interface{} { return md.info.Sys() }

// Reader Создаёт и возвращает интерфейс независимого объекта io.Reader
func (md *memData) Reader() DataReader { return bytes.NewReader(md.body) }
