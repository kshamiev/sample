package pages // import "application/controllers/resource/pages"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"net/url"
	"path"

	"application/controllers/resource/pool"
	"application/models/file"
	"application/models/filecache"
	modelsPages "application/models/pages"

	"gopkg.in/webnice/web.v1/route"
	"gopkg.in/webnice/web.v1/status"
)

// New creates a new object and return interface
func New() Interface {
	var pgi = &impl{
		Mfi:  file.New(),
		Mfc:  filecache.Get(),
		Pool: pool.New(),
		Pgm:  modelsPages.Get(),
	}
	return pgi
}

// Debug Set debug mode
func (pgi *impl) Debug(d bool) Interface {
	pgi.debug = d
	pgi.Pool.Debug(pgi.debug)
	pgi.Pgm.Debug(pgi.debug)
	pgi.Mfc.Debug(pgi.debug)
	return pgi
}

// DocumentRoot Устанавливает путь к корню веб сервера
func (pgi *impl) DocumentRoot(path string) Interface { pgi.rootPath = path; return pgi }

// TemplatePages Устанавливает путь к папке файлов шаблонов
func (pgi *impl) TemplatePages(path string) Interface {
	pgi.templatePages = path
	pgi.Pgm.TemplatePages(pgi.templatePages)
	return pgi
}

// ServerURL Устанавливает основной адрес веб сервера
func (pgi *impl) ServerURL(u string) Interface {
	var err error
	var su *url.URL

	pgi.serverURL = u
	if su, err = url.ParseRequestURI(u); err != nil {
		log.Criticalf("parse server URL error: %s", err)
		return pgi
	}
	pgi.serverScheme, pgi.serverDomain = su.Scheme, su.Host
	pgi.Pgm.ServerURL(pgi.serverURL)

	return pgi
}

// SetRouting Установка роутинга к статическим файлам
func (pgi *impl) SetRouting(rou route.Interface) Interface {
	var err error
	var index string
	var urns []string
	var indexHF http.HandlerFunc
	var content filecache.Data
	var i int
	var total uint64

	if urns, err = pgi.routeLoad(); err != nil {
		log.Criticalf("route settings error: %s", err)
	}
	// Проверка наличия не пустого index.html
	index = path.Join(pgi.rootPath, keyIndexHTML)
	if content, err = pgi.Mfc.Load(index); err == nil && content.Length() > 0 {
		indexHF = pgi.IndexHandlerFunc
	}
	// Инициализация модели pages
	if err = pgi.Pgm.Init(); err != nil {
		log.Criticalf("pages model init error: %s", err)
	}
	// Количество
	i = len(indexURN) + len(urns)
	pgi.indexURN = make([]string, 0, i)
	// Регистрация URN адресов для index
	rou.Group(func(sr route.Interface) {
		// Роутинг к статическим файлам находящимся в корне DocumentRoot
		total += pgi.SetRoutingStaticRootFiles(sr)
		// Index файл и его синонимы
		if indexHF != nil {
			// Роутинг к index
			for i = range indexURN {
				sr.Get(indexURN[i], indexHF)
				pgi.indexURN = append(pgi.indexURN, indexURN[i])
				total++
			}
			// Роутинг к синонимам index
			for i = range urns {
				sr.Get(urns[i], indexHF)
				pgi.indexURN = append(pgi.indexURN, urns[i])
				total++
			}
		}
		// Роутинг к страницам создаваемым на основе шаблонов
		total += pgi.Pgm.SetRouting(sr)
		// Если нет ни одного URN в группе - создание URN заглушки
		if total == 0 {
			log.Noticef("SET I AM A TEAPOT")
			sr.Get("/", func(wr http.ResponseWriter, rq *http.Request) {
				wr.WriteHeader(status.ImATeapot)
				if _, err = wr.Write(status.Bytes(status.ImATeapot)); err != nil {
					log.Errorf("response error: %s", err)
				}
			})
		}
	})

	return pgi
}
