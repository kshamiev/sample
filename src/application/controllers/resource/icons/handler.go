package icons // import "application/controllers/resource/icons"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"application/controllers/resource/pool"
	"application/models/filecache"
	"application/modules/rendering"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/status"
)

// Разбор запроса
func (ici *impl) parseRequest(req string) (ret *size) {
	var tmp []string
	var w, h uint64
	var it, tp string
	var i int
	var sz []size

	if tmp = rexIconSizeAndType.FindStringSubmatch(req); len(tmp) != 5 {
		return
	}
	w, _ = strconv.ParseUint(tmp[2], 10, 64) // nolint: gosec
	h, _ = strconv.ParseUint(tmp[3], 10, 64) // nolint: gosec
	it, tp = strings.ToLower(tmp[1]), strings.ToLower(tmp[4])
	switch it {
	case keyFavicon:
		sz = sizeFavicon
	case keyAppleTouchIcon:
		sz = sizeAppleTouchIcon
	case keyAndroidChrome:
		sz = sizeAndroidChrome
	case keyMstile:
		sz = sizeMstile
	default:
		return
	}
	for i = range sz {
		if uint(w) == sz[i].Width && uint(h) == sz[i].Height && tp == sz[i].Ext {
			ret = &sz[i]
		}
	}

	return
}

// Обработчик http запросов веб сервера
func (ici *impl) iconHandlerFunc(wr http.ResponseWriter, rq *http.Request, grp *group) {
	var err error
	var mdi filecache.Data
	var ims time.Time

	// Загрузка графического изображения
	mdi, err = ici.Mfc.Load(rq.URL.Path)
	switch err {
	case ici.Mfc.Errors().ErrNotFound():
		wr.WriteHeader(status.NotFound)
		return
	default:
		if err != nil {
			log.Errorf("Load file for request %q error: %s", rq.URL.Path, err)
			wr.WriteHeader(status.InternalServerError)
			return
		}
	}
	// Если клиент поддерживает If-Modified-Since, загружаем заголовки, проверяем
	ims, err = time.Parse(ifModifiedSinceTimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && mdi.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Выдача полного комплекта заголовков
	wr.Header().Set(header.ContentType, mdi.ContentType())
	wr.Header().Set(header.ContentLength, fmt.Sprintf("%d", mdi.Size()))
	wr.Header().Set(header.LastModified, mdi.ModTime().UTC().Format(ifModifiedSinceTimeFormat))
	wr.WriteHeader(status.Ok)
	if _, err = io.Copy(wr, mdi.Reader()); err != nil {
		log.Errorf("Assets response error: %s", err)
	}
}

// Функция http.HandleFunc формирующая http страницу с предпросмотром всех созданных иконок
func (ici *impl) iconsPreviewPageHandlerFunc(wr http.ResponseWriter, rq *http.Request) {
	const keyTemplate = `{{ template "preview" . }}`
	var err error
	var tpln string
	var mdi filecache.Data
	var rend rendering.Interface
	var tpl, content pool.ByteBufferInterface
	var vars *htmlPreviewTemplateVars
	var i int

	// Загрузка шаблона
	tpln = path.Join(ici.rootPath, htmlPreviewTemplate)
	if mdi, err = ici.Mfc.Load(tpln); err != nil {
		log.Errorf("filecache.Load(%q) error: %s", tpln, err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Получение и возврат интерфейса ByteBuffer
	tpl = ici.Pool.ByteBufferGet()
	defer ici.Pool.ByteBufferPut(tpl)
	_, _ = fmt.Fprint(tpl, keyTemplate)
	content = ici.Pool.ByteBufferGet()
	defer ici.Pool.ByteBufferPut(content)
	// Переменные шаблона
	vars = &htmlPreviewTemplateVars{Groups: make([]*htmlPreviewTemplateVarsGroup, 0, len(ici.groups))}
	for i = range ici.groups {
		vars.Groups = append(vars.Groups, &htmlPreviewTemplateVarsGroup{
			Group: ici.groups[i].name,
			Items: ici.groups[i].items,
		})
	}
	// Подготовка шаблонизатора и данных
	rend = rendering.New(rendering.Option{
		Directory: ici.rootPath,
	})
	err = rend.NewStandardTemplateRender().
		RenderTextData(content, vars, mdi.Reader(), tpl)
	if err != nil {
		log.Errorf("Template %s error: %s", htmlPreviewTemplate, err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Set(header.ContentType, mime.TextHTMLCharsetUTF8)
	wr.WriteHeader(status.Ok)
	if _, err = content.WriteTo(wr); err != nil {
		log.Errorf("Response error: %s", err)
		return
	}
}

// Функция http.HandleFunc формирующая manifest.json с описанием иконок для android сhrome
func (ici *impl) iconsManifestHandlerFunc(wr http.ResponseWriter, rq *http.Request) {
	const keyTemplate = `{{ template "manifest" . }}`
	var err error
	var tpln string
	var mdi filecache.Data
	var rend rendering.Interface
	var tpl, content, jsn pool.ByteBufferInterface
	var vars *manifestTemplateVars
	var icons *manifestIcons
	var enc *json.Encoder
	var tmp [][][]byte
	var i int

	// Загрузка шаблона
	tpln = path.Join(ici.rootPath, htmlPreviewTemplate)
	if mdi, err = ici.Mfc.Load(tpln); err != nil {
		log.Errorf("filecache.Load(%q) error: %s", tpln, err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Получение и возврат интерфейса ByteBuffer
	tpl = ici.Pool.ByteBufferGet()
	defer ici.Pool.ByteBufferPut(tpl)
	_, _ = fmt.Fprint(tpl, keyTemplate)
	content = ici.Pool.ByteBufferGet()
	defer ici.Pool.ByteBufferPut(content)
	jsn = ici.Pool.ByteBufferGet()
	defer ici.Pool.ByteBufferPut(jsn)
	// Формирование массива иконок в виде структуры
	icons = &manifestIcons{Icons: make([]*manifestIconsItems, 0, len(sizeAndroidChrome))}
	for i = range sizeAndroidChrome {
		icons.Icons = append(icons.Icons, &manifestIconsItems{
			Src:   ici.androidChromeImageName(&sizeAndroidChrome[i]),
			Sizes: fmt.Sprintf("%dx%d", sizeAndroidChrome[i].Width, sizeAndroidChrome[i].Height),
			Type:  sizeAndroidChrome[i].ContentType,
		})
	}
	enc = json.NewEncoder(jsn)
	enc.SetEscapeHTML(true)
	enc.SetIndent("", "  ")
	if err = enc.Encode(icons); err != nil {
		log.Errorf("Json encoder error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Подготовка данных для шаблонизатора
	tmp = rexJSONInside.FindAllSubmatch(jsn.Bytes(), -1)
	if len(tmp) == 1 && len(tmp[0]) == 4 {
		ici.Pool.ByteBufferPut(jsn)
		jsn = ici.Pool.ByteBufferGet()
		_, _ = fmt.Fprint(jsn, string(tmp[0][2]))
	}
	// Переменные шаблона
	vars = &manifestTemplateVars{Icons: string(jsn.Bytes())}
	// Подготовка шаблонизатора и данных
	rend = rendering.New(rendering.Option{
		Directory: ici.rootPath,
	})
	err = rend.NewStandardTemplateRender().
		RenderTextData(content, vars, mdi.Reader(), tpl)
	if err != nil {
		log.Errorf("Template %s error: %s", htmlPreviewTemplate, err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Set(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	wr.WriteHeader(status.Ok)
	if _, err = content.WriteTo(wr); err != nil {
		log.Errorf("Response error: %s", err)
		return
	}

}

// Функция http.HandleFunc формирующая browserconfig.xml с описанием иконок для microsoft ишак
func (ici *impl) iconsBrowserconfigHandlerFunc(wr http.ResponseWriter, rq *http.Request) {
	const keyTemplate = `{{ template "browserconfig" . }}`
	var err error
	var tpln string
	var mdi filecache.Data
	var rend rendering.Interface
	var tpl, content pool.ByteBufferInterface
	var vars *browserconfigTemplateVars
	var i int

	// Загрузка шаблона
	tpln = path.Join(ici.rootPath, htmlPreviewTemplate)
	if mdi, err = ici.Mfc.Load(tpln); err != nil {
		log.Errorf("filecache.Load(%q) error: %s", tpln, err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Получение и возврат интерфейса ByteBuffer
	tpl = ici.Pool.ByteBufferGet()
	defer ici.Pool.ByteBufferPut(tpl)
	_, _ = fmt.Fprint(tpl, keyTemplate)
	content = ici.Pool.ByteBufferGet()
	defer ici.Pool.ByteBufferPut(content)
	// Переменные шаблона
	vars = &browserconfigTemplateVars{Icons: make([]*browserconfigIconItem, 0, len(sizeMstile))}
	for i = range sizeMstile {
		vars.Icons = append(vars.Icons, &browserconfigIconItem{
			Width:       sizeMstile[i].Width,
			Height:      sizeMstile[i].Height,
			Ext:         sizeMstile[i].Ext,
			ContentType: sizeMstile[i].ContentType,
			Wide:        sizeMstile[i].Wide,
		})
	}
	// Подготовка шаблонизатора и данных
	rend = rendering.New(rendering.Option{
		Directory: ici.rootPath,
	})
	err = rend.NewStandardTemplateRender().
		RenderTextData(content, vars, mdi.Reader(), tpl)
	if err != nil {
		log.Errorf("Template %s error: %s", htmlPreviewTemplate, err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Set(header.ContentType, mime.ApplicationXMLCharsetUTF8)
	wr.WriteHeader(status.Ok)
	if _, err = content.WriteTo(wr); err != nil {
		log.Errorf("Response error: %s", err)
		return
	}
}

// Функция принимающая HTTP запрос:
// GET /safari-pinned-tab.svg
func (ici *impl) appleSafariPinnedTabVectorHandlerFunc(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var mdi filecache.Data
	var ims time.Time

	// Загрузка исходного файла
	if mdi, err = ici.Mfc.Load(path.Join(ici.rootPath, imgVector)); err != nil {
		log.Errorf("Source file error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Если клиент поддерживает If-Modified-Since, загружаем заголовки, проверяем
	ims, err = time.Parse(ifModifiedSinceTimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && mdi.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Выдача полного комплекта заголовков
	wr.Header().Set(header.ContentType, mdi.ContentType())
	wr.Header().Set(header.ContentLength, fmt.Sprintf("%d", mdi.Size()))
	wr.Header().Set(header.LastModified, mdi.ModTime().UTC().Format(ifModifiedSinceTimeFormat))
	wr.WriteHeader(status.Ok)
	if _, err = io.Copy(wr, mdi.Reader()); err != nil {
		log.Errorf("Assets response error: %s", err)
	}
}
