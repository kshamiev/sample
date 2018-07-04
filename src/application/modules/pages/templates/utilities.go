package templates // include "application/modules/pages/templates"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/webnice/kit.v1/models/file"
)

// LoadFileBodyInfo Загрузка информации о файле
func (pgs *impl) LoadFileBodyInfo(fileName string) (ret *TplBody, err error) {
	var fi os.FileInfo
	var fn string

	fn = path.Join(pgs.PagesRoot, fileName)
	if fi, err = os.Stat(fn); err != nil {
		return
	}
	ret = new(TplBody)
	ret.FullPath = fn
	ret.ModTime = fi.ModTime()
	ret.Name = file.New().GetFileName(fn)
	ret.NameBasis = pgs.getWordBase(ret.Name)
	ret.Path = `/` + strings.Replace(fileName, ret.Name, ``, -1)
	ret.IsInclude = rexIsInclude.MatchString(ret.Name)
	ret.Size = fi.Size()
	switch ret.NameBasis {
	case "index", "layout":
		ret.IsTemplate = true
	}
	return
}

// Выделение основания имени файла из полного имени файла
func (pgs *impl) getWordBase(fn string) string { return rexExtension.ReplaceAllString(fn, ``) }

// Сохранение в переменной
func (pgs *impl) pushFileInfo(fib *TplBody) {
	var i int
	var urn string
	var found bool

	urn = fib.Path
	if !fib.IsTemplate && !fib.IsInclude {
		urn = pgs.urnNormal(fib.Path + fib.NameBasis)
	}
	for i = range pgs.Tpls.Tpl {
		if pgs.Tpls.Tpl[i].UrnAddress == urn {
			found = true

			// Если был ранее добавлен одиночный файл
			if len(pgs.Tpls.Tpl[i].MapData) == 1 {
				var key string
				for key = range pgs.Tpls.Tpl[i].MapData {
				}
				if !pgs.Tpls.Tpl[i].MapData[key].IsTemplate && !pgs.Tpls.Tpl[i].MapData[key].IsInclude {
					log.Warningf("Ignore template file '%s', using folder same name",
						pgs.Tpls.Tpl[i].MapData[key].Path+pgs.Tpls.Tpl[i].MapData[key].Name,
					)
					delete(pgs.Tpls.Tpl[i].MapData, key)
				}
			}
			if fib.IsTemplate || fib.IsInclude {
				pgs.Tpls.Tpl[i].MapData[fib.NameBasis] = fib
			} else {
				log.Warningf("Ignore template file '%s', using folder same name", fib.Path+fib.Name)
			}

		}
	}
	if !found {
		pgs.Tpls.Tpl = append(pgs.Tpls.Tpl, &tpl{
			UrnAddress: urn,
			MapData:    map[string]*TplBody{fib.NameBasis: fib},
		})
	}
}

// LoadFileBodyData Загрузка тела файла с диска в память
func (pgs *impl) LoadFileBodyData(fib *TplBody) (err error) {
	var fle file.Interface
	var buf *bytes.Buffer
	var src string
	var fi os.FileInfo
	var infiniteRecursion map[string]bool
	var ok bool

	// Продакшн режим, если файл считан то не перечитываем его
	if !pgs.Debug && fib.Data != nil {
		return
	}

	// Девелоперский режим, если файл не изменился, то не перечитываем его
	if pgs.Debug && fib.Data != nil {
		if fi, err = os.Stat(fib.FullPath); err != nil {
			return
		}
		if fib.ModTime.Equal(fi.ModTime()) && fib.Size == fi.Size() {
			return
		}
	}

	// Чтение файла
	fle = file.New()
	src = fib.FullPath
	infiniteRecursion = make(map[string]bool)
	for src != "" {
		// Check infinite recursion
		if _, ok = infiniteRecursion[src]; ok {
			err = fmt.Errorf("Break infinite recursion")
			return
		}
		infiniteRecursion[src] = true

		if buf, err = fle.LoadFile(src); err != nil {
			return
		}
		fib.FullPath = src
		if src = pgs.getSourceLink(buf); src != "" {
			src = path.Join(pgs.PagesRoot, src)
		} else {
			fib.Data = buf
		}
	}

	// Обновление FileInfo
	if fi, err = os.Stat(fib.FullPath); err != nil {
		return
	}
	fib.Size = fi.Size()
	fib.ModTime = fi.ModTime()

	return
}

// Проверка наличия тэга <!-- #source: layout.tpl.html --> и возврат имени файла если тэг есть
func (pgs *impl) getSourceLink(buf *bytes.Buffer) (ret string) {
	var tmp []string
	if tmp = rexSourceLnk.FindStringSubmatch(buf.String()); len(tmp) > 1 {
		ret = rexSpaceFrst.ReplaceAllString(rexSpaceLast.ReplaceAllString(tmp[1], ``), ``)
	}
	return
}

// Нормализация URN адреса
func (pgs *impl) urnNormal(urn string) (ret string) {
	ret = rexSlashLast.ReplaceAllString(urn, ``) + `/`
	ret = `/` + rexSlashFrst.ReplaceAllString(ret, ``)
	return
}

// Поиск запрошенных шаблонов
func (pgs *impl) getTemplates(tps []string) (ret *tpls, err error) {
	var i, n int
	var found bool
	var urn string

	ret = new(tpls)
	for i = range tps {
		found = false
		urn = pgs.urnNormal(tps[i])
		for n = range pgs.Tpls.Tpl {
			if pgs.Tpls.Tpl[n].UrnAddress == urn {
				found = true
				ret.Tpl = append(ret.Tpl, pgs.Tpls.Tpl[n])
			}
		}
		if !found {
			err = fmt.Errorf("Template for urn '%s' not found", urn)
		}
	}
	return
}
