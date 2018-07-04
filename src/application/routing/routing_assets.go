package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/route"
import (
	"application/controllers"
	"application/middleware/minify"
)

// Assets Статические файлы и шаблоны страниц
func (rt *impl) Assets() Interface {
	// Static files
	rt.rou.Subroute("/assets", rt.assetsStatic)
	// GET /robots.txt
	rt.rou.Get("/robots.txt", controllers.RobotsController.RobotsTxt)
	// GET /sitemap.xml
	// GET /sitemap-index.xml

	// Root index page
	// Для всех URN приложения на ангуляре необходимо отвечать одним и тем же index.html
	rt.rou.Subroute("/", rt.assetsRoot)
	rt.rou.Subroute("/about", rt.assetsRoot)
	rt.rou.Subroute("/contacts", rt.assetsRoot)
	rt.rou.Subroute("/path/to/other/section", rt.assetsRoot)

	// GET /favicon*
	//  <link rel="icon" type="image/vnd.microsoft.icon" href="/favicon.ico" />
	//  <link rel="icon" type="image/png" href="/favicon.png" />
	//  <link rel="shortcut icon" href="/favicon.ico">
	//  <link rel="manifest" href="/assets/favicon-android-chrome-manifest.json">
	rt.rou.Get("/favicon.ico", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-16x16.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-32x32.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-36x36.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-48x48.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-70x70.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-72x72.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-96x96.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-128x128.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-144x144.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-150x150.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-192x192.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-310x150.png", controllers.FaviconController.FavIcon)
	rt.rou.Get("/favicon-310x310.png", controllers.FaviconController.FavIcon)

	// GET /apple-touch*.png
	//  <link rel="apple-touch-icon" href="/apple-touch-icon.png">
	//  <link rel="apple-touch-startup-image" href="/apple-touch-startup-image.png">
	//  <link rel="apple-touch-icon" sizes="57x57" href="/apple-touch-icon-57x57.png">
	//  <link rel="apple-touch-icon" sizes="60x60" href="/apple-touch-icon-60x60.png">
	//  <link rel="apple-touch-icon" sizes="72x72" href="/apple-touch-icon-72x72.png">
	//  <link rel="apple-touch-icon" sizes="76x76" href="/apple-touch-icon-76x76.png">
	//  <link rel="apple-touch-icon" sizes="114x114" href="/apple-touch-icon-114x114.png">
	//  <link rel="apple-touch-icon" sizes="120x120" href="/apple-touch-icon-120x120.png">
	//  <link rel="apple-touch-icon" sizes="144x144" href="/apple-touch-icon-144x144.png">
	//  <link rel="apple-touch-icon" sizes="152x152" href="/apple-touch-icon-152x152.png">
	//  <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon-180x180.png">
	rt.rou.Get("/apple-touch-icon-57x57.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-60x60.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-72x72.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-76x76.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-114x114.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-120x120.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-144x144.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-152x152.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-180x180.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-320x460.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-640x920.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-640x1096.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-748x1024.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-768x1004.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-1496x2048.png", controllers.AppleiconController.AppleIcon)
	rt.rou.Get("/apple-touch-icon-1536x2008.png", controllers.AppleiconController.AppleIcon)

	return rt
}

// Все URN на которые должен выдаваться index.html сделанный как Single Page Application (SPA)
func (rt *impl) assetsRoot(r route.Interface) {
	// В дебаг режиме минификация отключается
	if !rt.cfg.Debug() {
		r.Use(minify.Minify)
	}
	r.Get("/:filename", controllers.IndexController.Index)
}

// Все статические файлы не требующие никаких преобразований (чистая статиска)
func (rt *impl) assetsStatic(r route.Interface) {
	if !rt.cfg.Debug() {
		r.Use(minify.Minify)
	}
	r.Get("/*", controllers.AssetsController.Assets)
}
