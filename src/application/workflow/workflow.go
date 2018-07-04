package workflow // import "application/workflow"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import ()

// Get return interface of package
func Get() Interface {
	if singleton == nil {
		singleton = new(impl)
	}
	return singleton
}

// Register Registration of plug-in
func Register(p PluginInterface) { Get().(*impl).Register(p) }

// Register Registration of plug-in
func (wfw *impl) Register(p PluginInterface) {
	if p == nil {
		return
	}
	wfw.plugins = append(wfw.plugins, p)
}
