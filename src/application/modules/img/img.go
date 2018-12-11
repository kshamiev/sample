// Package img Пакет помогающий работать с графическими озображениями
// nolint: goimports
package img // import "application/modules/img"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"os"

	"golang.org/x/image/tiff"
	"image/gif"
	"image/jpeg"
)

// New creates new object
func New() Interface {
	var img = new(impl)
	return img
}

// Errors Ошибки известного состояни, которые могут вернуть функции пакета
func (img *impl) Errors() *Error { return Errors() }

// New Создаёт пустой объект графического изображения
// Объект обладает интерфейсом io.WriteCloser который можно использовать для загрузки данных графического
// объекта, после вызова Close(), записанные во Writer данные обрабатываются в форматонезависимый графический образ
// и присваиваются объекту
func (img *impl) New() Image {
	var ret = &imgItem{
		wr:      &bytes.Buffer{},
		optTIFF: new(tiff.Options),
		optGIF:  new(gif.Options),
		optJPEG: new(jpeg.Options),
	}
	return ret
}

// Open Загрузка объекта графического изображения из файла
func (img *impl) Open(filename string) (ret Image, err error) {
	var fh *os.File
	var i *imgItem

	fh, err = os.Open(filename) // nolint: gosec
	if os.IsNotExist(err) {
		err = img.Errors().ErrNotFound()
		return
	} else if err != nil {
		return
	}
	i = img.New().(*imgItem)
	if _, err = i.ReadFrom(fh); err != nil {
		return
	}
	i.fileName = filename
	ret = i

	return
}
