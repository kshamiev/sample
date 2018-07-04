package panic // import "application/controllers/apiV1/core/panic"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
)

// Interface is an interface of controller
type Interface interface {
	// Panic is a method for testing panic
	Panic(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of Controller
type impl struct {
}

// New Create new object and return interface
func New() Interface { return new(impl) }

// Panic is a method for testing panic
// GET /api/v1.0/panic
func (pnc *impl) Panic(wr http.ResponseWriter, rq *http.Request) {
	panic("Test panic")
}
