package sitemap // import "application/controllers/resource/sitemap"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"application/controllers/resource/pool"
	modelsSitemap "application/models/sitemap"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/status"
)

// SitemapHandlerFunc http.HandlerFunc получения запросов к файлам sitemap и sitemap-index
func (smi *impl) SitemapHandlerFunc(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var tmp []string
	var tpe modelsSitemap.Type
	var idx, from, to uint64

	// Разбор запроса
	if tmp = rexSitemap.FindStringSubmatch(rq.URL.Path); len(tmp) != 4 {
		wr.WriteHeader(status.NotFound)
		return
	}
	// Получение типа файла
	switch strings.ToLower(tmp[1]) {
	case modelsSitemap.Resource.String():
		tpe = modelsSitemap.Resource
	case modelsSitemap.ResourceIndex.String():
		tpe = modelsSitemap.ResourceIndex
	default:
		wr.WriteHeader(status.NotFound)
		return
	}
	// Если четвёртое значение равно 0 - передан индекс в третьем значении
	if tmp[3] == "0" || tmp[3] == "" {
		if idx, err = strconv.ParseUint(tmp[2], 10, 64); err != nil {
			idx = 0
		}
		if idx >= modelsSitemap.MaxIndex {
			wr.WriteHeader(status.RequestEntityTooLarge)
			return
		}
		from, to = idx*modelsSitemap.MaxRecords, idx*modelsSitemap.MaxRecords+modelsSitemap.MaxRecords-1
	} else {
		// Если четвёртое значение не равно 0 - передан диапазон от и до
		if from, err = strconv.ParseUint(tmp[2], 10, 64); err != nil {
			from = 0
		}
		if to, err = strconv.ParseUint(tmp[3], 10, 64); err != nil {
			to = 0
		}
		if from > to {
			wr.WriteHeader(status.BadRequest)
			return
		}
		if to-from > modelsSitemap.MaxRecords-1 {
			wr.WriteHeader(status.RequestEntityTooLarge)
			return
		}
	}
	// Проверка получившегося диапазона
	if from >= modelsSitemap.MaxRecords*modelsSitemap.MaxIndex || to > modelsSitemap.MaxRecords*modelsSitemap.MaxIndex {
		wr.WriteHeader(status.RequestEntityTooLarge)
		return
	}
	smi.sitemapHandlerFunc(wr, rq, tpe, from, to)
}

// Обработка запроса к sitemap и sitemap-index
func (smi *impl) sitemapHandlerFunc(wr http.ResponseWriter, rq *http.Request, tpe modelsSitemap.Type, from uint64, to uint64) {
	var err error
	var ims time.Time
	var count uint64
	var content pool.ByteBufferInterface

	// Если клиент поддерживает If-Modified-Since, загружаем заголовки
	ims, err = time.Parse(ifModifiedSinceTimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && smi.Smm.URLSetAt().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Получение и возврат интерфейса ByteBuffer
	content = smi.Pool.ByteBufferGet()
	defer smi.Pool.ByteBufferPut(content)
	count = to - from
	switch tpe {
	case modelsSitemap.Resource:
		err = smi.Smm.SitemapXMLWriteTo(content, smi.serverURL, from, count)
	case modelsSitemap.ResourceIndex:
		err = smi.Smm.SitemapIndexXMLWriteTo(content, smi.serverURL, from/modelsSitemap.MaxRecords)
	}
	switch err {
	case smi.Smm.Errors().ErrNotFound():
		wr.WriteHeader(status.NotFound)
		return
	case smi.Smm.Errors().ErrTooLarge():
		wr.WriteHeader(status.RequestEntityTooLarge)
		return
	default:
		if err != nil {
			log.Errorf("model error: %s", err)
			wr.WriteHeader(status.InternalServerError)
			return
		}
	}
	wr.Header().Set(header.ContentType, mime.ApplicationXMLCharsetUTF8)
	wr.Header().Set(header.LastModified, smi.Smm.URLSetAt().UTC().Format(ifModifiedSinceTimeFormat))
	wr.Header().Set(header.ContentLength, strconv.Itoa(content.Len()))
	wr.WriteHeader(status.Ok)
	if _, err = content.WriteTo(wr); err != nil {
		log.Errorf("response error: %s", err)
		return
	}
}
