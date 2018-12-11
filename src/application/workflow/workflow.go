package workflow // import "application/workflow"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"

// Get return interface of package
func Get() Interface {
	if singleton == nil {
		singleton = new(impl)
	}
	return singleton
}

// Register Registration of application component
func Register(p ComponentInterface) { Get().(*impl).Register(p) }

// Debug Enable or disable debug mode
func (wfw *impl) Debug(d bool) Interface { wfw.debug = d; return wfw }

// Register Registration of application component
func (wfw *impl) Register(p ComponentInterface) {
	if p == nil {
		if singleton.debug {
			log.Debugf("Register workflow application nil component. Action missed")
		}
		return
	}
	if singleton.debug {
		log.Debugf("Register workflow application component %q", packageName(p))
	}
	wfw.Components = append(wfw.Components, p)
}
