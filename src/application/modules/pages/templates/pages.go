package templates // include "application/modules/pages/templates"

//import "gopkg.in/webnice/web.v1/header"
//import "gopkg.in/webnice/web.v1/status"
//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/route"
import "gopkg.in/webnice/web.v1/method"
import (
	"container/list"
	"net/http"
)

func init() {
	singleton = new(impl)
	singleton.Tpls = new(tpls)
	singleton.Controllers = list.New()
	singleton.Urn = list.New()
}

// Register Регистрация интерфейса контроллера
func Register(ctl Controller) {
	singleton.Controllers.PushBack(ctl)
}

// Get Return interface to singleton object of module
func Get() Interface { return singleton }

// Error Последняя ошибка
func (pgs *impl) Error() error { return pgs.Err }

// Template Интерфейс ко всем загруженным шаблонам файлов
func (pgs *impl) Template() Template { return pgs.Tpls }

// MergeUrnMap URN шаблонов мержатся с URN контроллеров
func (pgs *impl) MergeUrnMap(um *UrnMap) {
	var item *list.Element
	var elm *UrnMap
	var found bool
	var i, n int

	// Поиск уже добавленного URN
	for item = pgs.Urn.Front(); item != nil; item = item.Next() {
		elm = item.Value.(*UrnMap)
		if elm.Urn == um.Urn && elm.Method == um.Method {
			found = true
			break
		}
	}

	// Если не нашелся, добавляем новый
	if !found {
		pgs.Urn.PushBack(um)
		return
	}

	// Если нашелся, мержим
	if found {
		elm.Hndl = append(elm.Hndl, um.Hndl...)
		for i = range um.Tpls.Tpl {
			found = false
			for n = range elm.Tpls.Tpl {
				if um.Tpls.Tpl[i].UrnAddress == elm.Tpls.Tpl[n].UrnAddress {
					found = true
				}
			}
			if !found {
				elm.Tpls.Tpl = append(elm.Tpls.Tpl, um.Tpls.Tpl[i])
			}
		}
	}
}

// Init Инициализация
func (pgs *impl) Init() (ret Interface) {
	var item *list.Element
	var elm Controller
	var riders []*Rider
	var um *UrnMap
	var i int

	ret = pgs

	// TODO
	// Настроить
	//pgs.PagesRoot
	//pgs.Debug

	// Инициализация зарегистрированных контроллеров
	for item = pgs.Controllers.Front(); item != nil; item = item.Next() {
		var rdrs []*Rider
		elm = item.Value.(Controller)
		if rdrs, pgs.Err = elm.Init(); pgs.Err != nil {
			return
		}
		riders = append(riders, rdrs...)
	}
	// Сообщение о завершении инициализации контроллеров html страниц
	log.Info("Controllers of html pages has initialized successfully")

	// Загрузка html страниц
	if pgs.LoadPages(); pgs.Err != nil {
		return
	}

	// Добавление всех шаблонов в карту URN
	for i = range pgs.Tpls.Tpl {
		um = new(UrnMap)
		um.Urn = pgs.Tpls.Tpl[i].UrnAddress
		um.Method = method.Get
		um.Tpls = new(tpls)
		um.Tpls.Tpl = append(um.Tpls.Tpl, pgs.Tpls.Tpl[i])
		pgs.Urn.PushBack(um)
	}

	// Добавление контроллеров в карту URN
	for i = range riders {
		um = new(UrnMap)
		um.Urn = pgs.urnNormal(riders[i].Hundler.Urn)
		um.Method = riders[i].Hundler.Method
		um.Hndl = append(um.Hndl, &riders[i].Hundler)
		// Поиск запрошенных шаблонов
		if um.Tpls, pgs.Err = pgs.getTemplates(riders[i].Templates); pgs.Err != nil {
			// Ошибка отсутствия шаблона
			return
		}
		pgs.MergeUrnMap(um)
	}

	// DEBUG
	//	for item = pgs.Urn.Front(); item != nil; item = item.Next() {
	//		elm := item.Value.(*UrnMap)
	//		debug.Dumper(elm)
	//	}
	// DEBUG

	return
}

// Routing Выгрузка настроек роутинга
func (pgs *impl) Routing() (ret http.Handler) {
	var rou route.Interface
	var item *list.Element

	rou = route.New()
	for item = pgs.Urn.Front(); item != nil; item = item.Next() {
		var elm = item.Value.(*UrnMap)
		switch elm.Method {
		case method.Get:
			//			irs.Get(elm.Urn, func(ctx *iris.Context) { pgs.irisHandlerFunc(elm, ctx) })
		case method.Post:
			//			irs.Post(elm.Urn, func(ctx *iris.Context) { pgs.irisHandlerFunc(elm, ctx) })
		case method.Put:
			//			irs.Put(elm.Urn, func(ctx *iris.Context) { pgs.irisHandlerFunc(elm, ctx) })
		case method.Patch:
			//			irs.Patch(elm.Urn, func(ctx *iris.Context) { pgs.irisHandlerFunc(elm, ctx) })
		case method.Delete:
			//			irs.Delete(elm.Urn, func(ctx *iris.Context) { pgs.irisHandlerFunc(elm, ctx) })
		case method.Head:
			//			irs.Head(elm.Urn, func(ctx *iris.Context) { pgs.irisHandlerFunc(elm, ctx) })
		case method.Options:
			//			irs.Options(elm.Urn, func(ctx *iris.Context) { pgs.irisHandlerFunc(elm, ctx) })
		default:
			log.Warningf("Ignore method '%s' for URN '%s'", elm.Method, elm.Urn)
		}
	}

	return rou
}

// irisHandlerFunc Функция-обработчик стандартного запроса фреймворка iris
func (pgs *impl) irisHandlerFunc(elm *UrnMap) {
	//	var err error
	//	var i int
	//	var key string
	//	var isBreak bool
	//
	//	// Обновление данных файлов шаблонов
	//	for i = range elm.Tpls.Tpl {
	//		for key = range elm.Tpls.Tpl[i].MapData {
	//			if err = pgs.LoadFileBodyData(elm.Tpls.Tpl[i].MapData[key]); err != nil {
	//				log.Errorf("Error LoadFileBodyData: %s", err.Error())
	//				ctx.SetStatusCode(status.InternalServerError)
	//				return
	//			}
	//		}
	//	}
	//
	//	// Если нет контроллеров, то выдаём контент как есть
	//	if len(elm.Hndl) == 0 {
	//		for i = range elm.Tpls.Tpl {
	//			for key = range elm.Tpls.Tpl[i].MapData {
	//				if key == "index" || len(elm.Tpls.Tpl) == 1 {
	//					if elm.Tpls.Tpl[i].MapData[key].Data == nil {
	//						ctx.SetStatusCode(status.NoContent)
	//						return
	//					}
	//					var tmpRq = new(http.Request)
	//					tmpRq.Header.Set(header.IfModifiedSince, ctx.RequestHeader(header.IfModifiedSince))
	//					if !httputil.IfModifiedSince(tmpRq, elm.Tpls.Tpl[i].MapData[key].ModTime) {
	//						ctx.SetStatusCode(status.NotModified)
	//						return
	//					}
	//					ctx.SetHeader(header.LastModified, elm.Tpls.Tpl[i].MapData[key].ModTime.UTC().Format(http.TimeFormat))
	//					if elm.Tpls.Tpl[i].MapData[key].Data.Len() == 0 {
	//						ctx.SetStatusCode(status.NoContent)
	//						return
	//					}
	//					ctx.HTML(status.Ok, elm.Tpls.Tpl[i].MapData[key].Data.String())
	//					return
	//				}
	//			}
	//		}
	//		return
	//	}
	//
	//	// Запуск контроллеров по очереди
	//	// Контроллеры фильтруются в соответствии с методом запроса
	//	isBreak = false
	//	for i = range elm.Hndl {
	//		if elm.Hndl[i].Method.String() != ctx.Method() {
	//			continue
	//		}
	//		if elm.Hndl[i].Func == nil {
	//			log.Warningf("Pages template controller func is null! Method '%s', URN: '%s'", ctx.Method(), elm.Urn)
	//			continue
	//		}
	//		// Если контролер вернул true, то выполнение других контроллеров пропускается
	//		if isBreak {
	//			continue
	//		}
	//		isBreak = elm.Hndl[i].Func(elm.Tpls, ctx)
	//	}
}
