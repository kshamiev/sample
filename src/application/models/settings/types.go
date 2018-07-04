package settings // import "application/models/settings"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/kit.v1/modules/db"
import (
	"time"
)

// Interface is an interface
type Interface interface {
	// StringSet Запись значения string
	StringSet(key string, value string)

	// StringGet Чтение значения string
	StringGet(key string) string

	// TimeSet Запись значения time.Time
	TimeSet(key string, value time.Time)

	// TimeGet Чтение значения time.Time
	TimeGet(key string) time.Time

	// UintSet Запись значения uint64
	UintSet(key string, value uint64)

	// UintGet Чтение значения uint64
	UintGet(key string) uint64

	// IntSet Запись значения int64
	IntSet(key string, value int64)

	// IntGet Чтение значения int64
	IntGet(key string) int64

	// DecimalSet Запись значения float64 как decimal
	DecimalSet(key string, value float64)

	// DecimalGet Чтение значения float64 как decimal
	DecimalGet(key string) float64

	// FloatSet Запись значения float64 как double
	FloatSet(key string, value float64)

	// FloatGet Чтение значения float64 как double
	FloatGet(key string) float64

	// BooleanSet Запись значения bool
	BooleanSet(key string, value bool)

	// BooleanGet Чтение значения bool
	BooleanGet(key string) bool

	// BlobSet Запись значения []byte
	BlobSet(key string, value []byte)

	// BlobGet Чтение значения []byte
	BlobGet(key string) []byte

	// ERRORS

	// Error Ошибка возникшая в результате последней операции
	Error() error
	// ErrKeyOrValueNotFound Key or value not found
	ErrKeyOrValueNotFound() error
	// ErrKeyIsNotUnique Key is not unique
	ErrKeyIsNotUnique() error
}

// is an implementation
type impl struct {
	db.Implementation       // Наследование
	LastError         error // Ошибка возникшая в результате последней операции
}
