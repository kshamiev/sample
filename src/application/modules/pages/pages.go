package pages // import "application/modules/pages"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	stdMime "mime"
	"path"
	"strings"
)

// New creates a new object and return interface
func New(documentRoot string) Interface {
	var pgs = new(impl)
	pgs.DocumentRoot = documentRoot
	return pgs
}

// Index Выгрузка отпределённых типов файлов из document root
func (pgs *impl) Index(rfn string) (ret *File, err error) {
	var tmp []string

	if rfn == "" {
		rfn = _IndexFile
	}
	rfn = path.Join(pgs.DocumentRoot, rfn)
	tmp = _RootFilesPattern.FindStringSubmatch(rfn)
	if strings.Index(rfn, pgs.DocumentRoot) < 0 || len(tmp) == 0 {
		err = pgs.ErrFileNotFound()
		return
	}
	ret, err = pgs.loadCached(rfn)
	ret.ContentType = stdMime.TypeByExtension("." + tmp[1])

	return
}

// Assets Выгрузка всех файлов из папки assets
func (pgs *impl) Assets(rfn string) (ret *File, err error) {
	const _AssetsPreffix = `assets`

	if i := strings.Index(rfn, "?"); i > 0 {
		rfn = rfn[:i]
	}
	rfn = path.Join(pgs.DocumentRoot, rfn)
	if strings.Index(rfn, path.Join(pgs.DocumentRoot, _AssetsPreffix)) != 0 {
		err = pgs.ErrFileNotFound()
		return
	}
	ret, err = pgs.loadCached(rfn)
	if tmp := _ExtractExtension.FindStringSubmatch(rfn); len(tmp) > 0 {
		ret.ContentType = stdMime.TypeByExtension("." + tmp[1])
	}

	return
}
