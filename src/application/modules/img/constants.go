package img // import "application/modules/img"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

const (
	// TypeUnknown Не известный формат графического файла
	TypeUnknown = Type(``)

	// TypeICO Формат графического файла ico
	TypeICO = Type(`ico`)

	// TypeBMP Формат графического файла bmp
	TypeBMP = Type(`bmp`)

	// TypeTIFF Формат графического файла tiff
	TypeTIFF = Type(`tiff`)

	// TypeGIF Формат графического файла gif
	TypeGIF = Type(`gif`)

	// TypeJPEG Формат графического файла jpeg
	TypeJPEG = Type(`jpeg`)

	// TypePNG Формат графического файла png
	TypePNG = Type(`png`)
)

// Type of image
type Type string

// String Convert type to string
func (t Type) String() string { return string(t) }
