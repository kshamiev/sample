package ico // import "application/modules/img/ico"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"image"
	"io"
)

const pngHeader = "\x89PNG\r\n\x1a\n"

// If the io.Reader does not also have ReadByte, then decode will introduce its own buffering.
type reader interface {
	io.Reader
	io.ByteReader
}

type decoder struct {
	r     reader
	num   uint16
	dir   []entry
	image []image.Image
	cfg   image.Config
}

type entry struct {
	Width   uint8
	Height  uint8
	Palette uint8
	_       uint8 // Reserved byte
	Plane   uint16
	Bits    uint16
	Size    uint32
	Offset  uint32
}
