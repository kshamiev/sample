// +build windows

package lockfile // import "application/modules/lockfile"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"syscall"
)

// For some reason these consts don't exist in syscall.
const (
	errorInvalidParameter = 87
	codeStillActive       = 259
)

func isRunning(pid int) (ret bool, err error) {
	var code uint32
	var procHnd syscall.Handle
	procHnd, err = syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, true, uint32(pid))
	if err != nil {
		if scerr, ok := err.(syscall.Errno); ok {
			if uintptr(scerr) == errorInvalidParameter {
				err = nil
				return
			}
		}
	}
	if err = syscall.GetExitCodeProcess(procHnd, &code); err != nil {
		return
	}
	ret = (code == codeStillActive)
	return
}
