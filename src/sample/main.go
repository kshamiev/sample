package main

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"

	"application/configuration"
	"application/workflow"

	// Plugins in order of initialization and execution
	// Окружение, процессоры, рандомизатор и т.п.
	_ "sample/environment"
	// Конфигурация
	_ "sample/configuration"
	// Настройка системы логирования
	_ "application/logging"
	// Применение миграций базы данных
	_ "sample/migrations"
	// Перехват сигналов прерывания процесса
	_ "sample/interrupt"
	// PID файла процесса
	_ "sample/pidfile"
	// Воркеры, процессы, службы, задачи (параллельные управляемые вычисления)
	_ "sample/daemon"
)

var (
	// Application build version
	// Sets build version from args:
	//   -X main.build=$(date -u +%Y%m%d.%H%M%S.%Z)
	// Default value: "dev"
	build = `dev`
)

func main() {
	// Ending an application with an error code
	os.Exit(Main())
}

// Main The application entry point
func Main() (exitCode int) {
	var err error
	var wfw = workflow.Get()

	// Initialize all registered plugins
	if exitCode, err = wfw.Initialize(version, build); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		_, _ = fmt.Fprintln(os.Stderr, "")
		return
	}
	defer wfw.Shutdown()
	// Running the command in the registered plugins
	if exitCode, err = wfw.Command(configuration.Get().Command()); err != nil {
		log.Error(err.Error())
		return
	}

	return
}
