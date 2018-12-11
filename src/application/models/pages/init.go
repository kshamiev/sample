package pages // import "application/models/pages"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"fmt"
	stddebug "runtime/debug"
	"strings"

	pagesTypes "application/models/pages/types"
)

// Init Инициализация объекта
func (pgm *impl) Init() (err error) {
	var elm *list.Element
	var item pagesTypes.Controller

	// Инициализация зарегистрированных контроллеров
	for elm = pgm.controllers.Front(); elm != nil; elm = elm.Next() {
		item = elm.Value.(pagesTypes.Controller)
		if err = pgm.InitController(item); err != nil {
			err = fmt.Errorf("controller %q init error: %s", packageName(item), err)
			return
		}
	}
	// Инициализация шаблонов
	if err = pgm.InitTemplates(); err != nil {
		err = fmt.Errorf("templates error: %s", err)
		return
	}

	return
}

// InitController Вызов функции инициализации контроллера с защитой от паники
func (pgm *impl) InitController(ctl pagesTypes.Controller) (err error) {
	var mft *pagesTypes.Manifest
	var rsp *pagesTypes.Response
	var i, m, n int
	var ok bool

	// Защита от паники
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Panic recovery:\n%v\n%s", e.(error), string(stddebug.Stack()))
		}
	}()
	// Вызов функции контроллера
	if mft, err = ctl.Init(pgm.serverURL); err != nil {
		err = fmt.Errorf("init error: %s", err)
		return
	}
	// Обработка rider
	pgm.urnSync.Lock()
	defer pgm.urnSync.Unlock()
	for i = range mft.Handlers {
		if _, ok = pgm.urn[mft.Handlers[i].URN]; !ok {
			pgm.urn[mft.Handlers[i].URN] = new(pagesTypes.Response)
		}
		rsp = pgm.urn[mft.Handlers[i].URN].(*pagesTypes.Response)
		// Добавление обработчика
		rsp.Handlers = append(rsp.Handlers, mft.Handlers[i])
		// Добавление URN шаблонов
		for m = range mft.TemplatesURN {
			ok = false
			for n = range rsp.TemplatesURN {
				if strings.EqualFold(mft.TemplatesURN[m], rsp.TemplatesURN[n]) {
					ok = true
				}
			}
			if !ok {
				rsp.TemplatesURN = append(rsp.TemplatesURN, mft.TemplatesURN[m])
			}
		}
	}

	return
}

// InitTemplates Загрузка шаблонов
func (pgm *impl) InitTemplates() (err error) {
	var tpl []*pagesTypes.Body
	var rsp *pagesTypes.Templates
	var i int
	var ok bool

	// Загрузка шаблонов
	if tpl, err = pgm.templateInit(); err != nil {
		return
	}
	// Создание карты URN шаблонов
	pgm.tplSync.Lock()
	defer pgm.tplSync.Unlock()
	for i = range tpl {
		if _, ok = pgm.tpl[tpl[i].Path]; !ok {
			pgm.tpl[tpl[i].Path] = new(pagesTypes.Templates)
		}
		rsp = pgm.tpl[tpl[i].Path].(*pagesTypes.Templates)
		*rsp = append(*rsp, tpl[i])
	}

	return
}
