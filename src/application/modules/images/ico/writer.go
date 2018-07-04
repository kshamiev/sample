package ico // import "application/modules/images/ico"

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
	var id icondir
	id.imageType = 1
	id.numImages = 1
	id.reserved = 0 // for structcheck: fix warning: unused struct field
	return id
}

func newIcondirentry() icondirentry {
	var ide icondirentry
	ide.colorPlanes = 1   // windows is supposed to not mind 0 or 1, but other icon files seem to have 1 here
	ide.bitsPerPixel = 32 // can be 24 for bitmap or 24/32 for png. Set to 32 for now
	ide.offset = 22       // 6 icondir + 16 icondirentry, next image will be this image size + 16 icondirentry, etc
	ide.reserved = 0      // for structcheck: fix warning: unused struct field
	ide.numColors = 0     // for structcheck: fix warning: unused struct field
	return ide
}

// Encode image to ico format
func Encode(w io.Writer, im image.Image) (err error) {
	var b = im.Bounds()
	var m = image.NewRGBA(b)
	var bb = new(bytes.Buffer)

	draw.Draw(m, b, im, b.Min, draw.Src)

	var id = newIcondir()
	var ide = newIcondirentry()

	pngbb := new(bytes.Buffer)
	pngwriter := bufio.NewWriter(pngbb)
	err = png.Encode(pngwriter, m)
	err = pngwriter.Flush()
	ide.sizeInBytes = uint32(len(pngbb.Bytes()))

	var bounds = m.Bounds()
	ide.imageWidth = uint8(bounds.Dx())
	ide.imageHeight = uint8(bounds.Dy())

	err = binary.Write(bb, binary.LittleEndian, id)
	err = binary.Write(bb, binary.LittleEndian, ide)

	_, err = w.Write(bb.Bytes())
	_, err = w.Write(pngbb.Bytes())

	return
}
