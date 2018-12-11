package file // import "application/models/file"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"bytes"
	"io"
	"os"
)

// Interface is an interface
type Interface interface {
	// CleanEmptyFolder Удаление пустых папок
	CleanEmptyFolder(pt string) error

	// Copy Копирует один файл в другой
	Copy(dst string, src string) (size int64, err error)

	// CopyWithSha512Sum Копирование контента с параллельным вычислением контрольной суммы алгоритмом SHA512
	CopyWithSha512Sum(dst io.Writer, src io.Reader) (written int64, sha512sum string, err error)

	// GetInfoSha512 Считывание информации о файле с контрольной суммой
	GetInfoSha512(filename string) (inf *InfoSha512, err error)

	// RecursiveFileList Поиск всех файлов начиная от path рекурсивно
	RecursiveFileList(path string) (ret []string, err error)

	// GetFileName Выделение из полного пути и имени файла, имя файла
	GetFileName(fileName string) string

	// LoadFile Загрузка файла в память и возврат в виде *bytes.Buffer
	LoadFile(name string) (data *bytes.Buffer, info os.FileInfo, err error)
}

// is an implementation
type impl struct {
}

// InfoSha512 Структура возвращаемой информации о файле
type InfoSha512 struct {
	Name   string // Название файла
	Size   int64  // Размер файла
	Sha512 string // Контрольная сумма файла
}
