package sitemap // import "application/controllers/resource/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// Все ошибки определены как константы
const (
	cNotFound                  = `Not found`
	cURIisToLong               = `URI address is too long`
	cPriorityNotCorrect        = `The priority is not correct. Value must be from 0.0 to 1.0`
	cChangeFrequencyNotCorrect = `Resource change frequency constant is not incorrect`
	cInvalidURL                = `Invalid URL`
)

// Константы указаны в объектах, адрес которых фиксирован всё время работы приложения
// Ошибку с ошибкой можно сравнивать по телу, по адресу и т.п.
var (
	errSingleton                 = &Error{}
	errNotFound                  = err(cNotFound)
	errURIisToLong               = err(cURIisToLong)
	errPriorityNotCorrect        = err(cPriorityNotCorrect)
	errChangeFrequencyNotCorrect = err(cChangeFrequencyNotCorrect)
	errInvalidURL                = err(cInvalidURL)
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

// ErrURIisToLong URI address is too long
func (e Error) ErrURIisToLong() error { return &errURIisToLong }

// ErrPriorityNotCorrect The priority is not correct. Value must be from 0.0 to 1.0
func (e Error) ErrPriorityNotCorrect() error { return &errPriorityNotCorrect }

// ErrChangeFrequencyNotCorrect Resource change frequency constant is not incorrect
func (e Error) ErrChangeFrequencyNotCorrect() error { return &errChangeFrequencyNotCorrect }

// ErrInvalidURL Invalid URL
func (e Error) ErrInvalidURL() error { return &errInvalidURL }
