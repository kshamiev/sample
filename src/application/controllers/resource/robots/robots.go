package robots // import "application/controllers/resource/robots"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"application/controllers/resource/pool"
	"application/controllers/resource/sitemap"
	"application/models/filecache"
	"application/modules/rendering"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/route"
	"gopkg.in/webnice/web.v1/status"
)

// New creates a new object and return interface
func New() Interface {
	var rbi = &impl{
		Mfc:  filecache.Get(),
		Pool: pool.New(),
	}
	return rbi
}

// Debug Set debug mode
func (rbi *impl) Debug(d bool) Interface {
	rbi.debug = d
	rbi.Pool.Debug(rbi.debug)
	rbi.Mfc.Debug(rbi.debug)
	return rbi
}

// DocumentRoot Устанавливает путь к корню веб сервера
func (rbi *impl) DocumentRoot(path string) Interface { rbi.rootPath = path; return rbi }

// ServerURL Устанавливает основной адрес веб сервера
func (rbi *impl) ServerURL(u string) Interface {
	var err error
	var su *url.URL

	rbi.serverURL = u
	if su, err = url.ParseRequestURI(u); err != nil {
		log.Criticalf("Parse server URL error: %s", err)
		return rbi
	}
	rbi.serverScheme, rbi.serverDomain = su.Scheme, su.Host

	return rbi
}

// Sitemap Устанавливает интерфейс sitemap
func (rbi *impl) Sitemap(smi sitemap.Interface) Interface {
	rbi.Smi = smi
	return rbi
}

// SetRouting Установка роутинга к статическим файлам
func (rbi *impl) SetRouting(rou route.Interface) Interface {
	rou.Get(keyRobots, rbi.RobotsTxt)
	return rbi
}

// Получение протокола на котором работает сервер
func (rbi *impl) getProto(rq *http.Request) (ret string) {
	defer func() {
		switch ret {
		case "http", "https":
			return
		default:
			ret = "http"
		}
	}()
	if ret = strings.ToLower(rq.Header.Get(header.XScheme)); ret != "" {
		return
	}
	if ret = strings.ToLower(rq.Header.Get(header.XForwardedProto)); ret != "" {
		return
	}
	if rq.TLS != nil {
		ret = `https`
	} else {
		ret = `http`
	}

	return
}

// RobotsTxt http.HandleFunc
// GET /robots.txt
func (rbi *impl) RobotsTxt(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var vars *templateVars
	var tpln string
	var tmpm []string
	var mdi filecache.Data
	var content pool.ByteBufferInterface
	var rend rendering.Interface
	var ims time.Time
	var i int

	tpln = path.Clean(path.Join(rbi.rootPath, keyRobotsTemplate))
	if strings.Index(tpln, rbi.rootPath) != 0 {
		log.Errorf("Pattern search outside of folder DocumentRoot: %q", tpln)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Загрузка шаблона
	if mdi, err = rbi.Mfc.Load(tpln); err != nil {
		log.Errorf("filecache.Load(%q) error: %s", tpln, err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	// Если клиент поддерживает If-Modified-Since, загружаем заголовки
	ims, err = time.Parse(ifModifiedSinceTimeFormat, rq.Header.Get(header.IfModifiedSince))
	if err == nil && !ims.IsZero() && mdi.ModTime().Before(ims.Add(time.Second*1)) {
		wr.WriteHeader(status.NotModified)
		return
	}
	// Получение и возврат интерфейса ByteBuffer
	content = rbi.Pool.ByteBufferGet()
	defer rbi.Pool.ByteBufferPut(content)
	// Подготовка шаблонизатора и данных
	vars = &templateVars{
		RequestScheme: rbi.getProto(rq),
		RequestDomain: rq.Host,
		ServerURL:     rbi.serverURL,
		ServerScheme:  rbi.serverScheme,
		ServerDomain:  rbi.serverDomain,
	}
	if rbi.Smi != nil {
		tmpm = rbi.Smi.Links()
		vars.Sitemap = make([]tvSitemap, 0, len(tmpm))
		for i = range tmpm {
			vars.Sitemap = append(vars.Sitemap, tvSitemap{
				URN: tmpm[i],
			})
		}
	}
	rend = rendering.New(rendering.Option{
		Directory: rbi.rootPath,
	})
	err = rend.NewStandardTemplateRender().
		RenderTextData(content, vars, mdi.Reader())
	if err != nil {
		log.Errorf("robots.txt template error: %s", err)
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Set(header.ContentType, mime.TextPlainCharsetUTF8)
	wr.Header().Set(header.LastModified, mdi.ModTime().UTC().Format(ifModifiedSinceTimeFormat))
	wr.WriteHeader(status.Ok)
	if _, err = content.WriteTo(wr); err != nil {
		log.Errorf("Response error: %s", err)
		return
	}
}
