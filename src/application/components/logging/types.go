package logging // import "application/components/logging"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"application/configuration"
	"application/workflow"
)

// Interface is an interface of package
type Interface workflow.ComponentInterface

// impl is an implementation of package
type impl struct {
	Cfg configuration.Interface
}
