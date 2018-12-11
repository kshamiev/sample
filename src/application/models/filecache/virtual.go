package filecache // import "application/models/filecache"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"
)

// Virtual Добавление виртуального файла в кеш, файл создаётся путём вызова функции
// После создания файла, он запоминается в памяти на время указанное в lifetime
func (mfc *impl) Virtual(name string, lifetime time.Duration, createFn CreateFn) (err error) {
	var mdo *memData

	if mfc.debug {
		lifetime = time.Duration(-1)
	}
	mdo = &memData{
		createAt:  time.Now(),
		lifetime:  lifetime,
		fullName:  name,
		creatorFn: createFn,
	}
	if err = mfc.CreateMemoryObject(mdo, false); err != nil {
		return
	}
	// Добавление в список объектов в памяти, индексирование
	mfc.mlist.PushBack(mdo)
	defer mfc.index()

	return
}
