package sitemap // import "application/models/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// Все ошибки определены как константы
const (
	cNotFound = `Not found`
	cTooLarge = `Requested data exceeded maximum supported data size`
)

// Константы указаны в объектах, адрес которых фиксирован всё время работы приложения
// Ошибку с ошибкой можно сравнивать по телу, по адресу и т.п.
var (
	errSingleton = &Error{}
	errNotFound  = err(cNotFound)
	errTooLarge  = err(cTooLarge)
)

type (
	// Error object of package
	Error struct{}
	err   string
)

// Error The error built-in interface implementation
func (e err) Error() string { return string(e) }

// Errors Все ошибки известного состояния, которые могут вернуть функции пакета
func Errors() *Error { return errSingleton }

// ERRORS:

// ErrNotFound Not found
func (e *Error) ErrNotFound() error { return &errNotFound }

// ErrTooLarge Requested data exceeded maximum supported data size
func (e *Error) ErrTooLarge() error { return &errTooLarge }
