package pages // import "application/controllers/resource/pages"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"bufio"
	"fmt"
	"net/url"
	"path"
	"strings"

	"application/models/filecache"
)

// Загрузка и разбор синтаксиса файла route.txt
func (pgi *impl) routeLoad() (ret []string, err error) {
	var tmp string
	var urnLines []string
	var i int

	if urnLines, err = pgi.routeReadContent(); err != nil || len(urnLines) == 0 {
		return
	}
	ret = make([]string, 0, len(urnLines))
	for i = range urnLines {
		if !pgi.routeCheckURN(urnLines[i]) {
			if pgi.debug {
				log.Warningf("route urn %q error", urnLines[i])
			}
			continue
		}
		if tmp, err = pgi.routeTransformToWebURN(urnLines[i]); err != nil {
			err = fmt.Errorf("URN %q error: %s", urnLines[i], err)
			return
		}
		ret = append(ret, tmp)
	}

	return
}

// Загрузка только полезных строк файла route.txt
func (pgi *impl) routeReadContent() (ret []string, err error) {
	var data filecache.Data
	var scanner *bufio.Scanner
	var fn string
	var buf []string

	fn = path.Join(pgi.rootPath, keyRouteTxt)
	switch data, err = pgi.Mfc.Load(fn); err {
	case pgi.Mfc.Errors().ErrNotFound():
		if err = nil; pgi.debug {
			log.Noticef("File %q not found", fn)
		}
		return
	default:
		if err != nil {
			err = fmt.Errorf("load file %q error: %s", fn, err)
			return
		}
	}
	scanner = bufio.NewScanner(data.Reader())
	for scanner.Scan() {
		if buf = rexRouteTxtLine.FindStringSubmatch(scanner.Text()); len(buf) > 2 {
			ret = append(ret, strings.TrimSpace(buf[2]))
		}
	}

	return
}

// Проверка URN адреса
func (pgi *impl) routeCheckURN(urn string) (ret bool) {
	const separator = `/`
	var tmp []string
	var i int

	tmp = strings.Split(urn, separator)
	for i = range tmp {
		if len(tmp[i]) == 0 && i > 0 {
			return
		}
		if tmp[i] == `*` && i == 1 {
			return
		}
		if tmp[i] == `?` && i == 1 && len(tmp) == 2 {
			return
		}
	}
	ret = true

	return
}

// Преобразование URN описанного в route.txt в формат gopkg.in/webnice/web.v1/route
func (pgi *impl) routeTransformToWebURN(urn string) (ret string, err error) {
	const (
		section   = `:section%d`
		separator = `/`
	)
	var base, u *url.URL
	var tmp []string
	var i int

	if base, err = url.Parse(pgi.serverURL); err != nil {
		err = fmt.Errorf("parse URL %q error: %s", pgi.serverURL, err)
		return
	}
	if u, err = url.Parse(urn); err != nil {
		err = fmt.Errorf("parse urn %q error: %s", urn, err)
		return
	}
	tmp = strings.Split(base.ResolveReference(u).RequestURI(), separator)
	for i = range tmp {
		if tmp[i] == `?` {
			tmp[i] = fmt.Sprintf(section, i)
		}
	}
	ret = strings.Join(tmp, separator)

	return
}
