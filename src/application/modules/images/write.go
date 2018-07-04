package images // import "application/modules/images"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"

	// Библиотека image
	_ "image"

	// Расширение image для формата png
	"image/png"

	// Расширение image для формата gif
	_ "image/gif"

	// Расширение image для формата jpeg
	_ "image/jpeg"

	// Расширение image для формата ico
	"application/modules/images/ico"

	// Расширение image для формата bmp
	_ "github.com/jsummers/gobmp"
)

// WriteIco Write image as ico format to writer
func (img *impl) WriteIco(w io.Writer, im *Image) (err error) {
	err = ico.Encode(w, im.Image)
	return
}

// WritePng Write image as png format to writer
func (img *impl) WritePng(w io.Writer, im *Image) (err error) {
	err = png.Encode(w, im.Image)
	return
}
