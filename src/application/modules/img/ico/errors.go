package ico // import "application/modules/img/ico"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// FormatError A FormatError reports that the input is not a valid ICO.
type FormatError string

// Error Возврат ошибки
func (e FormatError) Error() string { return "Invalid ICO format: " + string(e) }
