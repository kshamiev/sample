package images // import "application/modules/images"

import (
	"image"
)

// convertGrayscale Convert image to grayscale image
func (img *impl) convertGrayscale(in *Image, fn convertGrayscaleFunc) (ret *Image) {
	var rectangle image.Rectangle
	var pos, x, y int
	var gray *image.Gray

	rectangle = in.Image.Bounds()

	gray = new(image.Gray)
	gray = image.NewGray(rectangle)
	for y = rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x = rectangle.Min.X; x < rectangle.Max.X; x++ {
			gray.Pix[pos] = fn(in.Image.At(x, y)).Y
			pos++
		}
	}

	// Copy info
	ret = new(Image)
	ret.FileInfo = in.FileInfo
	ret.FileName = in.FileName
	ret.Type = ret.Type
	ret.Image = gray

	// Create image.Config
	img.createConfig(ret)

	return
}

// ConvertBlackAndWhite Convert image to grayscale image
func (img *impl) ConvertBlackAndWhite(in *Image) *Image {
	return img.convertGrayscale(in, grayscaleBlackAndWhite70percent)
}
