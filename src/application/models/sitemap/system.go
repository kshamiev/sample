package sitemap // import "application/models/sitemap"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Count Возвращает количество записей URN в sitemap
// Данные берутся из системной таблицы
func (smm *impl) Count() (ret uint64) {
	var err error
	var key, value []byte

	key = []byte(keyURLCount)
	if err = smm.db.View(func(tx *bolt.Tx) (err error) {
		value = tx.Bucket([]byte(keyBucketSystem)).Get(key)
		return
	}); err != nil {
		log.Errorf("Get %q from bucket %q error: %s", keyURLCount, keyBucketSystem, err)
		return
	}
	// Значение пустое
	if value == nil {
		return
	}
	if ret, err = strconv.ParseUint(string(value), 10, 64); err != nil {
		log.Errorf("parse count of URN records error: %s", err)
		return
	}

	return
}

// URLSetAt Функция возвращает дату и время последнего изменения в списке URL
// Данные берутся из системной таблицы
func (smm *impl) URLSetAt() (ret time.Time) {
	var err error
	var key, value []byte

	key = []byte(keyURLSetAt)
	if err = smm.db.View(func(tx *bolt.Tx) (err error) {
		value = tx.Bucket([]byte(keyBucketSystem)).Get(key)
		return
	}); err != nil {
		log.Errorf("Get %q from bucket %q error: %s", keyURLSetAt, keyBucketSystem, err)
		return
	}
	// Значение пустое
	if value == nil {
		return
	}
	if ret, err = time.ParseInLocation(time.RFC3339Nano, string(value), time.UTC); err != nil {
		log.Errorf("Time parse error: %s", err)
		return
	}

	return
}

// Обновление записи о количестве URN в sitemap в системной таблице
func (smm *impl) updateURLCount() (err error) {
	var count uint64

	// Подсчёт количества записей в key/value хранилище
	if err = smm.db.View(func(tx *bolt.Tx) (err error) {
		err = tx.Bucket([]byte(keyBucketURLSet)).
			ForEach(func(key, value []byte) (e error) {
				count++
				return
			})
		return
	}); err != nil {
		log.Warningf("Count records in bucket %q error: %s", keyBucketURLSet, err)
	}
	// Сохранение количества записей в системной табличке
	err = smm.db.Update(func(tx *bolt.Tx) (err error) {
		return tx.Bucket([]byte(keyBucketSystem)).
			Put([]byte(keyURLCount), []byte(strconv.FormatUint(count, 10)))
	})

	return
}

// Обновление даты и времени изменения URN в sitemap в системной таблице
func (smm *impl) updateURLSetAt(tm time.Time) (err error) {
	err = smm.db.Update(func(tx *bolt.Tx) (err error) {
		return tx.Bucket([]byte(keyBucketSystem)).
			Put([]byte(keyURLSetAt), []byte(tm.In(time.UTC).Format(time.RFC3339Nano)))
	})

	return
}
