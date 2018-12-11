package resource // import "application/controllers/resource"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	resourceAssets "application/controllers/resource/assets"
	resourceIcons "application/controllers/resource/icons"
	resourcePages "application/controllers/resource/pages"
	resourceRobots "application/controllers/resource/robots"
	resourceSitemap "application/controllers/resource/sitemap"
	webserverTypes "application/webserver/types"
)

// singleton object
var singleton *impl

// Interface is an interface of package
type Interface interface {
	// Debug Set debug mode
	Debug(d bool) Interface

	// WebServerConfiguration Конфигурация веб сервера
	WebServerConfiguration(wsc *webserverTypes.Configuration) Interface

	// Assets Возвращает интерфейс статических файлов
	Assets() resourceAssets.Interface

	// Icons Возвращает интерфейс генератора icons файлов
	Icons() resourceIcons.Interface

	// Pages Возвращает интерфейс генератора статических html страниц на основе шаблонов
	Pages() resourcePages.Interface

	// Robots Возвращает интерфейс robots.txt
	Robots() resourceRobots.Interface

	// Sitemap Возвращает интерфейс генератора sitemap и sitemap-index файлов
	Sitemap() resourceSitemap.Interface
}

// impl is an implementation of package
type impl struct {
	debug bool                          // =true - debug mode is on
	Wsc   *webserverTypes.Configuration // Конфигурация веб сервера
	Asi   resourceAssets.Interface      // Интерфейс статических файлов
	Ici   resourceIcons.Interface       // Интерфейс генератора favicon, apple-touch-icon, android-chrome, mstile
	Pgi   resourcePages.Interface       // Интерфейс генератора статических html страниц на основе шаблонов
	Rbi   resourceRobots.Interface      // Интерфейс robots.txt
	Smi   resourceSitemap.Interface     // Интерфейс генератора sitemap и sitemap-index файлов
}
