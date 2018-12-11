package workflow // import "application/workflow"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"sort"
	"strings"

	runtimeDebug "runtime/debug"
)

// Init Initialize all registered components
func (wfw *impl) Init(appVersion string, appBuild string) (exitCode uint8, err error) {
	defer func() {
		if e := recover(); e != nil {
			exitCode, err = ErrCatchPanic, fmt.Errorf("%s\nGoroutine stack is:\n%s", e, string(runtimeDebug.Stack()))
			return
		}
	}()
	// Сортировка компонентов для удовлетворения зависимостей
	wfw.sortRegistredComponents()
	// DEBUG
	//log.Debug(debug.DumperString(wfw.Components))
	// DEBUG
	if singleton.debug {
		log.Debugf("Workflow sorted all application components: %v", wfw.Components)
	}
	// Выполеннение инициализации компонентов
	for i := range wfw.Components {
		if singleton.debug {
			log.Debugf("Init component %q", packageName(wfw.Components[i]))
		}
		exitCode, err = wfw.Components[i].Init(appVersion, appBuild)
		if exitCode != ErrNone {
			err = fmt.Errorf(errText[exitCode], err)
			return
		}
	}

	return
}

// Сортировка зарегистрированных компонентов в соответствии с заявленными зависимостями
func (wfw *impl) sortRegistredComponents() {
	var tmp []string
	var mp map[string]int64
	var pkgname string
	var i, j int
	var ok bool

	mp = make(map[string]int64)
	for i = range wfw.Components {
		mp[packageName(wfw.Components[i])] = 0
	}
	for i = range wfw.Components {
		pkgname, tmp = packageName(wfw.Components[i]), wfw.Components[i].After()
		for j = range tmp {
			if _, ok = mp[tmp[j]]; ok {
				mp[pkgname] += int64(len(mp)) + mp[tmp[j]]
			}
		}
	}
	sort.Slice(wfw.Components, func(i, j int) (ret bool) {
		pgI, pgJ := packageName(wfw.Components[i]), packageName(wfw.Components[j])
		if mp[pgI] == mp[pgJ] {
			ret = strings.Compare(pgI, pgJ) == -1
		} else {
			ret = mp[pgI] < mp[pgJ]
		}
		return
	})
}
