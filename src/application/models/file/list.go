package file // import "application/models/file"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"io/ioutil"
	"os"
	"path"
)

// RecursiveFileList Поиск всех файлов начиная от path рекурсивно
// Возвращается слайс относительных имён файлов
func (fl *impl) RecursiveFileList(path string) (ret []string, err error) {
	ret, err = fl.RecursiveFileListLoop(path, "")
	return
}

// RecursiveFileListLoop Удобная для рекурсии функция
func (fl *impl) RecursiveFileListLoop(base, current string) (ret []string, err error) {
	var pf string
	var resp []string
	var fis []os.FileInfo
	var i int

	pf = path.Join(base, current)
	fis, err = ioutil.ReadDir(pf)
	for i = range fis {
		switch {
		case fis[i].IsDir():
			resp, err = fl.RecursiveFileListLoop(base, path.Join(current, fis[i].Name()))
			if err != nil {
				return
			}
			ret = append(ret, resp...)
		case !fis[i].IsDir():
			ret = append(ret, path.Join(current, fis[i].Name()))
		}
	}

	return
}
