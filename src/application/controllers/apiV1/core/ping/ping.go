package ping // import "application/controllers/apiV1/core/ping"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import (
	"net/http"
)

// Interface is an interface of controller
type Interface interface {
	// Ping is a method for checking service availability
	Ping(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of Controller
type impl struct {
}

// New Create new object and return interface
func New() Interface { return new(impl) }

// Ping is a method for checking service availability
// GET /api/v1.0/ping
func (ping *impl) Ping(wr http.ResponseWriter, rq *http.Request) {
	wr.WriteHeader(status.NoContent)
}
