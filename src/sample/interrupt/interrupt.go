package interrupt

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"os"

	"application/modules/interrupt"
	"application/workflow"

	"gopkg.in/webnice/job.v1/job"
)

// impl is an implementation of package
type impl struct {
	itp interrupt.Interface
}

func init() {
	workflow.Register(new(impl))
}

// Init Plug-in initialization function
func (cgn *impl) Init(appVersion string, appBuild string) (exitCode int, err error) {

	exitCode = workflow.ErrNone

	return
}

// Command Processing the command
func (cgn *impl) Command(cmd string) (exitCode int, err error, done bool) {
	exitCode = workflow.ErrNone
	cgn.itp = interrupt.New(func(sig os.Signal) {
		log.Alertf(`Received OS interrupt signal`+": %s", sig.String())
		job.Cancel()
	})
	cgn.itp.Start()
	log.Info(`Application interrupt protection has initialized successfully`)

	return
}

// Done The function is called when the application terminates
func (cgn *impl) Shutdown() {
	if cgn.itp == nil {
		return
	}
	_ = cgn.itp.Stop()
}
