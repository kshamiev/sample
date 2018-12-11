package fileinfo // import "application/models/fileinfo"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"os"
	"time"
)

// Interface is an interface of package
type Interface interface {
	os.FileInfo

	// CopyFrom Копирование значений из интерфейса переданного объекта и присваивание собственному объекту
	CopyFrom(src os.FileInfo) Interface

	// SetName Установка нового значения для Name()
	SetName(name string) Interface

	// SetSize Установка нового значения для Size()
	SetSize(size int64) Interface

	// SetMode Установка нового значения для Mode()
	SetMode(mode os.FileMode) Interface

	// SetModTime Установка нового значения для ModTime()
	SetModTime(modTime time.Time) Interface

	// SetIsDir Установка нового значения для IsDir()
	SetIsDir(isDir bool) Interface

	// SetSys Установка нового значения для Sys()
	SetSys(sys interface{}) Interface
}

// impl is an implementation of package
type impl struct { // nolint: maligned
	name    string      // base name of the file
	size    int64       // length in bytes for regular files; system-dependent for others
	mode    os.FileMode // file mode bits
	modTime time.Time   // modification time
	isDir   bool        // abbreviation for Mode().IsDir()
	sys     interface{} // underlying data source (can return nil)
}
