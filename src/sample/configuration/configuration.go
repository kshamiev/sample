package configuration

import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"

	"application/configuration"
	"application/workflow"
)

// impl is an implementation of package
type impl struct {
	conf configuration.Interface
}

func init() {
	workflow.Register(new(impl))
}

// Init Plug-in initialization function
func (cgn *impl) Init(appVersion string, appBuild string) (exitCode int, err error) {
	exitCode = workflow.ErrNone
	cgn.conf = configuration.Version(configuration.NewVersion(appVersion, appBuild))
	// Initializing the application configuration
	if err = cgn.conf.Init(); err != nil {
		exitCode = workflow.ErrInitConfiguration
		return
	}
	return
}

// Command Processing the command
func (cgn *impl) Command(cmd string) (exitCode int, err error, done bool) {
	exitCode = workflow.ErrNone
	// Print version and exit
	if cmd == `version` {
		fmt.Println(cgn.conf.Version().String())
		done = true
		return
	}
	// Go to the working directory
	if err = cgn.conf.GotoWorkingDirectory(); err != nil {
		err = fmt.Errorf("%q: %s", cgn.conf.WorkingDirectory(), err)
		exitCode, done = workflow.ErrCantChangeWorkDirectory, true
		return
	}
	// Only in the debug mode, the current configuration is displayed
	if cgn.conf.Debug() {
		log.Debug(debug.DumperString(cgn.conf.Configuration(), cgn.conf.Log()))
	}

	return
}

// Done The function is called when the application terminates
func (cgn *impl) Shutdown() {
}
