package images // import "application/modules/images"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"image"
	"image/color"
	"io"
	"os"
	"regexp"
)

// Interface is an interface
type Interface interface {
	ConvertGrayscale(*Image) *Image
	ConvertBlackAndWhite(*Image) *Image
	FindInFolderAndLoad(string, *regexp.Regexp) ([]*Image, error)
	LoadImageFile(string) (*Image, error)
	LoadImageReader(io.Reader) (*Image, error)
	Resize(im *Image, w, h uint) *Image
	WriteIco(io.Writer, *Image) error
	WritePng(io.Writer, *Image) error
}

// Implementation is an implementation of repository
type impl struct {
}

// Image object
type Image struct {
	FileName string       // Имя файла картинки
	FileInfo os.FileInfo  // Информация о оригинальном файле
	Image    image.Image  // Загруженная и декодированная картинка
	Config   image.Config // Информация о картинке
	Type     string       // Тип картинки
}

// convertGrayscaleFunc Функция конвертации цвета в grayscale
type convertGrayscaleFunc func(color.Color) color.Gray

// New creates new object
func New(args ...interface{}) Interface {
	var o = new(impl)
	// TODO
	// Если потребуется передавать аргументы, разбор args
	return o
}
