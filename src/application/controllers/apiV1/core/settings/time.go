package settings // import "application/controllers/apiV1/core/settings"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/context"
import "gopkg.in/webnice/kit.v1/modules/verify"
import (
	"encoding/json"
	"net/http"
	"time"
)

// New Create new object and return interface
func New() Interface { return new(impl) }

// Time Check time difference and return server time in UTC
// PUT /api/v1.0/settings/time
func (sts *impl) Time(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var buf []byte
	var ctx context.Interface
	var req *RequestTime
	var rsp *ResponseTime

	// Получение запроса
	ctx = context.New(rq)
	req = new(RequestTime)
	if buf, err = ctx.Data(req); err != nil {
		wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
		wr.WriteHeader(status.BadRequest)
		wr.Write(buf) // nolint: errcheck, gosec
		return
	}
	// Check zero time
	if req.Time.IsZero() {
		log.Noticef("Request witch zero time '%s'", req.Time.String())
		buf = verify.E4xx().Add(verify.Error{
			Field:      "time",
			FieldValue: req.Time.String(),
			Message:    `Submitted empty time, you must send a local client time`,
		}).Response().Json()
		wr.WriteHeader(status.BadRequest)
		wr.Write(buf) // nolint: errcheck, gosec
		return
	}
	// Ответ
	rsp = new(ResponseTime)
	rsp.Time = time.Now().In(time.UTC)
	rsp.Delta = rsp.Time.Sub(req.Time)
	// Кодирование json
	if buf, err = json.Marshal(rsp); err != nil {
		log.Errorf("json encode error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	wr.WriteHeader(status.Ok)
	wr.Write(buf) // nolint: errcheck, gosec
}
