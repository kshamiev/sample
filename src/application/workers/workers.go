package workers // import "application/workers"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/workers/web"

	_ "application/workers/cleaner" // Очистка системы
	_ "application/workers/sitemap" // Генератор sitemap ссылок

	"gopkg.in/webnice/job.v1/job"
)

// Init Инициализация специальных воркеров
func Init() (jbo job.Interface, err error) {
	var conf configuration.Interface

	jbo, conf = job.Get(), configuration.Get()
	if conf.Configuration() == nil {
		return
	}
	// Ручная регистрация веб серверов
	for i := range conf.Configuration().WEBServers {
		web.Init(&conf.Configuration().WEBServers[i])
	}
	// В дебаг режиме некоторые воркеры не запускаем, удаляем их регистрацию
	//if conf.Debug() {
	//	_ = jbo.Unregister("application/workers/stat/impl")
	//}

	return
}
