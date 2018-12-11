package retrycheck // import "application/middleware/retrycheck"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"strconv"

	"gopkg.in/webnice/web.v1/header"
)

// Handler Проверка, контроль и информирование о количестве запросов от одного клиента
func Handler(hndl http.Handler) http.Handler {
	var hfn = func(wr http.ResponseWriter, rq *http.Request) {
		wr.Header().Add(header.RetryAfter, strconv.FormatUint(RetryAfterDefault, 10))
		hndl.ServeHTTP(wr, rq)
	}

	return http.HandlerFunc(hfn)
}
