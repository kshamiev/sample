package environment // impoert "application/components/environment"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"

	"application/workflow"
)

// Interface is an interface of package
type Interface workflow.ComponentInterface

// impl is an implementation of package
type impl struct {
	Ctx    context.Context
	CtxCfn context.CancelFunc
}
