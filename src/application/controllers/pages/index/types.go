package index // import "application/controllers/pages/index"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	pagesTypes "application/models/pages/types"
)

var (
	indexURN = []string{
		`/`,
		`/index.htm`,
		`/index.html`,
	}
)

// Interface is an interface of package
type Interface interface {
	pagesTypes.Controller
}

// impl is an implementation of package
type impl struct {
	serverURL string // URL сервера
}
