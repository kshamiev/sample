package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/controllers"
	middlewareRetrycheck "application/middleware/retrycheck"
	//middlewareSession "application/middleware/session"

	"gopkg.in/webnice/web.v1/route"
)

// API Настройка роутинга к контроллерам API
func (rt *impl) API() {
	// API v1.0
	rt.Rou.Subroute("/api/v1.0", rt.routingAPI)
}

// API v1.0
// /api/v1.0/*
func (rt *impl) routingAPI(r route.Interface) {
	r.Use(middlewareRetrycheck.Handler)
	r.Subroute("/", rt.routingCore)

	// Сущность myitem
	r.Subroute("/myitem", rt.routingMyItem)
}

// Core
func (rt *impl) routingCore(r route.Interface) {
	r.Get("/ping", controllers.PingController.Ping)
	r.Get("/panic", controllers.PanicController.Panic)
	r.Put("/settings/time", controllers.SettingsControllerTime.Time)
	r.Get("/version", controllers.VersionController.Version)
	r.Subroute("/uploadfile", rt.routingStorageFile)
}

// Сущность myitem
func (rt *impl) routingMyItem(r route.Interface) {
	r.Options("/", controllers.MyitemController.Status)
	r.Post("/", controllers.MyitemController.Create)
	r.Get("/:id", controllers.MyitemController.Load)
	r.Delete("/:id", controllers.MyitemController.Delete)
	r.Get("/version", controllers.VersionController.Version)
}

// Upload file storage
func (rt *impl) routingStorageFile(r route.Interface) {
	r.Group(func(sr route.Interface) {
		//sr.Use(middlewareSession.ContextSessionHandler)            // Поднятие сессии авторизации
		//sr.Use(middlewareSession.GuardByAnyGroupHandler().Handler) // Доступ разрешен если поднялась сессия авторизации
		//sr.Use(middlewareSession.SessionLifeTimeUpdateHandler)     // Продление времени жизни сессии
		sr.Post("/", controllers.UploadFileController.UploadFile)
	})
}
