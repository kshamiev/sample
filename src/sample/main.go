package main

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"

	"application/configuration"
	"application/workflow"

	"gopkg.in/alecthomas/kingpin.v2"

	// Application components registration. Each component determines its own dependencies.
	// Components are initialized and run in the order specified by dependencies.
	// If the dependencies are not specified or equivalent, the initialization and run of the components is
	// carried out in the order of they registration
	_ "application/components/bootstrap"     // Выполняется после основных зависимостей, вспомогательный компонент
	_ "application/components/configuration" // Конфигурация и выполнение команды version
	_ "application/components/daemon"        // Воркеры, процессы, службы, выполнение команды daemon
	_ "application/components/environment"   // Окружение, процессоры, рандомизатор и т.п.
	_ "application/components/interrupt"     // Перехват сигналов прерывания приложения
	_ "application/components/logging"       // Настройка системы логирования
	_ "application/components/migrations"    // Применение миграций базы данных
	_ "application/components/pidfile"       // PID файла процесса
)

var (
	// Application build version
	// Sets application build version from args:
	//   -X main.build=$(date -u +%Y%m%d.%H%M%S.%Z)
	// Default value: "dev"
	build = `dev`
)

func main() {
	// Ending an application with an error code
	if exitCode := int(Main()); exitCode != 0 {
		os.Exit(exitCode)
	}
}

// Main The application entry point
func Main() (exitCode uint8) {
	var err error
	var wfw workflow.Interface

	// Waiting to complete shutdown of logging system
	defer log.Done()
	wfw = workflow.Get()
	// Running initialization of all registered components
	if exitCode, err = wfw.Init(version, build); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	// Stopping all registered components before return from Main
	defer func() {
		if exitCode != workflow.ErrNone || err != nil {
			return
		}
		exitCode, err = wfw.Stop()
		if exitCode != workflow.ErrNone || err != nil {
			log.Fatal(err)
		}
	}()
	// Running the command in the registered components
	exitCode, err = wfw.Start(
		configuration.Get().Command(),
		kingpin.Usage,
	)
	if exitCode != workflow.ErrNone || err != nil {
		log.Error(err)
		return
	}

	return
}
