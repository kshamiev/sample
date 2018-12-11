package launcher // import "application/modules/launcher"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"text/template"
)

// New Create new package object and return interface
func New() Interface {
	var lau = new(impl)
	return lau
}

// Launch Выполнение внешней команды или приложения
// При запуске приложения устанавливаются переменные окружения
// Если stdout, stderr не nil, то построчно пишется STDOUT и STDERR запущенного приложения
func (lau *impl) Launch(cmd []string, environment []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) (err error) {
	var command string
	var env []string
	var c = int(1)

	if len(cmd) == 0 {
		err = fmt.Errorf("Do not specify command and arguments that must run")
		return
	}

	// Подготовка к выполнению команды
	if command, err = exec.LookPath(cmd[0]); err != nil {
		return
	}
	lau.cmd = exec.Command(command, cmd[c:]...) // nolint: errcheck, gosec

	// Environment
	env = os.Environ()
	lau.cmd.Env = append(env, environment...)

	// Pipes
	if lau.stdout, err = lau.cmd.StdoutPipe(); err != nil {
		return
	}
	if lau.stderr, err = lau.cmd.StderrPipe(); err != nil {
		return
	}

	// Запуск приложения
	if err = lau.cmd.Start(); err != nil {
		return
	}

	// Сбор вывода приложения
	if stdout != nil {
		go lau.rdrw(lau.stdout, stdout)
	}
	if stderr != nil {
		go lau.rdrw(lau.stderr, stderr)
	}

	return
}

// ForkLaunch Выполнение внешней команды или приложения в защищенном виде через fork
func (lau *impl) ForkLaunch(cmd []string, environment []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) (err error) {
	var command, workDir string
	var env []string

	if len(cmd) == 0 {
		err = fmt.Errorf("Do not specify command and arguments that must run")
		return
	}
	// Подготовка к выполнению команды
	if command, err = exec.LookPath(cmd[0]); err != nil {
		return
	}
	// Environment
	env = append(os.Environ(), environment...)
	if workDir, err = os.Getwd(); err != nil {
		return
	}

	// STDOUT pipe
	if lau.stdout, lau.stdoutWr, err = os.Pipe(); err != nil {
		return
	}
	// STDERR pipe
	if lau.stderr, lau.stderrWr, err = os.Pipe(); err != nil {
		return
	}
	lau.pid, err = syscall.ForkExec(command, cmd, &syscall.ProcAttr{
		Dir: workDir,
		Env: env,
		Files: []uintptr{
			os.Stdin.Fd(),
			lau.stdoutWr.Fd(),
			lau.stderrWr.Fd(),
		},
	})
	if err != nil {
		return
	}
	// Сбор вывода приложения
	if stdout != nil {
		go lau.rdrw(lau.stdout, stdout)
	}
	if stderr != nil {
		go lau.rdrw(lau.stderr, stderr)
	}

	return
}

// Копирует по-строчно из одного потока в другой
func (lau *impl) rdrw(ri io.Reader, wi io.Writer) {
	var err error
	var rdrBf []byte
	var rdr = bufio.NewReader(ri)
	var wrt = bufio.NewWriter(wi)
	var count int64

	for err == nil {
		count++
		rdrBf, _, err = rdr.ReadLine()
		if len(rdrBf) > 0 {
			if count > 1 {
				wrt.WriteString("\n") // nolint: errcheck, gosec
			}
			wrt.Write(rdrBf) // nolint: errcheck, gosec
			wrt.Flush()      // nolint: errcheck, gosec
		}
	}
}

// Kill Завершение выполнения программы
func (lau *impl) Kill() { lau.signal(os.Kill) }

// Interrupt Прерывание выполнения программы
func (lau *impl) Interrupt() { lau.signal(os.Interrupt) }

// Отправка OS сигнала прерывания
func (lau *impl) signal(sig os.Signal) {
	var err error
	var prc *os.Process

	if lau.cmd != nil {
		if lau.cmd.Process == nil {
			return
		}
		if err = lau.cmd.Process.Signal(sig); err != nil {
			log.Errorf("Error send signal %q to PID %d: %s", sig, lau.cmd.Process.Pid, err)
		}
	}
	if lau.pid > 0 {
		if prc, err = os.FindProcess(lau.pid); err != nil {
			return
		}
		if err = prc.Signal(sig); err != nil {
			log.Errorf("Error send signal %q to PID %d: %s", sig, lau.pid, err)
			return
		}
	}
}

// Wait Ожидание завершения запущенной программы
// Возвращает код ошибки возвращенный приложением
func (lau *impl) Wait() (ret int, err error) {
	var prc *os.Process
	var pst *os.ProcessState

	ret = -1
	if lau.cmd != nil {
		defer lau.stdout.Close() // nolint: errcheck, gosec
		defer lau.stderr.Close() // nolint: errcheck, gosec

		// Ожидание завершения приложения
		if err = lau.cmd.Wait(); err != nil {
			// Получение кода ошибки с которым завершилось выполнение команды
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					ret = status.ExitStatus()
				}
			}
		} else {
			ret = 0
		}
	}

	if lau.pid > 0 {
		if prc, err = os.FindProcess(lau.pid); err != nil {
			return
		}
		if pst, err = prc.Wait(); err != nil {
			return
		}
		if pst.Success() {
			ret = 0
		} else {
			var tmp []string
			var code int64
			tmp = strings.Split(pst.String(), " ")
			if len(tmp) > 0 {
				if code, err = strconv.ParseInt(tmp[len(tmp)-1], 10, 64); err == nil {
					ret = int(code)
				} else {
					err = nil
				}
			}
			lau.stderrWr.Sync()  // nolint: errcheck, gosec
			lau.stdoutWr.Sync()  // nolint: errcheck, gosec
			lau.stderr.Close()   // nolint: errcheck, gosec
			lau.stderrWr.Close() // nolint: errcheck, gosec
			lau.stdout.Close()   // nolint: errcheck, gosec
			lau.stdoutWr.Close() // nolint: errcheck, gosec
		}
	}

	return
}

//lint:ignore U1000 Используется в другой реинкарнации, оставлена на будущее
// Прогон через шаблонизатор команды и всех её аргументов для замены переменных на их значения
func (lau *impl) template(cmd []string, vars interface{}) (ret []string, err error) {
	var rsp *bytes.Buffer
	var tpl *template.Template
	var i int
	// Ловля паники, шаблонизатор паникует, сука, вместо нормального возвращения ошибки
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Template panic: %s", e.(error).Error())
		}
	}()
	for i = range cmd {
		tpl = template.Must(template.New(fmt.Sprintf("arg%d", i)).Parse(cmd[i]))
		rsp = bytes.NewBufferString(``)
		if err = tpl.Execute(rsp, vars); err != nil {
			return
		}
		ret = append(ret, rsp.String())
	}
	return
}
