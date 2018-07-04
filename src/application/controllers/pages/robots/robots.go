package robots // import "application/controllers/pages/robots"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import (
	"bytes"
	"net/http"

	"application/configuration"
	"application/modules/rendering"
)

// New creates a new object and return interface
func New() Interface {
	var rbs = new(impl)
	return rbs
}

// Lazy initialization
func (rbs *impl) init() {
	if rbs.cfg == nil {
		rbs.cfg = configuration.Get()
	}
	if rbs.url == "" {
		// first web server
		for i := range rbs.cfg.Configuration().WEBServers {
			rbs.url = rbs.cfg.Configuration().WEBServers[i].Server.Address
			rbs.DocumentRoot = rbs.cfg.Configuration().WEBServers[i].DocumentRoot
			break
		}
	}
}

// RobotsTxt robots.txt
// GET /robots.txt
func (rbs *impl) RobotsTxt(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var v *templateVars
	var rend rendering.Interface
	var buf *bytes.Buffer

	rbs.init()
	v = &templateVars{
		RequestDomain: rq.Host,
		RequestURL:    rbs.url,
	}
	buf = &bytes.Buffer{}
	rend = rendering.New(rendering.Option{
		Directory: rbs.DocumentRoot,
		Reload:    rbs.cfg.Debug(),
	})
	err = rend.NewStandardTemplateRender().RenderHTML(buf, v, _TemplateName)
	if err != nil {
		log.Errorf("RobotsTxt template error: %s", err.Error())
		wr.WriteHeader(status.InternalServerError)
		return
	}
	wr.Header().Set(header.ContentType, mime.TextPlainCharsetUTF8)
	wr.WriteHeader(status.Ok)
	if _, err = buf.WriteTo(wr); err != nil {
		log.Errorf("Response error: %s", err.Error())
		return
	}
}
