package pages // import "application/controllers/resource/pages"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"path"
	"strconv"
	"time"

	"application/models/filecache"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/status"
)

// Выгрузка содержимого файла index.html
func (pgi *impl) IndexHandlerFunc(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var ims time.Time
	var indexFN string
	var content filecache.Data

	// Загрузка файла index.html
	indexFN = path.Join(pgi.rootPath, keyIndexHTML)
	if content, err = pgi.Mfc.Load(indexFN); err != nil {
		log.Errorf("filecache model error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Если клиент поддерживает If-Modified-Since, загружаем заголовки
	ims, err = time.Parse(ifModifiedSinceTimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && content.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Выдача всех заголовков и контента
	wr.Header().Set(header.ContentType, mime.TextHTMLCharsetUTF8)
	wr.Header().Set(header.LastModified, content.ModTime().UTC().Format(ifModifiedSinceTimeFormat))
	wr.Header().Set(header.ContentLength, strconv.FormatUint(content.Length(), 10))
	wr.WriteHeader(status.Ok)
	if _, err = content.WriteTo(wr); err != nil {
		log.Errorf("content response error: %s", err)
		return
	}
}
