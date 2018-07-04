package assets // import "application/controllers/pages/assets"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/header"
import (
	"fmt"
	"net/http"
	"time"

	"application/configuration"
	"application/modules/pages"
)

// New creates a new object and return interface
func New() Interface {
	var ass = new(impl)
	return ass
}

// Lazy initialization
func (ass *impl) init() {
	if ass.cfg == nil {
		ass.cfg = configuration.Get()
	}
	if ass.pgs == nil {
		// first web server
		for i := range ass.cfg.Configuration().WEBServers {
			ass.pgs = pages.New(ass.cfg.Configuration().WEBServers[i].DocumentRoot)
			break
		}
	}
}

// Assets Statisc files
func (ass *impl) Assets(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var rfn string
	var ims time.Time
	var fle *pages.File

	// Ленивая инициализация контроллера
	ass.init()
	// Запрашиваемый файл
	rfn = rq.RequestURI
	// Если что-то случилось с файлом
	if fle, err = ass.pgs.Assets(rfn); err == ass.pgs.ErrFileNotFound() {
		wr.WriteHeader(status.NotFound)
		return
	} else if err != nil {
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Если клиент поддерживает If-Modified-Since, загружаем заголовки
	ims, err = time.Parse(ifModifiedSinceTimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && fle.Info.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Выдача полного комплекта заголовков
	if fle.ContentType != "" {
		wr.Header().Set(header.ContentType, fle.ContentType)
	}
	wr.Header().Set(header.ContentLength, fmt.Sprintf("%d", fle.Body.Len()))
	wr.Header().Set(header.LastModified, fle.Info.ModTime().UTC().Format(ifModifiedSinceTimeFormat))
	wr.WriteHeader(status.Ok)
	if _, err = fle.Body.WriteTo(wr); err != nil {
		log.Errorf("Response error: %s", err.Error())
		return
	}
}
