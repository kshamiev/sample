// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package lockfile // import "application/modules/lockfile"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"os"
	"syscall"
)

func isRunning(pid int) (ret bool, err error) {
	var proc *os.Process
	if proc, err = os.FindProcess(pid); err != nil {
		return
	}
	if err = proc.Signal(syscall.Signal(0)); err != nil {
		return
	}
	ret = true
	return
}
