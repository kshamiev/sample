package interrupt // import "application/modules/interrupt"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

// Interface is an interface of package
type Interface interface {
	Start() Interface
	Stop() Interface
}

// impl is an implementation of package
type impl struct {
	DoExitUp   chan bool      // Канал сигнала завершении горутины
	DoExitDone chan bool      // Канал сигнала о горутина завершена
	Catch      CatchFn        // Функция вызывается при перехвате сигнала os
	Signal     chan os.Signal // Канал поступления сигналов от os
	IsRun      atomic.Value   // =true-горутина работает
}

// CatchFn Описание функции вызываемой при поступлении прерывания от OS
type CatchFn func(os.Signal)

// New Create new object
func New(fn CatchFn) Interface {
	var itp = new(impl)
	itp.Catch = fn
	itp.Signal = make(chan os.Signal, 100)
	itp.DoExitUp = make(chan bool, 1)
	itp.DoExitDone = make(chan bool, 1)
	itp.IsRun.Store(false)
	return itp
}

// Start Запуск перехвата сигналов прерывания
func (itp *impl) Start() Interface {
	defer func() {
		if re := recover(); re != nil {
			log.Criticalf("Interrupt panic recover: %v", re)
		}
	}()
	// Если горутина перехвата работает то второй раз старт не вызывается
	if itp.IsRun.Load().(bool) {
		return itp
	}
	itp.IsRun.Store(true)
	signal.Notify(itp.Signal,
		syscall.SIGABRT,
		syscall.SIGALRM,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	go itp.Catcher()
	return itp
}

// Stop Остановка перехвата сигналов прерывания
func (itp *impl) Stop() Interface {
	// Если горутина перехвата не работает то останавливать нечего
	if !itp.IsRun.Load().(bool) {
		return itp
	}
	signal.Stop(itp.Signal)
	// Выход после того как убедимся что горутина остановилась
	itp.DoExitUp <- true
	<-itp.DoExitDone
	return itp
}

// Catcher Получаем сигналы, вызываем вункцию
func (itp *impl) Catcher() {
	var exit bool
	var sig os.Signal
	for {
		if exit {
			break
		}
		select {
		case <-itp.DoExitUp:
			exit = true
		case sig = <-itp.Signal:
			if itp.Catch != nil {
				itp.Catch(sig)
			}
		}
	}
	itp.DoExitDone <- true
}
