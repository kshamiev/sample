package filecache // import "application/models/filecache"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// Все ошибки определены как константы
const (
	cNotFound           = "Not found"
	cNameNotSpecified   = "Name of object is not specified"
	cCreatorFnIsNil     = "Not specified a function to create a new object"
	cCreatorFnReturnNil = "Function of creating a new object returned an nil object"
	cCopyDataFailed     = "Copying data failed, destination data length do not match source data length"
)

// Константы указаны в объектах, адрес которых фиксирован всё время работы приложения
// Ошибку с ошибкой можно сравнивать по телу, по адресу и т.п.
var (
	errSingleton          = &Error{}
	errNotFound           = err(cNotFound)
	errNameNotSpecified   = err(cNameNotSpecified)
	errCreatorFnIsNil     = err(cCreatorFnIsNil)
	errCreatorFnReturnNil = err(cCreatorFnReturnNil)
	errCopyDataFailed     = err(cCopyDataFailed)
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

// ErrNameNotSpecified Name of object is not specified
func (e *Error) ErrNameNotSpecified() error { return &errNameNotSpecified }

// ErrCreatorFnIsNil Not specified a function to create a new object
func (e *Error) ErrCreatorFnIsNil() error { return &errCreatorFnIsNil }

// ErrCreatorFnReturnNil The function of creating a new object returned an nil object
func (e *Error) ErrCreatorFnReturnNil() error { return &errCreatorFnReturnNil }

// ErrCopyDataFailed Copying data failed, destination data length do not match source data length
func (e *Error) ErrCopyDataFailed() error { return &errCopyDataFailed }
