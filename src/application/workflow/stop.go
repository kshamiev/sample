package workflow // import "application/workflow"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"

	runtimeDebug "runtime/debug"
)

// Stop of all registered components
func (wfw *impl) Stop() (exitCode uint8, err error) {
	defer func() {
		if e := recover(); e != nil {
			exitCode, err = ErrCatchPanic, fmt.Errorf("%s\nGoroutine stack is:\n%s", e, string(runtimeDebug.Stack()))
			return
		}
	}()
	for i := len(wfw.Components) - 1; i >= 0; i-- {
		if singleton.debug {
			log.Debugf("Stop component %q", packageName(wfw.Components[i]))
		}
		exitCode, err = wfw.Components[i].Stop()
		if exitCode != ErrNone || err != nil {
			log.Fatalf(errText[exitCode], err)
		}
	}

	return
}
