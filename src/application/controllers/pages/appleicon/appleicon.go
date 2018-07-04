package appleicon // import "application/controllers/pages/appleicon"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import (
	"net/http"
	"path"
	"strconv"
	"time"

	"application/configuration"
	"application/modules/images"
)

// New creates a new object and return interface
func New() Interface {
	var ati = new(impl)
	return ati
}

// Lazy initialization
func (ati *impl) init() {
	if ati.cfg == nil {
		ati.cfg = configuration.Get()
	}
	if ati.DocumentRoot == "" {
		// first web server
		for i := range ati.cfg.Configuration().WEBServers {
			ati.DocumentRoot = ati.cfg.Configuration().WEBServers[i].DocumentRoot
			break
		}
	}
}

// AppleIcon Specifying a Webpage Icon for Web Clip
// Поиск в папке DocumentRoot файла favicon.* являющегося картинкой, конвертация в png, ресайз при необходимости и вывод
func (ati *impl) AppleIcon(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var im images.Interface
	var imgs []*images.Image
	var resp *images.Image
	var ims time.Time

	ati.init()
	im = images.New()
	// Чтение из папки DocumentRoot всех файлов являющихся картинкой и удовлетворяющих патерну
	imgs, err = im.FindInFolderAndLoad(path.Join(ati.DocumentRoot, "assets"), rexAppleIconName)
	if err != nil {
		log.Errorf("FavIcon load image error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Выбор картинки и ресайз
	resp = ati.appleTouchIconSelectAndResize(rq, imgs)
	// Если картинка не изменилась и был передан заголовок проверки модификации файла
	ims, err = time.Parse(http.TimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && resp.FileInfo.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Вывод результата
	wr.Header().Set(header.LastModified, resp.FileInfo.ModTime().UTC().Format(http.TimeFormat))
	wr.Header().Set(header.ContentType, mime.ImagePNG)
	wr.WriteHeader(status.Ok)
	if err = im.WritePng(wr, resp); err != nil {
		log.Errorf("Response error: %s", err.Error())
		return
	}
}

// Выбор картинки и ресайз
func (ati *impl) appleTouchIconSelectAndResize(rq *http.Request, imgs []*images.Image) (ret *images.Image) {
	var width uint
	var height uint
	var i int
	var im = images.New()

	width, height = ati.appleTouchIconDetectSize(rq)

	// Приоритетный результат
	// Если есть png, то выводится он
	// Если png нет, то выводится первая картинка конвертируемая в png
	for i = range imgs {
		if imgs[i].Type == _Png {
			ret = imgs[i]
		}
		if ret == nil {
			ret = imgs[0]
		}
	}

	// Resize
	if ret.Config.Width != int(width) || ret.Config.Height != int(height) {
		ret = im.Resize(ret, width, height)
	}

	return
}

// Определение на основе запроса размеров и формата требуемого favicon.ico
func (ati *impl) appleTouchIconDetectSize(rq *http.Request) (width uint, height uint) {
	var ssm []string
	var tmp uint64

	// Defaults
	width, height = 60, 60

	// Если нет совпадения, возвращаем дефолтовые значения
	if ssm = rexAppleIconSizeAndType.FindStringSubmatch(rq.URL.Path); len(ssm) != 5 {
		return
	}

	// Совпадение есть, проверяем значения
	tmp, _ = strconv.ParseUint(ssm[2], 0, 64)
	width = uint(tmp)
	tmp, _ = strconv.ParseUint(ssm[3], 0, 64)
	height = uint(tmp)

	// Размер
	// + не меньше 60x60
	// + не больше 2048x2048
	if width > 2048 || height > 2048 {
		width, height = 2048, 2048
	}
	if width < 60 {
		width, height = 60, 60
	}

	return
}
