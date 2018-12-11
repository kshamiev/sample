package myitem

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import "gopkg.in/webnice/web.v1/context"
import "gopkg.in/webnice/kit.v1/modules/verify"
import "gopkg.in/webnice/cpy.v1/cpy"
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	modelsMyitem "application/models/myitem"
	myitemTypes "application/models/myitem/types"
)

// New creates a new object and return interface
func New() Interface {
	var mic = new(impl)
	return mic
}

// Status Получение состояние данных в БД
// OPTIONS /api/v1.0/myitem
func (mic *impl) Status(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var data *modelsMyitem.StatusInfo
	var rsp *StatusResponse
	var buf []byte

	// Выполнение действия
	if data, err = modelsMyitem.New().Status(); err != nil {
		log.Errorf("Model error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Подготовка ответа
	rsp = new(StatusResponse)
	// СПЕЦИАЛЬНО СДЕЛАНО НЕ СОВПОДЕНИЕ ТИПА ПОЛЯ И ИМЕНИ ПОЛЯ И ПОКАЗАНО КАК ЭТО МОЖНО РЕШИТЬ
	if err = cpy.All(rsp, data); err != nil {
		log.Errorf("Copy data error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Кодирование json
	if buf, err = json.Marshal(rsp); err != nil {
		log.Errorf("json encode error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Add(header.RetryAfter, fmt.Sprintf("%d", uint64(retryAfter)))
	wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	wr.WriteHeader(status.Ok)
	if _, err = wr.Write(buf); err != nil {
		log.Errorf("response error: %s", err)
	}
}

// Create Сохранение в БД данных
// POST /api/v1.0/myitem
func (mic *impl) Create(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var ctx context.Interface
	var req *CreateRequest
	var rsp *CreateResponse
	var itt *myitemTypes.Myitem
	var buf []byte

	// Получение данных запроса
	ctx = context.New(rq)
	req = new(CreateRequest)
	if buf, err = ctx.Data(req); err != nil {
		// Ошибка проверки данных, используется "gopkg.in/go-playground/validator.v9"
		wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
		wr.Header().Add(header.RetryAfter, fmt.Sprintf("%d", uint64(retryAfter)))
		wr.WriteHeader(status.BadRequest)
		if _, err = wr.Write(buf); err != nil {
			log.Errorf("response error: %s", err)
		}
		return
	}
	// Выполнение действия
	itt, err = modelsMyitem.New().Create(req.Date, req.Number, req.Text)
	if err != nil {
		log.Errorf("Model error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Подготовка ответа
	rsp = new(CreateResponse)
	if err = cpy.All(rsp, itt); err != nil {
		log.Errorf("Copy data error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Кодирование json
	if buf, err = json.Marshal(rsp); err != nil {
		log.Errorf("json encode error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Add(header.RetryAfter, fmt.Sprintf("%d", uint64(retryAfter)))
	wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	wr.WriteHeader(status.Ok)
	if _, err = wr.Write(buf); err != nil {
		log.Errorf("response error: %s", err)
	}
}

// Load Получение данных из БД по ID
// GET /api/v1.0/myitem/:id
func (mic *impl) Load(wr http.ResponseWriter, rq *http.Request) {
	const idKey = `id`
	var err error
	var idSrc string
	var id uint64
	var mim modelsMyitem.Interface
	var rsp *LoadResponse
	var data *myitemTypes.Myitem
	var buf []byte

	// Получение данных запроса
	idSrc = context.New(rq).Route().Params().Get(idKey)
	id, err = strconv.ParseUint(idSrc, 10, 64)
	if err != nil {
		err = fmt.Errorf("Invalid ID of item: %q, error: %s", idSrc, err)
		log.Error(err)
		wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
		wr.WriteHeader(status.BadRequest)
		_, err = wr.Write(verify.E4xx().Code(-1).Add(verify.Error{Field: "id", FieldValue: idSrc, Message: err.Error()}).Message(err.Error()).Json())
		if err != nil {
			log.Errorf("response error: %s", err)
		}
		return
	}
	// Выполнение действия
	mim = modelsMyitem.New()
	data, err = mim.Load(id)
	// Проверка результата действия
	switch err {
	case mim.ErrNotFound():
		wr.WriteHeader(status.NotFound)
		return
	default:
		if err != nil {
			log.Errorf("Model Delete(%d) error: %s", id, err)
			wr.WriteHeader(status.InternalServerError)
			return
		}
	}
	// Подготовка ответа
	rsp = new(LoadResponse)
	if err = cpy.All(rsp, data); err != nil {
		log.Errorf("Copy data error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Кодирование json
	if buf, err = json.Marshal(rsp); err != nil {
		log.Errorf("json encode error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Add(header.RetryAfter, fmt.Sprintf("%d", uint64(retryAfter)))
	wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	wr.WriteHeader(status.Ok)
	if _, err = wr.Write(buf); err != nil {
		log.Errorf("response error: %s", err)
	}
}

// Delete Удаление данных в БД по ID
// DELETE /api/v1.0/myitem/:id
func (mic *impl) Delete(wr http.ResponseWriter, rq *http.Request) {
	const idKey = `id`
	var err error
	var mim modelsMyitem.Interface
	var idSrc string
	var id uint64

	// Получение данных запроса
	idSrc = context.New(rq).Route().Params().Get(idKey)
	id, err = strconv.ParseUint(idSrc, 10, 64)
	if err != nil {
		err = fmt.Errorf("Invalid ID of item: %q, error: %s", idSrc, err)
		log.Error(err)
		wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
		wr.WriteHeader(status.BadRequest)
		_, err = wr.Write(verify.E4xx().Code(-1).Add(verify.Error{Field: "id", FieldValue: idSrc, Message: err.Error()}).Message(err.Error()).Json())
		if err != nil {
			log.Errorf("response error: %s", err)
		}
		return
	}
	// Выполнение действия
	mim = modelsMyitem.New()
	err = mim.Delete(id)
	// Подготовка ответа
	switch err {
	case mim.ErrNotFound():
		wr.WriteHeader(status.NotFound)
		return
	default:
		if err != nil {
			log.Errorf("Model Delete(%d) error: %s", id, err)
			wr.WriteHeader(status.InternalServerError)
			return
		}
	}
	wr.WriteHeader(status.NoContent)
}
