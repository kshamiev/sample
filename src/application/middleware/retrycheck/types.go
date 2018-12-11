package retrycheck // import "application/middleware/retrycheck"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

const (
	// RetryAfterDefault Is an default value for header Retry-After
	RetryAfterDefault = uint64(3) // Количество секунд до следующего запроса
)
