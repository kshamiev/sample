package logging // import "application/logging"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/workflow"

	r "gopkg.in/webnice/log.v2/receiver"
	s "gopkg.in/webnice/log.v2/sender"
)

// impl is an implementation of package
type impl struct {
	cfg configuration.Interface
}

func init() {
	workflow.Register(new(impl))
}

// Init Plug-in initialization function
func (cgn *impl) Init(appVersion string, appBuild string) (exitCode int, err error) {
	var receiver s.Receiver

	exitCode = workflow.ErrNone
	cgn.cfg = configuration.Get()
	if cgn.cfg.Log() == nil {
		return
	}
	if cgn.cfg.Log().GraylogEnable {
		receiver = r.GelfReceiver.
			SetAddress(cgn.cfg.Log().GraylogProto, cgn.cfg.Log().GraylogAddress, cgn.cfg.Log().GraylogPort).
			SetCompression("gzip").
			Receiver
		s.Gist().AddSender(receiver)
	}
	if cgn.cfg.Debug() {
		s.Gist().AddSender(r.Default.Receiver)
	}
	log.Gist().StandardLogSet()
	log.Info(`Logging system has initialized successfully`)

	return
}

// Command Processing the command
func (cgn *impl) Command(cmd string) (exitCode int, err error, done bool) {
	exitCode = workflow.ErrNone
	return
}

// Done The function is called when the application terminates
func (cgn *impl) Shutdown() {
	log.Gist().StandardLogUnset()
	log.Done()
}
