package workflow // import "application/workflow"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"fmt"

	runtimeDebug "runtime/debug"
)

// Initialize all registered plugins
func (wfw *impl) Initialize(appVersion string, appBuild string) (exitCode int, err error) {
	defer func() {
		if e := recover(); e != nil {
			exitCode, err = ErrCatchPanic, fmt.Errorf("%s\nGoroutine stack is:\n%s", e, string(runtimeDebug.Stack()))
			return
		}
	}()

	for i := range wfw.plugins {
		exitCode, err = wfw.plugins[i].Init(appVersion, appBuild)
		if exitCode != ErrNone {
			err = fmt.Errorf(errText[exitCode], err)
			return
		}
	}

	return
}
