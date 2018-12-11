package pages // import "application/models/pages"

import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"application/models/filecache"
	pagesTypes "application/models/pages/types"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/method"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/route"
	"gopkg.in/webnice/web.v1/status"
)

// Get singleton object and return interface
func Get() Interface {
	if singleton == nil {
		singleton = &impl{
			Mfc:         filecache.Get(),
			controllers: list.New(),
			urn:         make(map[string]pagesTypes.Responder),
			urnSync:     new(sync.RWMutex),
			tpl:         make(map[string]pagesTypes.Templater),
			tplSync:     new(sync.RWMutex),
		}
	}
	return singleton
}

// RegisterController Регистрация контроллера
func RegisterController(ctl pagesTypes.Controller) { Get().RegisterController(ctl) }

// Debug Set debug mode
func (pgm *impl) Debug(d bool) Interface {
	pgm.debug = d
	pgm.Mfc.Debug(pgm.debug)
	return pgm
}

// TemplatePages Путь к корню файлов шаблонов
func (pgm *impl) TemplatePages(p string) Interface { pgm.templatePages = p; return pgm }

// ServerURL Устанавливает основной адрес веб сервера
func (pgm *impl) ServerURL(u string) Interface {
	var err error
	var su *url.URL

	pgm.serverURL = u
	if su, err = url.ParseRequestURI(u); err != nil {
		log.Criticalf("Parse server URL error: %s", err)
		return pgm
	}
	pgm.serverScheme, pgm.serverDomain = su.Scheme, su.Host

	return pgm
}

// RegisterController Регистрация контроллера
func (pgm *impl) RegisterController(ctl pagesTypes.Controller) Interface {
	pgm.controllers.PushBack(ctl)
	return pgm
}

// SetRouting Настройка роутинга
func (pgm *impl) SetRouting(rou route.Interface) (count uint64) {
	var urn string
	var urnCount int
	var hf http.HandlerFunc
	var i int

	pgm.urnSync.RLock()
	pgm.tplSync.RLock()
	defer pgm.urnSync.RUnlock()
	defer pgm.tplSync.RUnlock()

	// Подсчёт количества urn
	for urn = range pgm.urn {
		for i = range pgm.urn[urn].(*pagesTypes.Response).Handlers {
			if !pgm.urn[urn].(*pagesTypes.Response).Handlers[i].Private {
				urnCount++
			}
		}
	}

	// Если нет ничего - выход
	if urnCount+len(pgm.tpl) == 0 {
		return
	}
	// Регистрация URN для контроллеров
	rou.Group(func(sr route.Interface) {
		// Контроллеры
		for urn = range pgm.urn {

			hf = func(wr http.ResponseWriter, rq *http.Request) {

				log.Debugf("* m: %s, urn: %s", rq.Method, rq.RequestURI)
				wr.Header().Add(header.ContentType, mime.TextPlainCharsetUTF8)
				wr.WriteHeader(status.Ok)
				_, _ = fmt.Fprintf(wr, "Method: %q\nURN: %q\n", rq.Method, rq.RequestURI)

			}

			for i = range pgm.urn[urn].(*pagesTypes.Response).Handlers {
				if pgm.urn[urn].(*pagesTypes.Response).Handlers[i].Private {
					continue
				}
				switch pgm.urn[urn].(*pagesTypes.Response).Handlers[i].Method {
				case method.Get:
					rou.Get(urn, hf)
				case method.Connect:
					rou.Connect(urn, hf)
				case method.Head:
					rou.Head(urn, hf)
				case method.Options:
					rou.Options(urn, hf)
				case method.Delete:
					rou.Delete(urn, hf)
				case method.Patch:
					rou.Patch(urn, hf)
				case method.Post:
					rou.Post(urn, hf)
				case method.Put:
					rou.Put(urn, hf)
				case method.Trace:
					rou.Trace(urn, hf)
				}
				log.Debugf("+ routing %s %q", pgm.urn[urn].(*pagesTypes.Response).Handlers[i].Method.String(), urn)
			}
		}

		// Шаблоны без контроллеров
		for urn = range pgm.tpl {

			log.Debug(debug.DumperString(urn, pgm.tpl[urn]))

		}

	})

	log.Debug(debug.DumperString(pgm.urn))

	return
}
