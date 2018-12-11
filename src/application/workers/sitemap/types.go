package sitemap // import "application/workers/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"gopkg.in/webnice/job.v1/types"
)

// Interface is an interface of package
type Interface types.WorkerInterface

// impl is an implementation of package
type impl struct {
	ID string // Уникальный идентификатор процесса
}
