package uploadfile // import "application/controllers/apiV1/core/uploadfile"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"application/models/filestore"

	"gopkg.in/webnice/kit.v1/modules/verify"
	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/status"
)

// New creates a new object and return interface
func New() Interface {
	var ufc = new(impl)
	return ufc
}

// Ленивая инициализация и возврат интерфейса загрузки и хранения файлов
func (ufc *impl) lazyUploadFile() filestore.Interface {
	if ufc.ufm == nil {
		ufc.ufm = filestore.New()
	}
	return ufc.ufm
}

// UploadFile Загрузка файла на сервер
func (ufc *impl) UploadFile(wr http.ResponseWriter, rq *http.Request) {
	const mimeFileValue = `multipart/form-data`
	var err error
	var verrors []verify.Error
	var formFieldName, filename string
	var item *Response
	var rsp []*Response
	var mfh *multipart.FileHeader
	var fh io.ReadCloser
	var buf []byte
	var i int

	// 1Gb памяти для парсинга multipart в памяти
	if err = rq.ParseMultipartForm(1 << 30); err != nil {
		wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
		wr.WriteHeader(status.BadRequest)
		_, err = wr.Write(verify.E4xx().Code(-1).Message(err.Error()).Json())
		if err != nil {
			log.Errorf("response error: %s", err)
		}
		return
	}
	for formFieldName = range rq.MultipartForm.File {
		for i = range rq.MultipartForm.File[formFieldName] {
			mfh = rq.MultipartForm.File[formFieldName][i]
			if filename, err = strconv.Unquote(fmt.Sprintf(`"%s"`, mfh.Filename)); err == nil {
				mfh.Filename = filename
			}
			item = &Response{
				Field:       formFieldName,
				Filename:    mfh.Filename,
				Size:        uint64(mfh.Size),
				ContentType: mfh.Header.Get(header.ContentType),
			}
			if fh, err = mfh.Open(); err != nil {
				verrors = append(verrors, verify.Error{
					Field:      fmt.Sprintf("%s[%d]", formFieldName, i),
					FieldValue: mimeFileValue,
					Message:    fmt.Sprintf("Error open source: %s", err),
				})
				continue
			}
			fh.Close() // nolint: gosec, errcheck
			item.ID, err = ufc.lazyUploadFile().
				NewTemporaryFile(item.Filename, item.Size, item.ContentType, fh)
			if err != nil {
				verrors = append(verrors, verify.Error{
					Field:      fmt.Sprintf("%s[%d]", formFieldName, i),
					FieldValue: mimeFileValue,
					Message:    fmt.Sprintf("Error store new temporary file: %s", err),
				})
				continue
			}
			rsp = append(rsp, item)
		}
	}
	// Вывод ошибок
	if len(verrors) > 0 {
		wr.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
		wr.WriteHeader(status.InternalServerError)
		ve := verify.E5xx().Code(-1).Message("Save file error")
		for _, veo := range verrors {
			ve.Add(veo)
		}
		if _, err = wr.Write(ve.Json()); err != nil {
			log.Errorf("response error: %s", err)
		}
		return
	}
	// Кодирование json
	if len(rsp) == 0 {
		buf = []byte(`[]`)
	} else if buf, err = json.Marshal(rsp); err != nil {
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
