package pages // import "application/models/pages"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"application/models/file"
	"application/models/filecache"
	pagesTypes "application/models/pages/types"
)

// Загрузка html шаблонов
func (pgm *impl) templateInit() (ret []*pagesTypes.Body, err error) {
	var mfi file.Interface
	var item *pagesTypes.Body
	var files []string
	var i int

	mfi = file.New()
	if files, err = mfi.RecursiveFileList(pgm.templatePages); err != nil {
		return
	}
	for i = range files {
		// Пропускаем лишние файлы
		if !rexFilenameFilter.MatchString(files[i]) {
			continue
		}
		if item, err = pgm.templateLoad(files[i]); err != nil {
			log.Criticalf("template load error: %s", err)
			err = nil
			continue
		}
		ret = append(ret, item)
	}

	return
}

// Загрузка файла шаблона
func (pgm *impl) templateLoad(f string) (ret *pagesTypes.Body, err error) {
	var fullname, source string
	var data filecache.Data
	var buf []string

	fullname = path.Join(pgm.templatePages, f)
	if data, err = pgm.Mfc.Load(fullname); err != nil {
		err = fmt.Errorf("load file %q error: %s", fullname, err)
		return
	}
	// Поиск окончательного файла, если файл является ссылкой
	if source, err = pgm.templateGetSourceLink(data.Reader()); err != nil {
		return
	}
	ret = &pagesTypes.Body{
		FullName:  fullname,
		Name:      data.Name(),
		IsInclude: rexIsInclude.MatchString(data.Name()),
		Source:    source,
	}
	// Выделение корня имени файла
	if buf = rexFilenameFilter.FindStringSubmatch(ret.Name); len(buf) > 1 {
		ret.NameBasis = buf[1]
	}
	// Определение роли шаблона
	switch ret.NameBasis {
	case keyIndex, keyLayout:
		ret.IsTemplate = true
	}
	// Выделение пути к шаблону или группе шаблонов
	ret.Path = strings.Replace(f, ret.Name, ``, -1)
	ret.Path = string(os.PathSeparator) + strings.TrimRight(ret.Path, string(os.PathSeparator))

	return
}

// Проверка наличия тега <!-- #source: layout.tpl.html -->
// Рекурсивный поиск конечного исходного файла с защитой от вечного цикла
func (pgm *impl) templateGetSourceLink(rd io.Reader) (ret string, err error) {
	var src, fname string
	var get func(rd io.Reader) (string, error)
	var infiniteRecursionProtector map[string]bool
	var data filecache.Data
	var rdr io.Reader
	var ok bool

	get = func(r io.Reader) (ret string, err error) {
		var buf string
		var tmp []string
		buf, err = bufio.NewReader(r).ReadString(0x0)
		if err != nil && err != io.EOF {
			return
		}
		err = nil
		if tmp = rexSourceLnk.FindStringSubmatch(buf); len(tmp) > 1 {
			ret = strings.TrimSpace(tmp[1])
		}
		return
	}
	// Поиск источника с защитой от вечного цикла
	infiniteRecursionProtector, rdr = make(map[string]bool), rd
	for {
		if src, err = get(rdr); err != nil {
			break
		}
		if src == "" {
			break
		} else if _, ok = infiniteRecursionProtector[src]; ok {
			err = fmt.Errorf("break infinite recursion link file %q", src)
			break
		}
		ret, infiniteRecursionProtector[src], fname =
			path.Join(pgm.templatePages, src), true, path.Join(pgm.templatePages, src)
		if data, err = pgm.Mfc.Load(fname); err != nil {
			err = fmt.Errorf("load file %q error: %s", fname, err)
			break
		}
		rdr = data.Reader()
	}

	return
}
