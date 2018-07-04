package pidfile // import "application/modules/pidfile"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"os"
	"strings"

	"application/modules/lockfile"
)

const (
	_ProcesFinished = `process already finished`
)

// Interface is an interface of package
type Interface interface {
	// Lock tries to own the lock
	Lock() error
	// Unlock a lock again, if we owned it
	Unlock() error
}

// impl is an implementation of package
type impl struct {
	FileName string
	Error    error
	File     lockfile.Lockfile
}

// New Create new object
func New(fileName string) Interface {
	var pidf = new(impl)
	pidf.FileName = fileName
	pidf.File, pidf.Error = lockfile.New(pidf.FileName)
	return pidf
}

// Lock tries to own the lock
func (pidf *impl) Lock() (err error) {
	if pidf.Error != nil {
		err = pidf.Error
		return
	}

	err = pidf.File.TryLock()
	if err != nil {
		if strings.Contains(err.Error(), _ProcesFinished) {
			err = pidf.Delete()
			err = pidf.File.TryLock()
		}
	}
	return
}

// Unlock a lock again, if we owned it. Returns any error that happend during release of lock.
func (pidf *impl) Unlock() (err error) {
	if pidf.Error != nil {
		err = pidf.Error
		return
	}
	err = pidf.File.Unlock()
	return
}

// Delete pid file
func (pidf *impl) Delete() (err error) {
	err = os.Remove(pidf.FileName)
	return
}
