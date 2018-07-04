package environment

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"math/rand"
	"runtime"
	"time"

	"application/configuration"
	"application/workflow"
)

// impl is an implementation of package
type impl struct {
}

func init() {
	workflow.Register(new(impl))
}

// Init Plug-in initialization function
func (cgn *impl) Init(appVersion string, appBuild string) (exitCode int, err error) {
	exitCode = workflow.ErrNone
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
	return
}

// Command Processing the command
func (cgn *impl) Command(cmd string) (exitCode int, err error, done bool) {

	exitCode = workflow.ErrNone
	if cmd == `daemon` {
		log.Info(`Application initialized successfully`)
	}
	if configuration.Get().Debug() {
		// Запуск GC только в режиме дебага
		// Так как GC по умолчанию оптимизирован и запускается при накоплении 100% переменных от общей кучи
		// в режиме дебага запускаем GC раз в секунду для отладки финализаторов и проверки утечки памяти
		go func() {
			for {
				runtime.GC()
				time.Sleep(time.Second)
			}
		}()
	}

	return
}

// Done The function is called when the application terminates
func (cgn *impl) Shutdown() {
}
