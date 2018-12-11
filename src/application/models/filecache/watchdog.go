package filecache // import "application/models/filecache"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"time"
)

// Горутина очистки памяти от объектов с просроченным временем жизни
func (mfc *impl) watchdog() {
	const tmc = time.Second
	var tic = time.NewTicker(tmc)
	defer tic.Stop()
	for {
		select {
		case <-mfc.wdctx.Done():
			return
		case <-tic.C:
			mfc.watchdogCleaner()
		}
	}
}

// Просмотр и пометка просроченных объектов в кеше в памяти
func (mfc *impl) watchdogCleaner() {
	var elm *list.Element
	var item *memData
	var reindex bool

	for elm = mfc.mlist.Front(); elm != nil; elm = elm.Next() {
		item = elm.Value.(*memData)
		if item.lifetime == time.Duration(0) {
			continue
		}
		if item.lifetime > time.Duration(0) && time.Since(item.createAt.Add(item.lifetime)) < 0 {
			item.lifetime, reindex = -1, true
		}
		if item.lifetime == time.Duration(-1) {
			item.loaded = false
		}
	}
	if !reindex {
		return
	}
	mfc.index()
}
