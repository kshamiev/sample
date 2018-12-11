package controllers // import "application/controllers"

import (
	"application/controllers/apiV1/core/panic"
	"application/controllers/apiV1/core/ping"
	"application/controllers/apiV1/core/settings"
	"application/controllers/apiV1/core/uploadfile"
	"application/controllers/apiV1/core/version"
	"application/controllers/apiV1/myitem"
	"application/controllers/internalServerError"
	"application/controllers/resource"
)

var (
	// CORE

	// InternalServerErrorController is an interface for controller implementation
	InternalServerErrorController = internalServerError.New()

	// API v.1

	// VersionController is an interface for controller implementation
	VersionController = version.New()
	// PingController is an interface for controller implementation
	PingController = ping.New()
	// PanicController is an interface for controller implementation
	PanicController = panic.New()
	// SettingsControllerTime controller interface
	SettingsControllerTime = settings.New()
	// UploadFileController controller interface
	UploadFileController = uploadfile.New()

	// MyitemController controller interface
	MyitemController = myitem.New()

	// Контроллер статических и шаблонных ресурсов

	// ResourceController controller interface
	ResourceController = resource.Get()
)
