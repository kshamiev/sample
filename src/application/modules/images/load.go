package images // import "application/modules/images"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"image"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

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
)

// LoadImageFile Load image from disk
func (img *impl) LoadImageFile(fileName string) (ret *Image, err error) {
	var fh *os.File

	fh, err = os.Open(fileName)
	if err != nil {
		return
	}
	defer func() {
		_ = fh.Close()
	}()

	// Чтение картинки из Reader
	if ret, err = img.LoadImageReader(fh); err != nil {
		return
	}

	// FileInfo
	if ret.FileInfo, err = fh.Stat(); err != nil {
		return
	}

	return
}

// LoadImageReader Load image from io.Reader
func (img *impl) LoadImageReader(fh io.Reader) (ret *Image, err error) {
	// image паникует, сука...
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			log.Warning("Recovery panic in image.Image: %s", err.Error())
		}
	}()

	ret = new(Image)

	// Content
	if ret.Image, ret.Type, err = image.Decode(fh); err != nil {
		return
	}

	// Create image.Config
	img.createConfig(ret)

	return
}

// createConfig Create image.Config
func (img *impl) createConfig(im *Image) {
	var rect image.Rectangle

	rect = im.Image.Bounds()
	im.Config.Width = rect.Max.X
	im.Config.Height = rect.Max.Y
	im.Config.ColorModel = im.Image.ColorModel()
}

// FindInFolderAndLoad Чтение из папки всех файлов являющихся картинкой и удовлетворяющих патерну
func (img *impl) FindInFolderAndLoad(dirName string, pattern *regexp.Regexp) (ret []*Image, err error) {
	var dir []os.FileInfo
	var i int

	if dir, err = ioutil.ReadDir(dirName); err != nil {
		return
	}

	for i = range dir {
		var im *Image
		var fn, fileName string

		if dir[i].IsDir() {
			continue
		}

		fileName = strings.ToLower(dir[i].Name())
		if !pattern.MatchString(fileName) {
			continue
		}

		fn = path.Join(dirName, fileName)
		if im, err = img.LoadImageFile(fn); err != nil {
			continue
		}

		im.FileName = fileName
		ret = append(ret, im)
	}

	return
}
