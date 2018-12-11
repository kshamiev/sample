package ico // import "application/modules/img/ico"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"image/png"
	"io"
)

type icondir struct {
	reserved  uint16
	imageType uint16
	numImages uint16
}

type icondirentry struct {
	imageWidth   uint8
	imageHeight  uint8
	numColors    uint8
	reserved     uint8
	colorPlanes  uint16
	bitsPerPixel uint16
	sizeInBytes  uint32
	offset       uint32
}

func newIcondir() icondir {
	var id = icondir{
		imageType: 1, // Тип картинки
		numImages: 1, // Количество картинок
		reserved:  0, // for structcheck: fix warning: unused struct field
	}
	return id
}

func newIcondirentry() icondirentry {
	var ide = icondirentry{
		colorPlanes:  1,  // windows is supposed to not mind 0 or 1, but other icon files seem to have 1 here
		bitsPerPixel: 32, // can be 24 for bitmap or 24/32 for png. Set to 32 for now
		offset:       22, // 6 icondir + 16 icondirentry, next image will be this image size + 16 icondirentry, etc
		reserved:     0,  // for structcheck: fix warning: unused struct field
		numColors:    0,  // for structcheck: fix warning: unused struct field
	}
	return ide
}

// Encode image to ico format
func Encode(w io.Writer, im image.Image) (err error) {
	var b, bounds image.Rectangle
	var m *image.RGBA
	var bb, pngbb *bytes.Buffer
	var pngwriter *bufio.Writer
	var id icondir
	var ide icondirentry

	b = im.Bounds()
	m = image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)
	id, ide, pngbb = newIcondir(), newIcondirentry(), &bytes.Buffer{}
	pngwriter = bufio.NewWriter(pngbb)
	if err = png.Encode(pngwriter, m); err != nil {
		return
	}
	_ = pngwriter.Flush() // nolint: errcheck, gosec
	ide.sizeInBytes = uint32(len(pngbb.Bytes()))
	bounds = m.Bounds()
	ide.imageWidth, ide.imageHeight, bb = uint8(bounds.Dx()), uint8(bounds.Dy()), &bytes.Buffer{}
	_ = binary.Write(bb, binary.LittleEndian, id)  // nolint: errcheck, gosec
	_ = binary.Write(bb, binary.LittleEndian, ide) // nolint: errcheck, gosec
	if _, err = w.Write(bb.Bytes()); err != nil {
		return
	}
	if _, err = w.Write(pngbb.Bytes()); err != nil {
		return
	}

	return
}
