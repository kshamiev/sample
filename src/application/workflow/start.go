package workflow // import "application/workflow"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"

	runtimeDebug "runtime/debug"
)

// Start Running start command in the all registered components
// if all of components ignored the command, runs the usageFunc
func (wfw *impl) Start(cmd string, usageFunc func()) (exitCode uint8, err error) {
	var done, help bool

	defer func() {
		if e := recover(); e != nil {
			exitCode, err = ErrCatchPanic, fmt.Errorf("%s\nGoroutine stack is:\n%s", e, string(runtimeDebug.Stack()))
			return
		}
	}()
	help = true
	for i := range wfw.Components {
		if singleton.debug {
			log.Debugf("Execute command %q in component %q", cmd, packageName(wfw.Components[i]))
		}
		done, exitCode, err = wfw.Components[i].Start(cmd)
		if exitCode != ErrNone || err != nil {
			err = fmt.Errorf(errText[exitCode], err)
			return
		}
		if done {
			help = false
			break
		}
	}
	if help {
		usageFunc()
	}

	return
}
