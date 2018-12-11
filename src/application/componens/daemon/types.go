package daemon // import "application/componens/daemon"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/workflow"

	"gopkg.in/webnice/job.v1/job"
)

// Interface is an interface of package
type Interface workflow.ComponentInterface

// impl is an implementation of package
type impl struct {
	Cfg configuration.Interface
	Jbo job.Interface
}
