package cleaner // import "application/workers/cleaner"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"context"
	"time"

	"application/configuration"
	modelsFilestore "application/models/filestore"

	"gopkg.in/webnice/job.v1/job"
	"gopkg.in/webnice/job.v1/types"
)

// Injecting a process object to goroutine management tools
func init() {
	job.Get().RegisterWorker(New())
}

// New creates a new object and return interface
func New() Interface {
	var clr = new(impl)
	clr.Ctx, clr.CtxCfn = context.WithCancel(context.Background())
	clr.Cfg = configuration.Get()
	return clr
}

// Info Функция конфигурации процесса,
// - в функцию передаётся уникальный идентификатор присвоенный процессу
// - функция должна вернуть конфигурацию или nil, если возвращается nil, то применяется конфигурация по умолчанию
func (clr *impl) Info(id string) (ret *types.Configuration) {
	clr.ID, ret = id, &types.Configuration{
		Autostart:      true,
		Fatality:       true,
		Restart:        true,
		RestartTimeout: time.Second,
	}
	return
}

// Prepare Функция выполнения действий подготавливающих воркер к работе
// Завершение с ошибкой означает, что процесс не удалось подготовить к запуску
func (clr *impl) Prepare() (err error) { return }

// Cancel Функция прерывания работы
func (clr *impl) Cancel() (err error) { clr.CtxCfn(); return }

// Worker Функция-реализация процесса, данная функция будет запущена в горутине
// до тех пор пока функция не завершился воркер считается работающим
func (clr *impl) Worker() (err error) {
	var errors []error
	var ufm modelsFilestore.Interface
	var secondsTicker *time.Ticker
	var minutesTicker *time.Ticker
	var exit bool

	// Создание генераторов разной периодичности
	secondsTicker = time.NewTicker(time.Second)
	defer secondsTicker.Stop()
	minutesTicker = time.NewTicker(time.Minute)
	defer minutesTicker.Stop()

	// Интерфейсы моделей
	ufm = modelsFilestore.New()
	// Цикл прерывается только через Cancel() или panic :)
	for {
		errors = errors[:0]
		select {
		// Выход по команде прерывания Cancel()
		case <-clr.Ctx.Done():
			err = clr.Ctx.Err()
			exit = true

		// Секундный таймер
		case <-secondsTicker.C:

			//

		// Минутный таймер
		case <-minutesTicker.C:
			// Очистка filestore от старых данных
			if err = ufm.CleanOldData(); err != nil {
				err, errors = nil, append(errors, err)
			}
		}
		if exit {
			break
		}

		if len(errors) > 0 {
			for _, err = range errors {
				log.Errorf("Error: %s", err)
			}
		}
		err = nil
	}

	return
}
