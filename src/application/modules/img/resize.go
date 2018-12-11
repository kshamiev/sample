// Package img Пакет помогающий работать с графическими озображениями
// nolint: goimports
package img // import "application/modules/img"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"image"
	"image/color"

	_ "application/modules/img/ico" // Расширение image для формата ico
	_ "golang.org/x/image/bmp"      // Расширение image для формата bmp
	_ "golang.org/x/image/tiff"     // Расширение image для формата tiff
	_ "image/gif"                   // Расширение image для формата gif
	_ "image/jpeg"                  // Расширение image для формата jpeg
	_ "image/png"                   // Расширение image для формата png

	"github.com/disintegration/imaging"
)

// Resize Resize image
func (img *impl) Resize(im Image, w, h uint) (ret Image) {
	var min int
	var newImg image.Image

	if w != h {
		newImg, min = imaging.New(int(w), int(h), color.Transparent), int(w)
		if w > h {
			min = int(h)
		}
		newImg = imaging.PasteCenter(newImg, imaging.Resize(im.Image(), min, min, imaging.Lanczos))
	} else {
		newImg = imaging.Resize(im.Image(), int(w), int(h), imaging.Lanczos)
	}
	ret = img.New().
		SetImage(newImg).
		SetType(im.Type()).
		SetConfig(image.Config{
			Width:      newImg.Bounds().Max.X,
			Height:     newImg.Bounds().Max.Y,
			ColorModel: newImg.ColorModel(),
		}).
		SetFileInfo(im.FileInfo()).
		SetFilename(im.Filename())

	return
}
