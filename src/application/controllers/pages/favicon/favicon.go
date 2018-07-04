package favicon // import "application/controllers/pages/favicon"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import (
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"application/configuration"
	"application/modules/images"
)

// New creates a new object and return interface
func New() Interface {
	var fin = new(impl)
	return fin
}

// Lazy initialization
func (fin *impl) init() {
	if fin.cfg == nil {
		fin.cfg = configuration.Get()
	}
	if fin.DocumentRoot == "" {
		// first web server
		for i := range fin.cfg.Configuration().WEBServers {
			fin.DocumentRoot = fin.cfg.Configuration().WEBServers[i].DocumentRoot
			break
		}
	}
}

// FavIcon Cоздания на лету favicon.ico и всех возможных вариантов размеров и форматов favicon
//  <link rel="icon" type="image/vnd.microsoft.icon" href="/favicon.ico" />
//  <link rel="icon" type="image/png" href="/favicon.png" />
// Поиск в папке DocumentRoot файла favicon.* являющегося картинкой, конвертация в icon, ресайз при необходимости и вывод
func (fin *impl) FavIcon(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var im images.Interface
	var imgs []*images.Image
	var resp *images.Image
	var ims time.Time

	fin.init()
	im = images.New()
	// Чтение из папки DocumentRoot всех файлов являющихся картинкой и удовлетворяющих патерну
	imgs, err = im.FindInFolderAndLoad(path.Join(fin.DocumentRoot, "assets"), rexFaviconName)
	if err != nil {
		log.Errorf("FavIcon load image error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Выбор картинки и ресайз
	resp = fin.faviconSelectAndResize(rq, imgs)
	// Если картинка не изменилась и был передан заголовок проверки модификации файла
	ims, err = time.Parse(http.TimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && resp.FileInfo.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Вывод результата
	wr.Header().Set(header.LastModified, resp.FileInfo.ModTime().UTC().Format(http.TimeFormat))

	switch resp.Type {
	case _Ico:
		wr.Header().Set(header.ContentType, mime.FavIcon)
		wr.WriteHeader(status.Ok)
		err = im.WriteIco(wr, resp)
	case _Png:
		wr.Header().Set(header.ContentType, mime.ImagePNG)
		wr.WriteHeader(status.Ok)
		err = im.WritePng(wr, resp)
	}
	if err != nil {
		log.Errorf("Response error: %s", err.Error())
		return
	}
}

// Выбор картинки и ресайз
func (fin *impl) faviconSelectAndResize(rq *http.Request, imgs []*images.Image) (ret *images.Image) {
	var format string
	var width uint
	var height uint
	var i int
	var im = images.New()

	format, width, height = fin.faviconDetectSize(rq)

	// Приоритетный результат
	// Если есть ico, то выводится он
	// Если ico нет, то выводится первая картинка конвертируемая в ico
	for i = range imgs {
		if imgs[i].Type == format {
			ret = imgs[i]
		}
		if ret == nil {
			ret = imgs[i]
		}
	}

	// Resize
	if ret.Config.Width != int(width) || ret.Config.Height != int(height) {
		ret = im.Resize(ret, width, height)
	}

	return
}

// Определение на основе запроса размеров и формата требуемого favicon.ico
func (fin *impl) faviconDetectSize(rq *http.Request) (format string, width uint, height uint) {
	var ssm []string
	var tmp uint64

	// Defaults
	format, width, height = _Ico, 32, 32

	// Если нет совпадения, возвращаем дефолтовые значения
	if ssm = rexFaviconSizeAndType.FindStringSubmatch(rq.URL.Path); len(ssm) != 4 {
		return
	}

	// Совпадение есть, проверяем значения
	format = ssm[3]
	tmp, _ = strconv.ParseUint(ssm[1], 0, 64)
	width = uint(tmp)
	tmp, _ = strconv.ParseUint(ssm[2], 0, 64)
	height = uint(tmp)

	// Формат
	switch strings.ToLower(format) {
	case _Ico, _Png:
		format = strings.ToLower(format)
	default:
		format = _Ico
	}

	// Размер
	// + не меньше 16x!6
	// + не больше 460x460
	if width > 460 || height > 460 {
		width, height = 460, 460
	}
	if width < 16 {
		width, height = 16, 16
	}

	return
}
