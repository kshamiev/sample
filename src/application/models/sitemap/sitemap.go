package sitemap // import "application/models/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

// New creates a new object and return interface
func New(cachePath string) (Interface, error) {
	var smm = &impl{
		dbPath: path.Join(cachePath, keyDb),
		dbSync: new(sync.Mutex),
	}
	return smm, smm.initDB()
}

// Errors Ошибки известного состояни, которые могут вернуть функции пакета
func (smm *impl) Errors() *Error { return Errors() }

// Инициализация встроенной базы данных
func (smm *impl) initDB() (err error) {
	var tx *bolt.Tx

	// Защита от гонки
	smm.dbSync.Lock()
	defer smm.dbSync.Unlock()
	// DEBUG
	//log.Noticef(" - OPEN %q", smm.dbPath)
	// DEBUG
	// Открытие БД
	smm.db, err = bolt.Open(smm.dbPath, os.FileMode(0660), nil)
	if err != nil {
		err = fmt.Errorf("open database %q error: %s", smm.dbPath, err)
		return
	}
	if tx, err = smm.db.Begin(true); err != nil {
		err = fmt.Errorf("begin transaction error: %s", err)
		return
	}
	// Создание ведра для URL
	if _, err = tx.CreateBucketIfNotExists([]byte(keyBucketURLSet)); err != nil {
		err = fmt.Errorf("create bucket %q error: %s", keyBucketURLSet, err)
		return
	}
	// Создание ведра для системных переменных
	if _, err = tx.CreateBucketIfNotExists([]byte(keyBucketSystem)); err != nil {
		err = fmt.Errorf("create bucket %q error: %s", keyBucketSystem, err)
		return
	}
	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("commit transaction error: %s", err)
		return
	}

	return
}

// Debug Set debug mode
func (smm *impl) Debug(d bool) Interface { smm.debug = d; return smm }

// Close model
func (smm *impl) Close() (err error) {
	// Защита от гонки
	smm.dbSync.Lock()
	defer smm.dbSync.Unlock()
	// DEBUG
	//log.Noticef(" - CLOSE %q", smm.dbPath)
	// DEBUG
	// На всякий случай делаем синхронизацию
	if err = smm.db.Sync(); err != nil {
		err = fmt.Errorf("Sync database error: %s", err)
		return
	}
	err = smm.db.Close()

	return
}

// Reset delete all data and create new database
func (smm *impl) Reset() (ret Interface, err error) {
	// Закрытие БД
	ret, err = smm, smm.Close()
	if err != nil {
		err = fmt.Errorf("Close database error: %s", err)
		return
	}
	// Удаление с защитой от гонки
	smm.dbSync.Lock()
	err = os.Remove(smm.dbPath)
	smm.dbSync.Unlock()
	if err != nil {
		err = fmt.Errorf("Delete file %q error: %s", smm.dbPath, err)
		return
	}
	// Создание новой БД
	if err = smm.initDB(); err != nil {
		err = fmt.Errorf("init database error: %s", err)
		return
	}
	// Обновление даты и времени изменения в списке URL
	if err = smm.updateURLSetAt(time.Now()); err != nil {
		err = fmt.Errorf("update URLSetAt error: %s", err)
		return
	}
	// Обновление Count записи в системной табличке
	if err = smm.updateURLCount(); err != nil {
		err = fmt.Errorf("update URLCount error: %s", err)
	}

	return
}
