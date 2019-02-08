package filestore // import "application/models/filestore"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

// SetDefaultStoragePath Установка пути по умолчанию к хранилищу файлов
func SetDefaultStoragePath(storagePath string) {
	if storagePath == "" {
		return
	}
	defaultStoragePath = storagePath
}

// New creates a new object and return interface
func New() Interface {
	var ufm = &impl{
		storagePath: defaultStoragePath,
	}
	return ufm
}

// Errors Ошибки известного состояни, которые могут вернуть функции пакета
func (ufm *impl) Errors() *Error { return Errors() }
