package index // import "application/controllers/pages/index"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/context"
import (
	"fmt"
	"net/http"
	"time"

	"application/configuration"
	"application/modules/pages"
)

// New creates a new object and return interface
func New() Interface {
	var idx = new(impl)
	return idx
}

// Lazy initialization
func (idx *impl) init() {
	if idx.cfg == nil {
		idx.cfg = configuration.Get()
	}
	if idx.pgs == nil {
		// first web server
		for i := range idx.cfg.Configuration().WEBServers {
			idx.pgs = pages.New(idx.cfg.Configuration().WEBServers[i].DocumentRoot)
			break
		}
	}
}

// Index page
// GET /
func (idx *impl) Index(wr http.ResponseWriter, rq *http.Request) {
	const keyFilename = `filename`
	var err error
	var rfn string
	var ctx context.Interface
	var ims time.Time
	var fle *pages.File

	// Ленивая инициализация контроллера
	idx.init()

	// Запрашиваемый файл
	ctx = context.New(rq)
	rfn = ctx.Route().Params().Get(keyFilename)

	// Если что-то случилось с файлом
	if fle, err = idx.pgs.Index(rfn); err == idx.pgs.ErrFileNotFound() {
		wr.WriteHeader(status.NotFound)
		return
	} else if err != nil {
		wr.WriteHeader(status.InternalServerError)
		return
	}

	// Если клиент поддерживает If-Modified-Since, загружаем заголовки
	ims, err = time.Parse(_IfModifiedSinceTimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && fle.Info.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}

	// Выдача полного комплекта заголовков
	if fle.ContentType != "" {
		wr.Header().Set(header.ContentType, fle.ContentType)
	}
	wr.Header().Set(header.ContentLength, fmt.Sprintf("%d", fle.Body.Len()))
	wr.Header().Set(header.LastModified, fle.Info.ModTime().UTC().Format(_IfModifiedSinceTimeFormat))
	wr.WriteHeader(status.Ok)
	if _, err = fle.Body.WriteTo(wr); err != nil {
		log.Errorf("Response error: %s", err.Error())
		return
	}
}
