package environment // impoert "application/componens/environment"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"application/configuration"
	"application/workflow"
)

func init() { workflow.Register(New()) }

// New Create object and return interface
func New() Interface {
	return new(impl)
}

// After Возвращает массив зависимостей - имена пакетов компонентов, которые должны быть запущены до этого компонента
func (cpn *impl) After() []string {
	return []string{}
}

// Init Функция инициализации компонента
func (cpn *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
	return
}

// Start Выполнение компонента
func (cpn *impl) Start(cmd string) (done bool, exitCode uint8, err error) {
	if configuration.Get().Debug() {
		cpn.Ctx, cpn.CtxCfn = context.WithCancel(context.Background())
		go cpn.gcSheduler()
	}

	return
}

// Stop Функция завершения работы компонента
func (cpn *impl) Stop() (exitCode uint8, err error) {
	if cpn.CtxCfn != nil {
		cpn.CtxCfn()
	}

	return
}

// Горутина постоянного запуска GC для очистик памяти. Выполняется только в режиме дебага приложения
// Так как GC по умолчанию оптимизирован и запускается при накоплении 100% переменных от общей кучи
// в режиме дебага запускаем GC раз в секунду для отладки финализаторов и проверки утечки памяти
func (cpn *impl) gcSheduler() {
	var tic = time.NewTicker(time.Second)
	defer tic.Stop()
	for {
		select {
		case <-tic.C:
			runtime.GC()
		case <-cpn.Ctx.Done():
			return
		}
	}
}
