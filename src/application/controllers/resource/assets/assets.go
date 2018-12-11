package assets // import "application/controllers/resource/assets"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"application/controllers/resource/pool"
	"application/models/file"
	"application/models/filecache"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/route"
	"gopkg.in/webnice/web.v1/status"
)

// New creates a new object and return interface
func New() Interface {
	var asi = &impl{
		Mfc:  filecache.Get(),
		Pool: pool.New(),
	}
	return asi
}

// Debug Set debug mode
func (asi *impl) Debug(d bool) Interface {
	asi.debug = d
	asi.Pool.Debug(asi.debug)
	asi.Mfc.Debug(asi.debug)
	return asi
}

// DocumentRoot Устанавливает путь к корню веб сервера
func (asi *impl) DocumentRoot(path string) Interface { asi.rootPath = path; return asi }

// SetRouting Установка роутинга к статическим файлам
func (asi *impl) SetRouting(rou route.Interface) Interface {
	var err error
	var mfi file.Interface
	var list []string
	var apath string
	var i int

	mfi, apath = file.New(), path.Join(asi.rootPath, preffixPath)
	if list, err = mfi.RecursiveFileList(apath); err != nil {
		return asi
	}
	for i = range list {
		rou.Get(path.Join(preffixPath, list[i]), asi.Assets)
	}

	return asi
}

// CleanFilePath Очистка запрашиваемого файла и приведение его к полному пути и имени файла
func (asi *impl) CleanFilePath(requestFile string) (ret string) {
	ret = path.Clean(path.Join(asi.rootPath, requestFile))
	// Если результат содержит путь к базовой директории, то всё корректно
	if strings.Index(ret, path.Join(asi.rootPath, preffixPath)) == 0 {
		return
	}
	// Если не содержит, то возвращаем пустую строку
	ret = ret[:0]

	return
}

// Assets http.HandleFunc
// GET /assets/*
func (asi *impl) Assets(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var mdi filecache.Data
	var ims time.Time
	var fnm string

	// Имя запрашиваемого файла
	if fnm = asi.CleanFilePath(rq.RequestURI); fnm == "" {
		wr.WriteHeader(status.NotFound)
		return
	}
	mdi, err = asi.Mfc.Load(fnm)
	switch err {
	case asi.Mfc.Errors().ErrNotFound():
		wr.WriteHeader(status.NotFound)
		return
	default:
		if err != nil {
			log.Errorf("filecache.Load(%q) error: %s", fnm, err)
			wr.WriteHeader(status.InternalServerError)
			return
		}
	}
	// Если клиент поддерживает If-Modified-Since, загружаем заголовки
	ims, err = time.Parse(ifModifiedSinceTimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && mdi.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Выдача полного комплекта заголовков
	if mdi.ContentType() != "" {
		wr.Header().Set(header.ContentType, mdi.ContentType())
	}
	wr.Header().Set(header.ContentLength, fmt.Sprintf("%d", mdi.Length()))
	wr.Header().Set(header.LastModified, mdi.ModTime().UTC().Format(ifModifiedSinceTimeFormat))
	wr.WriteHeader(status.Ok)
	if _, err = io.Copy(wr, mdi.Reader()); err != nil {
		log.Errorf("Assets response error: %s", err)
	}
}
