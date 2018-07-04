package controllers // import "application/controllers"

import (
	"application/controllers/apiV1/core/panic"
	"application/controllers/apiV1/core/ping"
	"application/controllers/apiV1/core/settings"
	"application/controllers/apiV1/core/version"
	"application/controllers/apiV1/myitem"
	"application/controllers/internal_server_error"
	"application/controllers/pages/appleicon"
	"application/controllers/pages/assets"
	"application/controllers/pages/favicon"
	"application/controllers/pages/index"
	"application/controllers/pages/robots"
)

var (
	// CORE

	// InternalServerErrorController is an interface for controller implementation
	InternalServerErrorController = internal_server_error.New()

	// API v.1

	// VersionController is an interface for controller implementation
	VersionController = version.New()
	// PingController is an interface for controller implementation
	PingController = ping.New()
	// PanicController is an interface for controller implementation
	PanicController = panic.New()
	// SettingsControllerTime controller interface
	SettingsControllerTime = settings.New()

	// MyitemController controller interface
	MyitemController = myitem.New()

	// STATIC

	// SpaController Контроллер выдающий шаблоны SPA
	IndexController = index.New()
	// AssetsController Контроллер статических файлов
	AssetsController = assets.New()
	// RobotsController robots.txt
	RobotsController = robots.New()
	// FaviconController Контроллер создания на лету favicon.ico и всех возможных вариантов размеров и форматов favicon
	FaviconController = favicon.New()
	// AppleiconController Контроллер создания на лету apple-touch-icon и всех возможных вариантов размеров и форматов
	AppleiconController = appleicon.New()
)
