package workflow // import "application/workflow"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"fmt"
	"os"

	runtimeDebug "runtime/debug"
)

// Shutdown of all plugins
func (wfw *impl) Shutdown() {
	defer func() {
		if e := recover(); e != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\nGoroutine stack is:\n%s\n", e, string(runtimeDebug.Stack()))
			return
		}
	}()
	for i := len(wfw.plugins) - 1; i >= 0; i-- {
		wfw.plugins[i].Shutdown()
	}
}
