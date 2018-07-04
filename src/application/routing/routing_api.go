package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/route"
import (
	"application/controllers"
)

// RoutingAPI Настройка роутинга к контроллерам API
func (rt *impl) RoutingAPI() {
	// API v1.0
	rt.rou.Subroute("/api/v1.0", rt.routingApi)
}

// API v1.0
// /api/v1.0/*
func (rt *impl) routingApi(r route.Interface) {
	r.Subroute("/", rt.routingCore)

	// Сущность myitem
	r.Subroute("/myitem", rt.routingMyItem)
}

// Core
func (rt *impl) routingCore(r route.Interface) {
	r.Get("/version", controllers.VersionController.Version)
	r.Get("/ping", controllers.PingController.Ping)
	r.Get("/panic", controllers.PanicController.Panic)
	r.Put("/settings/time", controllers.SettingsControllerTime.Time)
}

// Сущность myitem
func (rt *impl) routingMyItem(r route.Interface) {
	r.Options("/", controllers.MyitemController.Status)
	r.Post("/", controllers.MyitemController.Create)
	r.Get("/:id", controllers.MyitemController.Load)
	r.Delete("/:id", controllers.MyitemController.Delete)
}
