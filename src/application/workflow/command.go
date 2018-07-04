package workflow // import "application/workflow"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"fmt"

	runtimeDebug "runtime/debug"
)

// Command Running the command with registered plugins
func (wfw *impl) Command(cmd string) (exitCode int, err error) {
	var done, help bool

	defer func() {
		if e := recover(); e != nil {
			exitCode, done, err = ErrCatchPanic, true, fmt.Errorf("%s\nGoroutine stack is:\n%s", e, string(runtimeDebug.Stack()))
			return
		}
	}()

	help = true
	for i := range wfw.plugins {
		exitCode, err, done = wfw.plugins[i].Command(cmd)
		if exitCode != ErrNone {
			err = fmt.Errorf(errText[exitCode], err)
			return
		}
		if done {
			help = false
			break
		}
	}
	if help {
		fmt.Println("Please use --help for help")
	}

	return
}
