package webserver // import "application/webserver"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"path"

	"application/models/file"
	"application/models/filecache"
)

// WarmingUpCache Прогрев файлового кеша
func (wso *impl) CacheWarmingUp(rootPath string) (err error) {
	var mfi file.Interface
	var mfc filecache.Interface
	var files []string
	var fnm string
	var i int

	mfi, mfc = file.New(), filecache.Get().Debug(wso.debug)
	if wso.debug {
		log.Infof("Cache warming up for path %q", rootPath)
	}
	if files, err = mfi.RecursiveFileList(rootPath); err != nil {
		return
	}
	for i = range files {
		fnm = path.Join(rootPath, files[i])
		_, err = mfc.Load(fnm)
		if err != nil {
			log.Warningf("Cache warming up for file %q error: %s", fnm, err)
		}
	}
	if wso.debug {
		log.Infof("Summary cache size is %d bytes", mfc.Size())
	}

	return
}
