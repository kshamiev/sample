package cleaner // import "application/workers/cleaner"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"time"

	modelsMyitem "application/models/myitem"

	"gopkg.in/webnice/job.v1/job"
	"gopkg.in/webnice/job.v1/types"
)

// impl is an implementation of package
type impl struct {
	ID string
}

// Injecting a process object to goroutine management tools
func init() {
	job.Get().RegisterWorker(newWorker())
}

// Конструктор объекта
func newWorker() types.WorkerInterface {
	return &impl{}
}

// Info Функция конфигурации процесса,
// - в функцию передаётся уникальный идентификатор присвоенный процессу
// - функция должна вернуть конфигурацию или nil, если возвращается nil, то применяется конфигурация по умолчанию
func (wrk *impl) Info(id string) (ret *types.Configuration) {
	wrk.ID, ret = id, &types.Configuration{
		Autostart:      true,
		Fatality:       true,
		Restart:        true,
		RestartTimeout: time.Minute / 2,
		CancelTimeout:  time.Minute * 3 / 4,
	}
	return
}

// Prepare Функция выполнения действий подготавливающих воркер к работе
// Завершение с ошибкой означает, что процесс не удалось подготовить к запуску
func (wrk *impl) Prepare() (err error) { return }

// Cancel Функция прерывания работы
func (wrk *impl) Cancel() (err error) { return }

// Worker Функция-реализация процесса, данная функция будет запущена в горутине
// до тех пор пока функция не завершился воркер считается работающим
func (wrk *impl) Worker() (err error) {
	var mim modelsMyitem.Interface

	// Очистка БД от удалённых записей
	mim = modelsMyitem.New()
	if err = mim.Clean(); err != nil {
		log.Warningf("Model Cleanner error: %s", err)
		// Если выйти с err != nil  - то сработает Fatality=true
		err = nil
	}

	// Очистка от чего-то еще

	// -- // --

	return
}
