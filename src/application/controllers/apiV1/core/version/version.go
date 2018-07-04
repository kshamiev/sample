package version // import "application/controllers/apiV1/core/version"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import (
	"encoding/json"
	"net/http"

	"application/configuration"
)

// Interface is an interface of controller
type Interface interface {
	// Version is a method for checking service availability
	Version(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of Controller
type impl struct {
}

// VersionResponse Структура ответа
type VersionResponse struct {
	Version string `json:"version"`
}

// New Create new object and return interface
func New() Interface { return new(impl) }

// Version Метод возвращает текущую версию приложения
// GET /api/v1.0/ping
func (vrs *impl) Version(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var resp *VersionResponse
	var buf []byte

	resp = &VersionResponse{
		Version: configuration.Get().Version().String(),
	}
	if buf, err = json.Marshal(resp); err != nil {
		log.Errorf("json encode error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	wr.WriteHeader(status.Ok)
	wr.Write(buf)
}
