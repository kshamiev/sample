package pages // import "application/modules/pages"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"io"
	"os"
)

// Загрузка файла в память
// TODO
// Запоминать в памяти все файлы на 1 час если не изменилась длинна и время модификации файла
func (pgs *impl) loadCached(rfn string) (ret *File, err error) {
	var fh *os.File
	var l int64

	ret = new(File)
	ret.Info, err = os.Stat(rfn)
	if os.IsNotExist(err) {
		err = pgs.ErrFileNotFound()
		return
	} else if err != nil {
		return
	}
	if ret.Info.IsDir() {
		err = pgs.ErrFileNotFound()
		return
	}
	if fh, err = os.Open(rfn); err != nil {
		return
	}
	defer func() { _ = fh.Close() }()

	ret.Body = &bytes.Buffer{}
	l, err = io.Copy(ret.Body, fh)
	if l != ret.Info.Size() {
		err = pgs.ErrIncorrectFileLength()
		return
	}

	return
}
