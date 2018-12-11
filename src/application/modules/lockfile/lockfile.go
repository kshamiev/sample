package lockfile // import "application/modules/lockfile"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Lockfile is a pid file which can be locked
type Lockfile string

// TmpError is a type of error where a retry after a random amount of sleep should help to mitigate it.
type TmpError string

func (t TmpError) Error() string { return string(t) }

// Temporary returns always true.
// It exists, so you can detect it via
//	if te, ok := err.(interface{ Temporary() bool }); ok {
//		fmt.Println("I am a temporary error situation, so wait and retry")
//	}
func (t TmpError) Temporary() bool { return true }

// Various errors returned by this package
var (
	ErrBusy          = TmpError("Locked by other process")             // If you get this, retry after a short sleep might help
	ErrNotExist      = TmpError("Lockfile created, but doesn't exist") // If you get this, retry after a short sleep might help
	ErrNeedAbsPath   = errors.New("Lockfiles must be given as absolute path names")
	ErrInvalidPath   = errors.New("Lockfiles must be in existing folder")
	ErrInvalidPid    = errors.New("Lockfile contains invalid pid for system")
	ErrDeadOwner     = errors.New("Lockfile contains pid of process not existent on this system anymore")
	ErrRogueDeletion = errors.New("Lockfile owned by me has been removed unexpectedly")
)

// New describes a new filename located at the given absolute path.
func New(pth string) (lf Lockfile, err error) {
	if !filepath.IsAbs(pth) {
		return Lockfile(""), ErrNeedAbsPath
	}
	err = os.MkdirAll(path.Dir(pth), os.FileMode(0755))
	if err != nil {
		return Lockfile(""), ErrInvalidPath
	}
	lf = Lockfile(pth)
	return
}

// GetOwner returns who owns the lockfile.
func (l Lockfile) GetOwner() (proc *os.Process, err error) {
	var content []byte
	var pid int
	var running bool
	var name = string(l)

	// Ok, see, if we have a stale lockfile here
	content, err = ioutil.ReadFile(name) // nolint: errcheck, gosec
	if err != nil {
		return
	}
	// try hard for pids. If no pid, the lockfile is junk anyway and we delete it.
	pid, err = scanPidLine(content)
	if err != nil {
		return
	}
	running, err = isRunning(pid)
	if err != nil {
		return
	}
	if !running {
		err = ErrDeadOwner
		return
	}
	proc, err = os.FindProcess(pid)
	if err != nil {
		return
	}

	return

}

// TryLock tries to own the lock.
// It Returns nil, if successful and and error describing the reason, it didn't work out.
// Please note, that existing lockfiles containing pids of dead processes
// and lockfiles containing no pid at all are simply deleted.
func (l Lockfile) TryLock() (err error) {
	var tmplock *os.File
	var fiTmp, fiLock os.FileInfo
	var proc *os.Process
	var name = string(l)

	// This has been checked by New already. If we trigger here,
	// the caller didn't use New and re-implemented it's functionality badly.
	// So panic, that he might find this easily during testing.
	if !filepath.IsAbs(name) {
		panic(ErrNeedAbsPath)
	}

	tmplock, err = ioutil.TempFile(filepath.Dir(name), "")
	if err != nil {
		return
	}

	var cleanup = func() {
		tmplock.Close()           // nolint: errcheck, gosec
		os.Remove(tmplock.Name()) // nolint: errcheck, gosec
	}
	defer cleanup()

	if err = writePidLine(tmplock, os.Getpid()); err != nil {
		return
	}

	// return value intentionally ignored, as ignoring it is part of the algorithm
	os.Link(tmplock.Name(), name) // nolint: errcheck, gosec

	fiTmp, err = os.Lstat(tmplock.Name())
	if err != nil {
		return
	}
	fiLock, err = os.Lstat(name)
	if err != nil {
		// tell user that a retry would be a good idea
		if os.IsNotExist(err) {
			err = ErrNotExist
			return
		}
		return
	}

	// Success
	if os.SameFile(fiTmp, fiLock) {
		return nil
	}

	proc, err = l.GetOwner()
	switch err {
	default:
		// Other errors -> defensively fail and let caller handle this
		return
	case nil:
		if proc.Pid != os.Getpid() {
			err = ErrBusy
			return
		}
	case ErrDeadOwner, ErrInvalidPid:
		// cases we can fix below
	}

	// clean stale/invalid lockfile
	err = os.Remove(name)
	if err != nil {
		// If it doesn't exist, then it doesn't matter who removed it.
		if !os.IsNotExist(err) {
			return
		}
	}

	// now that the stale lockfile is gone, let's recurse
	err = l.TryLock()
	return
}

// Unlock a lock again, if we owned it. Returns any error that happened during release of lock.
func (l Lockfile) Unlock() error {
	proc, err := l.GetOwner()
	switch err {
	case ErrInvalidPid, ErrDeadOwner:
		return ErrRogueDeletion
	case nil:
		if proc.Pid == os.Getpid() {
			// we really own it, so let's remove it.
			return os.Remove(string(l))
		}
		// Not owned by me, so don't delete it.
		return ErrRogueDeletion
	default:
		// This is an application error or system error.
		// So give a better error for logging here.
		if os.IsNotExist(err) {
			return ErrRogueDeletion
		}
		// Other errors -> defensively fail and let caller handle this
		return err
	}
}

func writePidLine(w io.Writer, pid int) error {
	_, err := io.WriteString(w, fmt.Sprintf("%d\n", pid))
	return err
}

func scanPidLine(content []byte) (int, error) {
	if len(content) == 0 {
		return 0, ErrInvalidPid
	}

	var pid int
	if _, err := fmt.Sscanln(string(content), &pid); err != nil {
		return 0, ErrInvalidPid
	}

	if pid <= 0 {
		return 0, ErrInvalidPid
	}
	return pid, nil
}
