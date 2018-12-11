package settings // import "application/models/settings"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// Все ошибки определены как константы
const (
	cKeyOrValueNotFound = `Key or value not found`
	cKeyIsNotUnique     = `Key is not unique`
)

// Константы указаны в объектах, адрес которых фиксирован всё время работы приложения
// Ошибку с ошибкой можно сравнивать по телу, по адресу и т.п.
var (
	errSingleton          = &Error{}
	errKeyOrValueNotFound = err(cKeyOrValueNotFound)
	errKeyIsNotUnique     = err(cKeyIsNotUnique)
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

// ErrKeyOrValueNotFound Key or value not found
func (e *Error) ErrKeyOrValueNotFound() error { return &errKeyOrValueNotFound }

// ErrKeyIsNotUnique Key is not unique
func (e *Error) ErrKeyIsNotUnique() error { return &errKeyIsNotUnique }
