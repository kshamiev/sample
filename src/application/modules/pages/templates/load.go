package templates // include "application/modules/pages/templates"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"gopkg.in/webnice/kit.v1/models/file"
)

// LoadPages Загрузка html страниц
func (pgs *impl) LoadPages() {
	var fle file.Interface
	var files []string
	var fib *TplBody
	var i int

	if pgs.PagesRoot == "" {
		log.Info(`HTML templates folder is not available`)
		return
	}
	// Загрузка списка всех файлов
	fle = file.New()
	if files, pgs.Err = fle.RecursiveFileList(pgs.PagesRoot); pgs.Err != nil {
		log.Errorf("Error reading list of html pages templates '%s': %s", pgs.PagesRoot, pgs.Err.Error())
		return
	}
	// Получение информации по всем файлам, группировка по URN
	for i = range files {
		if fib, pgs.Err = pgs.LoadFileBodyInfo(files[i]); pgs.Err != nil {
			return
		}
		// Загрузка данных из файлов в память
		// С проверкой линковки файлов с помощью тега: <!-- #source: layout.tpl.html -->
		if err := pgs.LoadFileBodyData(fib); err != nil {
			log.Errorf("Error LoadFile('%s'): %s", fib.FullPath, err.Error())
			continue
		}
		pgs.pushFileInfo(fib)
	}
	log.Info(`Templates html pages has loaded successfully`)
}
