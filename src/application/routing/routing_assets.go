package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/controllers"
	"application/middleware/minify"

	"gopkg.in/webnice/web.v1/route"
)

// Assets Статические файлы и шаблоны страниц
func (rt *impl) Assets() Interface {
	// Конфигурация контроллера ресурсов (singleton) через пакет "application/controllers"
	controllers.ResourceController.
		Debug(rt.debug).
		WebServerConfiguration(rt.Wsc)
	// ------------------------------
	// Статические файлы /assets
	// GET /assets/*
	rt.Rou.Group(func(sr route.Interface) {
		// В продакшн режиме минифицируются на лету
		if !rt.debug {
			sr.Use(minify.Minify)
		}
		controllers.ResourceController.
			Assets().
			Debug(rt.debug).
			SetRouting(sr)
	})
	// ------------------------------
	// Файл информирования роботов /robots.txt
	// GET /robots.txt
	controllers.ResourceController.
		Robots().
		Debug(rt.debug).
		SetRouting(rt.Rou)
	// ------------------------------
	// Генератор иконок для всех платформ и браузеров
	// GET /favicon*
	// GET /apple-touch*.png
	// GET /android-chrome-*.png
	// GET /mstile*.png
	// GET /icons_preview_page
	rt.Rou.Group(func(sr route.Interface) {
		// В продакшн режиме минифицируются на лету
		if !rt.debug {
			sr.Use(minify.Minify)
		}
		controllers.ResourceController.
			Icons().
			Debug(rt.debug).
			SetRouting(sr)
	})
	// ------------------------------
	// Генератор sitemap.xml и sitemap-index.xml
	// GET /sitemap.xml
	// GET /sitemap.xml/:number
	// GET /sitemap-index.xml
	// GET /sitemap-index.xml/:number
	rt.Rou.Group(func(sr route.Interface) {
		// В продакшн режиме минифицируются на лету
		if !rt.debug {
			sr.Use(minify.Minify)
		}
		controllers.ResourceController.
			Sitemap().
			Debug(rt.debug).
			SetRouting(sr)
	})
	// ------------------------------
	// Генератор статических HTML страниц на основе шаблонов
	// Обработка /route.txt для обеспечения работы SPA/PWA приложения
	// Принцип работы:
	// - Если в DocumentRoot есть не пустой файл index.html,
	//   то при запросах /, /index.htm и index.html возвращается его содержимое
	// - Если в DocumentRoot нет файла index.html или он пустой,
	//   то выполняется контроллер Pages, обрабатывая шаблоны из templates/pages
	// - Если в DocumentRoot нет файла index.html или он пустой и контроллер Pages возвращает ошибку,
	//   то роутинг не настраивается, а при обращении клиента веб сервер возвращает 404 ошибку
	// - Если есть контент для страницы /, то загружается и обрабатывается файл настроек
	//   дополнительного роутинга: /route.txt
	// ------------------------------
	// Обеспечение выгрузки из DocumentRoot файлов:
	// - /:filename - Любой файл находящийся в корне, за исключением файлов - шаблонов
	// ------------------------------
	rt.Rou.Group(func(sr route.Interface) {
		// В продакшн режиме минифицируются на лету
		if !rt.debug {
			sr.Use(minify.Minify)
		}
		controllers.ResourceController.
			Pages().
			Debug(rt.debug).
			SetRouting(sr)
	})

	return rt
}
