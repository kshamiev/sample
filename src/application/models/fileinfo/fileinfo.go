package fileinfo // import "application/models/fileinfo"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"os"
	"path"
	"time"
)

// New creates a new object and return interface
func New() Interface {
	var fif = new(impl)
	return fif
}

// Name The base name of the file
func (fif *impl) Name() string { return fif.name }

// Size Length in bytes for regular files; system-dependent for others
func (fif *impl) Size() int64 { return fif.size }

// Mode file mode bits
func (fif *impl) Mode() os.FileMode { return fif.mode }

// ModTime modification time
func (fif *impl) ModTime() time.Time { return fif.modTime }

// IsDir abbreviation for Mode().IsDir()
func (fif *impl) IsDir() bool { return fif.isDir }

// Sys underlying data source (can return nil)
func (fif *impl) Sys() interface{} { return fif.sys }

// CopyFrom Копирование значений из интерфейса переданного объекта и присваивание собственному объекту
func (fif *impl) CopyFrom(src os.FileInfo) Interface {
	fif.SetName(src.Name()).
		SetSize(src.Size()).
		SetMode(src.Mode()).
		SetModTime(src.ModTime()).
		SetIsDir(src.IsDir()).
		SetSys(src.Sys())
	return fif
}

// SetName Установка нового значения для Name()
func (fif *impl) SetName(name string) Interface { fif.name = path.Base(name); return fif }

// SetSize Установка нового значения для Size()
func (fif *impl) SetSize(size int64) Interface { fif.size = size; return fif }

// SetMode Установка нового значения для Mode()
func (fif *impl) SetMode(mode os.FileMode) Interface { fif.mode = mode; return fif }

// SetModTime Установка нового значения для ModTime()
func (fif *impl) SetModTime(modTime time.Time) Interface { fif.modTime = modTime; return fif }

// SetIsDir Установка нового значения для IsDir()
func (fif *impl) SetIsDir(isDir bool) Interface { fif.isDir = isDir; return fif }

// SetSys Установка нового значения для Sys()
func (fif *impl) SetSys(sys interface{}) Interface { fif.sys = sys; return fif }
