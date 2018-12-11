package img // import "application/modules/img"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import (
//	"image"
//	"image/color"
//)

// convertGrayscale Convert image to grayscale image
//func (img *impl) convertGrayscale(in *Image, fn func(color.Color) color.Gray) (ret *Image) {
//	var rectangle image.Rectangle
//	var pos, x, y int
//	var gray *image.Gray
//
//	rectangle = in.Image.Bounds()
//	gray = image.NewGray(rectangle)
//	for y = rectangle.Min.Y; y < rectangle.Max.Y; y++ {
//		for x = rectangle.Min.X; x < rectangle.Max.X; x++ {
//			gray.Pix[pos] = fn(in.Image.At(x, y)).Y
//			pos++
//		}
//	}
//	// Copy info
//	ret = &Image{
//		FileName: in.FileName,
//		FileInfo: in.FileInfo,
//		Image:    gray,
//		Type:     in.Type,
//	}
//	// Новый конфиг
//	img.createConfig(ret)
//
//	return
//}

// ConvertBlackAndWhite Convert image to grayscale image
//func (img *impl) ConvertBlackAndWhite(in *Image) *Image {
//	return img.convertGrayscale(in, grayscaleBlackAndWhite70percent)
//}
