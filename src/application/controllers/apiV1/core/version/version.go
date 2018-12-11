package version // import "application/controllers/apiV1/core/version"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"application/configuration"
	"application/models/goose"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/status"
)

// Interface is an interface of controller
type Interface interface {
	// Version Метод возвращает текущую версию
	Version(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of Controller
type impl struct {
}

// Response Структура ответа
type Response struct {
	Version     string    `json:"version"`     // Текущая версия приложения
	VersionDb   string    `json:"versionDb"`   // Текущая версия схемы базы данных (уникальный идентификатор миграции)
	VersionDbAt time.Time `json:"versionDbAt"` // Дата и время применения миграции базы данных
}

// New Create new object and return interface
func New() Interface { return new(impl) }

// Version Метод возвращает текущую версию
// GET /api/v1.0/version
func (vrs *impl) Version(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var ver *goose.DbVersion
	var rsp *Response
	var buf []byte

	rsp = &Response{
		Version: configuration.Get().Version().String(),
	}
	if ver, err = goose.New().CurrentVersion(); err == nil {
		rsp.VersionDb = fmt.Sprintf("%d", ver.VersionID)
		rsp.VersionDbAt = ver.TimeStamp.MustValue()
	} else {
		log.Errorf("Goose DB model error: %s", err)
	}
	if buf, err = json.Marshal(rsp); err != nil {
		log.Errorf("json encode error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	wr.WriteHeader(status.Ok)
	if _, err = wr.Write(buf); err != nil {
		log.Errorf("response error: %s", err)
	}
}
