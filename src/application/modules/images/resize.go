package images // import "application/modules/images"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	// Библиотека image
	_ "image"

	// Расширение image для формата png
	_ "image/png"

	// Расширение image для формата gif
	_ "image/gif"

	// Расширение image для формата jpeg
	_ "image/jpeg"

	// Расширение image для формата ico
	_ "application/modules/images/ico"

	// Расширение image для формата bmp
	_ "github.com/jsummers/gobmp"

	"github.com/nfnt/resize"
)

// Resize Resize image
func (img *impl) Resize(im *Image, w, h uint) (ret *Image) {
	ret = new(Image)

	ret.Type = im.Type
	ret.FileName = im.FileName
	ret.FileInfo = im.FileInfo

	// Resize
	ret.Image = resize.Resize(w, h, im.Image, resize.Lanczos3)

	// Create image.Config
	img.createConfig(ret)

	return
}
