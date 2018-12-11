package sitemap // import "application/controllers/resource/sitemap"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/controllers/resource/pool"
	"application/models/filecache"
	modelsSitemap "application/models/sitemap"

	"gopkg.in/webnice/web.v1/route"
)

// New creates a new object and return interface
func New() Interface {
	var err error
	var smi = &impl{
		Mfc:  filecache.Get(),
		Pool: pool.New(),
	}
	smi.Smm, err = modelsSitemap.New(configuration.Get().CachePath())
	if err != nil {
		log.Criticalf("create model of sitemap error: %s", err)
	}

	return smi
}

// Close opened connections
func (smi *impl) Close() error { return smi.Smm.Close() }

// Debug Set debug mode
func (smi *impl) Debug(d bool) Interface {
	smi.debug = d
	smi.Smm.Debug(smi.debug)
	smi.Pool.Debug(smi.debug)
	return smi
}

// Errors Ошибки известного состояни, которые могут вернуть функции пакета
func (smi *impl) Errors() *Error { return Errors() }

// DocumentRoot Устанавливает путь к корню веб сервера
func (smi *impl) DocumentRoot(path string) Interface { smi.rootPath = path; return smi }

// ServerURL Устанавливает основной адрес веб сервера
func (smi *impl) ServerURL(url string) (Interface, error) {
	if !rexURLCheck.MatchString(url) {
		return smi, smi.Errors().ErrInvalidURL()
	}
	smi.serverURL = rexSlashLast.ReplaceAllString(url, ``)
	return smi, nil
}

// SetRouting Установка роутинга к статическим файлам
func (smi *impl) SetRouting(rou route.Interface) Interface {
	rou.Get(urnPattern, smi.SitemapHandlerFunc)
	return smi
}

// Add Добавление URI ресурса в sitemap.xml
func (smi *impl) Add(items []*modelsSitemap.Location) (err error) {
	for i := range items {
		if len(items[i].URN) > 2048 {
			err = smi.Errors().ErrURIisToLong()
			return
		}
		if items[i].Priority == -1.0 {
			items[i].Priority = modelsSitemap.DefaultPriority
		}
		if items[i].Priority < 0.0 || items[i].Priority > 1.0 {
			err = smi.Errors().ErrPriorityNotCorrect()
			return
		}
		switch items[i].Change {
		case modelsSitemap.CfAlways, modelsSitemap.CfHourly, modelsSitemap.CfDaily, modelsSitemap.CfWeekly,
			modelsSitemap.CfMonthly, modelsSitemap.CfYearly, modelsSitemap.CfNever:
		default:
			err = smi.Errors().ErrChangeFrequencyNotCorrect()
			return
		}
	}
	_, err = smi.Smm.Add(items)

	return
}

// Del Удаление URI ресурса из sitemap.xml
func (smi *impl) Del(urn string) (err error) {
	_, err = smi.Smm.Del(urn)
	return
}

// Count Количество занисей в sitemap
func (smi *impl) Count() uint64 { return smi.Smm.Count() }
