// +build darwin dragonfly freebsd linux netbsd openbsd plan9 solaris

package paths // import "application/configuration/paths"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"os"
	"os/user"
	"path/filepath"

	"application/configuration/osext"
)

// SearchСonfigurationPaths Initialize the list of ways the location of the configuration file
func SearchСonfigurationPaths() (pattern []string) {
	var err error
	var pt, filename, currentFolder string
	var parent, home string
	var ps = string(os.PathSeparator)
	var usr *user.User

	pt, _ = osext.Executable() // nolint: errcheck, gosec
	filename = filepath.Base(pt)
	parent = filepath.Clean(filepath.Dir(pt) + ps + `..`)
	usr, err = user.Current()
	if err != nil {
		home = os.Getenv(`HOME`)
	} else {
		home = usr.HomeDir
	}

	// In the folder with the executable file
	pattern = append(pattern, pt+configExtension)

	// The parent folder of the executable file
	pattern = append(pattern, parent+ps+filename+configExtension)

	// The conf directory is located in the parent folder of the executable file
	pattern = append(pattern, parent+ps+`conf`+ps+filename+configExtension)

	// In the /etc folder
	pattern = append(pattern, `/etc`+ps+filename+configExtension)

	// In the folder [filename] is located in /etc
	pattern = append(pattern, `/etc`+ps+filename+ps+filename+configExtension)

	// In the folder [filename] is located in /opt
	pattern = append(pattern, `/opt`+ps+filename+ps+filename+configExtension)

	// In the folder [filename] is located in /usr/local/etc
	pattern = append(pattern, `/usr/local/etc`+ps+filename+configExtension)

	// In your ~/.config/[filename] is located in your home folder
	pattern = append(pattern, home+ps+`.config`+ps+filename+ps+filename+configExtension)

	// The current folder
	currentFolder, err = os.Getwd()
	if err == nil {
		if currentFolder[len(currentFolder)-1] != '/' {
			currentFolder += ps
		}
		pattern = append(pattern, currentFolder+filename+configExtension)
	}
	return
}
