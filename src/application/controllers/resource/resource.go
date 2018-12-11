package resource // import "application/controllers/resource"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	resourceAssets "application/controllers/resource/assets"
	resourceIcons "application/controllers/resource/icons"
	resourcePages "application/controllers/resource/pages"
	resourceRobots "application/controllers/resource/robots"
	resourceSitemap "application/controllers/resource/sitemap"
	webserverTypes "application/webserver/types"
)

func init() {
	singleton = new(impl)
}

// Get an singleton object of package interface
func Get() Interface { return singleton }

// Debug Set debug mode
func (rsc *impl) Debug(d bool) Interface { rsc.debug = d; return rsc }

// WebServerConfiguration Конфигурация веб сервера
func (rsc *impl) WebServerConfiguration(wsc *webserverTypes.Configuration) Interface {
	rsc.Wsc = wsc
	return rsc
}

// Assets Возвращает интерфейс статических файлов
func (rsc *impl) Assets() resourceAssets.Interface {
	if rsc.Asi == nil {
		rsc.Asi = resourceAssets.New().
			DocumentRoot(rsc.Wsc.DocumentRoot)
	}
	rsc.Asi.Debug(rsc.debug)
	return rsc.Asi
}

// Icons Возвращает интерфейс генератора icons файлов
func (rsc *impl) Icons() resourceIcons.Interface {
	if rsc.Ici == nil {
		rsc.Ici = resourceIcons.New().
			DocumentRoot(rsc.Wsc.DocumentRoot)
	}
	rsc.Ici.Debug(rsc.debug)
	return rsc.Ici
}

// Pages Возвращает интерфейс генератора статических html страниц на основе шаблонов
func (rsc *impl) Pages() resourcePages.Interface {
	if rsc.Pgi == nil {
		rsc.Pgi = resourcePages.New().
			ServerURL(rsc.Wsc.Server.Address).
			DocumentRoot(rsc.Wsc.DocumentRoot).
			TemplatePages(rsc.Wsc.Pages)
	}
	rsc.Pgi.Debug(rsc.debug)
	return rsc.Pgi
}

// Robots Возвращает интерфейс robots.txt
func (rsc *impl) Robots() resourceRobots.Interface {
	if rsc.Rbi == nil {
		rsc.Rbi = resourceRobots.New().
			DocumentRoot(rsc.Wsc.DocumentRoot).
			ServerURL(rsc.Wsc.Server.Address).
			Sitemap(rsc.Sitemap())
	}
	rsc.Rbi.Debug(rsc.debug)
	return rsc.Rbi
}

// Sitemap Возвращает интерфейс генератора sitemap и sitemap-index файлов
func (rsc *impl) Sitemap() resourceSitemap.Interface {
	var err error

	if rsc.Smi == nil {
		rsc.Smi, err = resourceSitemap.New().
			DocumentRoot(rsc.Wsc.DocumentRoot).
			ServerURL(rsc.Wsc.Server.Address)
		if err != nil {
			log.Warningf("controllers/resource/sitemap error: %s", err)
		}
	}
	rsc.Smi.Debug(rsc.debug)
	return rsc.Smi
}
