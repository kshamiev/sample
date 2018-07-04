package launcher // import "application/modules/launcher"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"os"
	"os/exec"
)

// Interface is an interface of package
type Interface interface {
	// Launch Выполнение внешней команды или приложения
	// При запуске приложения устанавливаются переменные окружения
	// Если stdout, stderr не nil, то построчно пишется STDOUT и STDERR запущенного приложения
	Launch(cmd []string, environment []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error

	// ForkLaunch Выполнение внешней команды или приложения в защищенном виде через fork
	ForkLaunch(cmd []string, environment []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error

	// Interrupt Прерывание выполнения программы
	Interrupt()

	// Kill Завершение выполнения программы
	Kill()

	// Wait Ожидание завершения запущенной программы
	// Возвращает код ошибки возвращенный приложением
	Wait() (int, error)
}

// impl is an implementation of package
type impl struct {
	stdout   io.ReadCloser
	stdoutWr *os.File
	stderr   io.ReadCloser
	stderrWr *os.File
	cmd      *exec.Cmd
	pid      int
}
