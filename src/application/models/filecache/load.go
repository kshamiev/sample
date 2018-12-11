package filecache // import "application/models/filecache"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"container/list"
	"fmt"
	stdMime "mime"
	"os"
	stddebug "runtime/debug"
	"time"

	"application/models/file"

	"gopkg.in/webnice/web.v1/mime"
)

// Load Функция загружает объект в память и предоставляет интерфейс доступа к объекту в режиме только для чтения
// Если объекта в памяти нет, то он запрашивается у функции CreateFn (функция по умолчанию - чтение объекта с диска)
func (mfc *impl) Load(name string) (Data, error) { return mfc.LoadWithLifetime(name, time.Duration(0)) }

// LoadWithLifetime Функция загружает объект в память и предоставляет интерфейс доступа к объекту в режиме только для чтения
// Если объекта в памяти нет, то он запрашивается у функции CreateFn (функция по умолчанию - чтение объекта с диска)
func (mfc *impl) LoadWithLifetime(name string, lifetime time.Duration) (Data, error) {
	var err error
	var elm *list.Element
	var mdo *memData
	var ok bool

	if mfc.debug {
		lifetime = time.Duration(-1)
	}
	// Protect concurent map access
	mfc.mhash.RLock()
	elm, ok = mfc.mhash.Map[name]
	mfc.mhash.RUnlock()
	// Объект найден в памяти
	if ok {
		mdo = elm.Value.(*memData)
		if !mdo.loaded {
			err = mfc.CreateMemoryObject(mdo, true)
		}
		return mdo, err
	}
	// Создание объекта в памяти путём чтения его с диска
	mdo = &memData{
		lifetime:  lifetime,
		fullName:  name,
		creatorFn: mfc.loadFile,
	}
	if err = mfc.CreateMemoryObject(mdo, true); err != nil {
		return mdo, err
	}
	// Добавление в список объектов в памяти, индексирование
	mfc.mlist.PushBack(mdo)
	defer mfc.index()

	return mdo, err
}

// CreateMemoryObject Создание объекта в памяти используя функцию создания
func (mfc *impl) CreateMemoryObject(mdo *memData, loadNow bool) (err error) {
	if mdo.fullName == "" {
		err = mfc.Errors().ErrNameNotSpecified()
		return
	}
	if loadNow && mdo.creatorFn != nil && !mdo.loaded {
		if err = mfc.loadBody(mdo); err != nil {
			return
		}
	} else if mdo.creatorFn == nil {
		err = mfc.Errors().ErrCreatorFnIsNil()
		return
	}

	return
}

// Загрузка тела объекта
func (mfc *impl) loadBody(mdo *memData) (err error) {
	var nwo *MemoryObject
	var n int

	// При высове внешней функции возможна паника!
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Panic recovery:\n%v\n%s", e.(error), string(stddebug.Stack()))
		}
	}()
	nwo, err = mdo.creatorFn(mdo.fullName)
	if err == os.ErrNotExist {
		err = mfc.Errors().ErrNotFound()
		return
	} else if err != nil {
		return
	}
	if nwo == nil {
		err = mfc.Errors().ErrCreatorFnReturnNil()
		return
	}
	mdo.body = make([]byte, len(nwo.Body))
	if n = copy(mdo.body, nwo.Body); n != len(nwo.Body) {
		err = mfc.Errors().ErrCopyDataFailed()
		return
	}
	if !mfc.debug {
		mdo.lifetime = nwo.Lifetime
	}
	mdo.createAt = time.Now()
	mdo.contentType = nwo.ContentType
	mdo.info = nwo.Info
	mdo.loaded = true

	return
}

// Функция создания объекта путём чтения файла с диска
func (mfc *impl) loadFile(name string) (ret *MemoryObject, err error) {
	var buf *bytes.Buffer
	var mfi file.Interface
	var fhi os.FileInfo
	var tmp []string

	mfi = file.New()
	buf, fhi, err = mfi.LoadFile(name)
	if os.IsNotExist(err) {
		err = mfc.Errors().ErrNotFound()
		return
	} else if err != nil {
		err = fmt.Errorf("Loading file %q error: %s", name, err)
		return
	}
	ret = &MemoryObject{
		Lifetime:  time.Duration(0),
		Body:      buf.Bytes(),
		Info:      fhi,
		CreatorFn: mfc.loadFile,
	}
	if mfc.debug {
		ret.Lifetime = time.Duration(-1)
	}
	// Определение Content-Type по расширению имени файла
	if tmp, ret.ContentType = extractExtension.FindStringSubmatch(name), mime.OctetStream; len(tmp) > 1 {
		ret.ContentType = stdMime.TypeByExtension("." + tmp[1])
	}

	return
}
