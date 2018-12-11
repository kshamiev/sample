package cleaner // import "application/workers/cleaner"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"

	"application/configuration"

	"gopkg.in/webnice/job.v1/types"
)

// Interface is an interface of package
type Interface interface {
	types.WorkerInterface
}

// impl is an implementation of package
type impl struct {
	ID     string
	Ctx    context.Context
	CtxCfn context.CancelFunc
	Cfg    configuration.Interface
}
