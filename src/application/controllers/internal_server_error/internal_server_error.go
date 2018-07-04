package internal_server_error // import "application/controllers/internal_server_error"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/context"
import "gopkg.in/webnice/kit.v1/modules/verify"
import (
	"fmt"
	"net/http"
	"strings"
)

// Interface is an interface of controller
type Interface interface {
	// InternalServerError Кастомный обработчик ответа внутренней ошибки сервера
	InternalServerError(http.ResponseWriter, *http.Request)
}

// impl is an implementation of Controller
type impl struct {
}

// New Create new object and return interface
func New() Interface { return new(impl) }

// InternalServerError Кастомный обработчик ответа внутренней ошибки сервера
// TODO
// В будущем сделать ответ на основе HTML шаблона для не JSON запросов
func (ise *impl) InternalServerError(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var buf []byte

	err = context.New(rq).Errors().InternalServerError(nil)
	if strings.Contains(rq.Header.Get(header.ContentType), mime.ApplicationJSON) {
		// JSON
		buf = verify.E5xx().Code(-1).Message(err.Error()).Response().Json()
		wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
		wr.WriteHeader(status.InternalServerError)
		wr.Write(buf)
	} else {
		// TEXT/HTML
		wr.Header().Add(header.ContentType, mime.TextPlainCharsetUTF8)
		wr.WriteHeader(status.InternalServerError)
		fmt.Fprintf(wr, "%s", status.Text(status.InternalServerError))
	}
	log.Errorf("Error:\n%s", err.Error())
}
