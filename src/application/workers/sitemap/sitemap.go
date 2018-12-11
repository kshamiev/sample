package sitemap // import "application/workers/sitemap"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"application/controllers"
	modelsSitemap "application/models/sitemap"

	"gopkg.in/webnice/job.v1/job"
	"gopkg.in/webnice/job.v1/types"
)

// Injecting a process object to goroutine management tools
func init() {
	job.Get().RegisterWorker(New())
}

// New creates a new object and return interface
func New() Interface {
	var wsm = new(impl)
	return wsm
}

// Info Функция конфигурации процесса,
// - в функцию передаётся уникальный идентификатор присвоенный процессу
// - функция должна вернуть конфигурацию или nil, если возвращается nil, то применяется конфигурация по умолчанию
func (wsm *impl) Info(id string) (ret *types.Configuration) {
	wsm.ID, ret = id, &types.Configuration{
		Autostart:      false,
		Restart:        false,
		Fatality:       true,
		RestartTimeout: time.Second * 5,
	}

	return
}

// Prepare Функция выполнения действий подготавливающих воркер к работе
// Завершение с ошибкой означает, что процесс не удалось подготовить к запуску
func (wsm *impl) Prepare() (err error) {
	return
}

// Cancel Функция прерывания работы
func (wsm *impl) Cancel() (err error) {
	return
}

// Worker Функция-реализация процесса, данная функция будет запущена в горутине
// до тех пор пока функция не завершился воркер считается работающим
func (wsm *impl) Worker() (err error) {
	var items []*modelsSitemap.Location
	var item *modelsSitemap.Location
	var buf []byte
	var max int
	var i int
	var del []string

	//max = int(1000000)
	max = int(2000000)
	log.Noticef(" - create sitemap elements")
	items = make([]*modelsSitemap.Location, 0, max)
	for i = 0; i < max; i++ {
		buf, _ = time.Now().MarshalBinary() // nolint: gosec
		item = &modelsSitemap.Location{
			URN: fmt.Sprintf(
				"/%09d/path/to/resource/%s/%s",
				i,
				hex.EncodeToString(buf),
				strconv.FormatInt(time.Now().UnixNano(), 10),
			),
			ModTime:  time.Now(),
			Change:   modelsSitemap.CfAlways,
			Priority: -1,
		}
		items = append(items, item)
	}
	log.Noticef(" - sitemap elements begin add to key/value storage")
	if err = controllers.ResourceController.
		Sitemap().
		Add(items); err != nil {
		log.Errorf(" - add urn error: %s", err)
	}
	log.Noticef(" - create sitemap elements done")

	// Удаление
	del = []string{}
	log.Noticef(" - delete %q", del)
	for _, delItem := range del {
		if err = controllers.ResourceController.
			Sitemap().
			Del(delItem); err != nil {
			log.Errorf(" - delete urn %q error: %s", delItem, err)
		}
	}
	err = nil
	log.Noticef(" - delete done")
	log.Noticef(" - count: %d", controllers.ResourceController.Sitemap().Count())

	return
}
