package filecache // import "application/models/filecache"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"context"
	"sort"
	"strings"
)

// Get return interface of singleton object
func Get() Interface {
	if singleton == nil {
		singleton = newObject()
	}
	return singleton
}

// Debug Set debug mode
func (mfc *impl) Debug(d bool) Interface { mfc.debug = d; return mfc }

// Errors Ошибки известного состояни, которые могут вернуть функции пакета
func (mfc *impl) Errors() *Error { return Errors() }

// Создание нового объекта пакета
func newObject() (mfc *impl) {
	mfc = &impl{
		mlist: list.New(),
		mhash: &memHash{
			Map: make(map[string]*list.Element),
		},
	}
	mfc.wdctx, mfc.wdcfn = context.WithCancel(context.Background())
	go mfc.watchdog()

	return
}

// Пересоздание индекса
func (mfc *impl) index() {
	var elm *list.Element
	var item *memData
	var fixErrors []*list.Element
	var i int
	var ok bool

	mfc.mhash.Lock()
	defer mfc.mhash.Unlock()
	mfc.mhash.Map = make(map[string]*list.Element)
	for elm = mfc.mlist.Front(); elm != nil; elm = elm.Next() {
		item = elm.Value.(*memData)
		if _, ok = mfc.mhash.Map[item.fullName]; !ok {
			mfc.mhash.Map[item.fullName] = elm
		} else {
			fixErrors = append(fixErrors, mfc.mhash.Map[item.fullName])
			mfc.mhash.Map[item.fullName] = elm
		}
	}
	if len(fixErrors) <= 0 {
		return
	}
	// Удаление ошибок
	for i = range fixErrors {
		mfc.mlist.Remove(fixErrors[i])
	}
}

// List Возвращает список объектов в памяти
func (mfc *impl) List() (ret []string) {
	var elm *list.Element
	var item *memData

	ret = make([]string, 0, mfc.mlist.Len())
	for elm = mfc.mlist.Front(); elm != nil; elm = elm.Next() {
		item = elm.Value.(*memData)
		ret = append(ret, item.fullName)
	}
	sort.SliceStable(ret, func(i, j int) bool { return strings.Compare(ret[i], ret[j]) == -1 })

	return
}

// Size Возвращает суммарный размер всех объектов в памяти
func (mfc *impl) Size() (ret uint64) {
	var elm *list.Element
	var item *memData

	for elm = mfc.mlist.Front(); elm != nil; elm = elm.Next() {
		item = elm.Value.(*memData)
		ret += item.Length()
	}

	return
}

// IsExist Проверка существования файла в кеше
func (mfc *impl) IsExist(name string) (ok bool) {
	mfc.mhash.RLock()
	_, ok = mfc.mhash.Map[name]
	mfc.mhash.RUnlock()
	return
}
