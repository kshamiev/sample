package paths // import "application/configuration/paths"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"os"
	"path/filepath"
	"strings"

	"application/configuration/osext"
)

// Initialize the list of ways the location of the configuration file
func SearchСonfigurationPaths() (pattern []string) {
	var path0, path1, path3, path4, path5, path6, path7, f, fn, p, ps string
	var arr []string

	arr = strings.Split(filepath.Base(os.Args[0]), `.`)
	arr[len(arr)-1] = configExtension[1:]
	fn = arr[len(arr)-2]
	f = `/` + strings.Join(arr, `.`)

	p, _ = osext.Executable()
	p = filepath.Dir(p)
	ps = "/"
	path0 = strings.Replace(p, `\`, `/`, -1)
	path1 = strings.Replace(filepath.Dir(path0), `\`, `/`, -1)
	path3 = strings.Replace(os.Getenv("SYSTEMROOT"), `\`, `/`, -1)
	path4 = strings.Replace(os.Getenv("PROGRAMFILES"), `\`, `/`, -1)
	path5 = strings.Replace(os.Getenv("COMMONPROGRAMFILES"), `\`, `/`, -1)
	path6 = strings.Replace(os.Getenv("PROGRAMDATA"), `\`, `/`, -1)
	path7 = strings.Replace(os.Getenv("LOCALAPPDATA"), `\`, `/`, -1)

	// In the folder with the executable file
	pattern = append(pattern, path1+f)

	// In the system folder of OS
	pattern = append(pattern, path3+f)

	// In the application folders
	pattern = append(pattern, path4+ps+fn+f)
	pattern = append(pattern, path5+ps+fn+f)
	pattern = append(pattern, path6+ps+fn+f)

	// In your application home folder
	pattern = append(pattern, path7+ps+fn+f)

	// The current folder
	pattern = append(pattern, path0+f)
	return
}
