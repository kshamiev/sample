package images // import "application/modules/images"

import (
	"image/color"
	//"math"
)

// ConvertGrayscale Convert image to grayscale image
func (img *impl) ConvertGrayscale(in *Image) *Image {
	return img.convertGrayscale(in, grayscaleValue)
}

// grayscaleAverage The formula used for conversion is: Y = (r + g + b) / 3.
//func grayscaleAverage(c color.Color) color.Gray {
//	r, g, b, _ := c.RGBA()
//	Y := (10*(r+g+b) + 5) / 10
//	return color.Gray{uint8(Y >> 8)}
//}

// grayscaleGrayValue The formula used for conversion is: Y = max(r,g,b)
func grayscaleValue(c color.Color) color.Gray {
	r, g, b, _ := c.RGBA()
	Y := max(r, max(g, b))
	return color.Gray{uint8(Y >> 8)}
}

// grayscaleLuma The formula used for conversion is: Y = 0.299*r + 0.587*g + 0.114*b.
// The same formula is used by color.GrayModel.Convert().
//func grayscaleLuma(c color.Color) color.Gray {
//	r, g, b, _ := c.RGBA()
//	Y := (299*r + 587*g + 114*b + 500) / 1000
//	return color.Gray{uint8(Y >> 8)}
//}

// grayscaleLuma709 The formula used for conversion is: Y = 0.299*r + 0.587*g + 0.114*b.
// The same formula is used by color.GrayModel.Convert().
//func grayscaleLuma709(c color.Color) color.Gray {
//	r, g, b, _ := c.RGBA()
//	Y := (2125*r + 7154*g + 721*b + 5000) / 10000
//	return color.Gray{uint8(Y >> 8)}
//}

// The formula used for conversion is: Y' = 0.2125*R' + 0.7154*G' + 0.0721*B'
// where r, g and b are gamma expanded with gamma 2.2 and final Y is Y'
// gamma compressed again.
// The same formula is used by color.GrayModel.Convert().
//func grayscaleLuminance(c color.Color) color.Gray {
//	rr, gg, bb, _ := c.RGBA()
//	r := math.Pow(float64(rr), 2.2)
//	g := math.Pow(float64(gg), 2.2)
//	b := math.Pow(float64(bb), 2.2)
//	y := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
//	Y := uint16(y + 0.5)
//	return color.Gray{uint8(Y >> 8)}
//}

// grayscaleLightness The formula used for conversion is: Y = (max(r,g,b) + min(r,g,b)) / 2.
//func grayscaleLightness(c color.Color) color.Gray {
//	r, g, b, _ := c.RGBA()
//	max := max(r, max(g, b))
//	min := min(r, min(g, b))
//	Y := (10*(min+max) + 5) / 20
//	return color.Gray{uint8(Y >> 8)}
//}

// grayscaleRed converts color.Color c to grayscale using the R component.
//func grayscaleRed(c color.Color) color.Gray {
//	r, _, _, _ := c.RGBA()
//	return color.Gray{uint8(r >> 8)}
//}

// grayscaleGreen converts color.Color c to grayscale using the G component.
//func grayscaleGreen(c color.Color) color.Gray {
//	_, g, _, _ := c.RGBA()
//	return color.Gray{uint8(g >> 8)}
//}

// grayscaleBlue converts color.Color c to grayscale using the B component.
//func grayscaleBlue(c color.Color) color.Gray {
//	_, _, b, _ := c.RGBA()
//	return color.Gray{uint8(b >> 8)}
//}

// grayscaleBlackAndWhite50percent The formula used for conversion is: if sum(r+g+b)/3 > 50% then-> white else-> black
//func grayscaleBlackAndWhite50percent(c color.Color) (ret color.Gray) {
//	r, g, b, _ := c.RGBA()
//	t := uint64(uint64(r)+uint64(g)+uint64(b)) / 3
//	if t > 65535/2 {
//		ret = color.Gray{255}
//	} else {
//		ret = color.Gray{0}
//	}
//	return
//}

// grayscaleBlackAndWhite60percent The formula used for conversion is: if sum(r+g+b)/3 > 60% then-> white else-> black
//func grayscaleBlackAndWhite60percent(c color.Color) (ret color.Gray) {
//	r, g, b, _ := c.RGBA()
//	t := uint64(uint64(r)+uint64(g)+uint64(b)) / 3
//	if t > 65535/100*60 {
//		ret = color.Gray{255}
//	} else {
//		ret = color.Gray{0}
//	}
//	return
//}

// grayscaleBlackAndWhite70percent The formula used for conversion is: if sum(r+g+b)/3 > 70% then-> white else-> black
func grayscaleBlackAndWhite70percent(c color.Color) (ret color.Gray) {
	r, g, b, _ := c.RGBA()
	t := (uint64(r) + uint64(g) + uint64(b)) / 3
	if t > 65535/100*70 {
		ret = color.Gray{255}
	} else {
		ret = color.Gray{0}
	}
	return
}

func max(a, b uint32) (ret uint32) {
	ret = b
	if a > b {
		ret = a
	}
	return
}

//func min(a, b uint32) (ret uint32) {
//	ret = b
//	if a < b {
//		ret = a
//	}
//	return
//}
