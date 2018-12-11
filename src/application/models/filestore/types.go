package filestore // import "application/models/filestore"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"regexp"

	filestoreTypes "application/models/filestore/types"

	"gopkg.in/webnice/kit.v1/modules/db"
)

const (
	pathStorageTemporary = `tmp`   // Относительный путь к временным файлам
	pathStoragePermanent = `files` // Относительный путь к постоянным файлам
)

var (
	// Путь по умолчанию к папке хранения контента
	defaultStoragePath = `storage`

	// Выделение из имени файла его расширения
	rexFileExt = regexp.MustCompile(`\.([^.]+)$`)
)

// Interface is an interface of package
type Interface interface {
	// TEMPORARY STORAGE

	// NewTemporaryFile Создание нового временного файла
	NewTemporaryFile(filename string, size uint64, contentType string, inpFh io.Reader) (id uint64, err error)

	// OpenTemporaryFile Открытие для чтения временного файла по его ID
	OpenTemporaryFile(fileID uint64) (info *filestoreTypes.FilesTemporary, fh io.ReadCloser, err error)

	// PERMANENT STORAGE

	// NewPermanentFileFromTemporaryFile Создание постоянного файла из временного путём копирования
	NewPermanentFileFromTemporaryFile(fileID uint64) (ret *filestoreTypes.Filestore, err error)

	// PermanentFileInfo Загрузка информации о файле по ID
	PermanentFileInfo(fileID uint64) (info *filestoreTypes.Filestore, err error)

	// CleanOldData Очистка filestore от устаревших данных
	CleanOldData() error

	// ERRORS

	// ErrNotFound Not found
	ErrNotFound() error
}

// impl is an implementation of package
type impl struct {
	db.Implementation        // Наследование интерфейса работы с БД
	storagePath       string // Путь к хранилищу файлов
}
