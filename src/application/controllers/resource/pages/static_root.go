package pages // import "application/controllers/resource/pages"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"application/models/filecache"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/route"
	"gopkg.in/webnice/web.v1/status"
)

// SetRoutingStaticRootFiles Установка роутинга к статическим файлам
// размещённым в корне DocumentRoot
func (pgi *impl) SetRoutingStaticRootFiles(rou route.Interface) (count uint64) {
	const preffixPath = `/`
	var err error
	var arr []string
	var i, n int
	var next bool

	arr, err = pgi.Mfi.RecursiveFileList(pgi.rootPath)
	if err != nil {
		return
	}
	for i = range arr {
		// Только файлы в корне
		if !rexRootFiles.MatchString(arr[i]) {
			continue
		}
		next = false
		for n = range excludeStaticRootFiles {
			if strings.EqualFold(arr[i], excludeStaticRootFiles[n]) {
				next = true
			}
		}
		// Пропуск технических файлов
		if next {
			continue
		}
		rou.Get(preffixPath+arr[i], pgi.staticRootFiles)
		count++
	}

	return
}

// CleanFilePath Очистка запрашиваемого файла и приведение его к полному пути и имени файла
func (pgi *impl) CleanFilePath(requestFile string) (ret string) {
	const preffixPath = `/`
	ret = path.Clean(path.Join(pgi.rootPath, requestFile))
	// Если результат содержит путь к базовой директории, то всё корректно
	if strings.Index(ret, path.Join(pgi.rootPath, preffixPath)) == 0 {
		return
	}
	// Если не содержит, то возвращаем пустую строку
	ret = ret[:0]

	return
}

// HTTP handler
func (pgi *impl) staticRootFiles(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var mdi filecache.Data
	var ims time.Time
	var fnm string

	// Имя запрашиваемого файла
	if fnm = pgi.CleanFilePath(rq.RequestURI); fnm == "" {
		wr.WriteHeader(status.NotFound)
		return
	}
	switch mdi, err = pgi.Mfc.Load(fnm); err {
	case pgi.Mfc.Errors().ErrNotFound():
		wr.WriteHeader(status.NotFound)
		return
	default:
		if err != nil {
			log.Errorf("file cache load(%q) error: %s", fnm, err)
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
		log.Errorf("staticRootFiles response error: %s", err)
	}
}
