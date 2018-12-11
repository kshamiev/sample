package sitemap // import "application/models/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	msgpack "github.com/vmihailenco/msgpack"
	bolt "go.etcd.io/bbolt"
)

// Add Добавление URN в sitemap
func (smm *impl) Add(items []*Location) (Interface, error) {
	var err error
	var i, max int

	for i = 0; len(items) > 0; i++ {
		if len(items) > bolt.DefaultMaxBatchSize {
			max = bolt.DefaultMaxBatchSize
		} else {
			max = len(items)
		}
		if err = smm.add(items[:max]); err != nil {
			return smm, err
		}
		if len(items) >= bolt.DefaultMaxBatchSize {
			items = items[bolt.DefaultMaxBatchSize:]
		} else {
			items = items[len(items):]
		}
	}
	// Обновление Count записи в системной табличке
	if err = smm.updateURLCount(); err != nil {
		err = fmt.Errorf("update URLCount error: %s", err)
	}

	return smm, err
}

// Добавление среза URN в sitemap, срез не более чем DefaultMaxBatchSize
func (smm *impl) add(items []*Location) (err error) {
	var last time.Time

	if err = smm.db.Batch(func(tx *bolt.Tx) (err error) {
		var bt *bolt.Bucket
		var key, value []byte
		var i int

		bt = tx.Bucket([]byte(keyBucketURLSet))
		for i = range items {
			if last.IsZero() || items[i].ModTime.After(last) {
				last = items[i].ModTime
			}
			if key, value, err = smm.makeKeyValue(items[i]); err != nil {
				err = fmt.Errorf("make key/value error: %s", err)
				return
			}
			if err = bt.Put(key, value); err != nil {
				err = fmt.Errorf("put key %q to bucket %q error: %s", key, keyBucketURLSet, err)
				return
			}
		}

		return
	}); err != nil {
		err = fmt.Errorf("add item to sitemap error: %s", err)
		return
	}
	// Обновление даты и времени изменения в списке URL
	if err = smm.updateURLSetAt(last); err != nil {
		err = fmt.Errorf("update URLSetAt error: %s", err)
	}

	return
}

// Создание key/value из структуры URN
func (smm *impl) makeKeyValue(item *Location) (key []byte, value []byte, err error) {
	var tmp [32]byte

	key = make([]byte, 32)
	tmp = sha256.Sum256([]byte(item.URN))
	_ = copy(key, tmp[:])
	if value, err = msgpack.Marshal(item); err != nil {
		return
	}

	return
}

// Декодирование URN из value
func (smm *impl) decodeValue(value []byte) (item *Location, err error) {
	item = new(Location)
	err = msgpack.Unmarshal(value, item)
	return
}

// GetByRange Выгружает URN в соответствии с порядковыми номерами ключей
// Функция вернёт данные начиная с [from] элемента, массив не более чем [count] элементов
func (smm *impl) GetByRange(from uint64, count uint64) (ret []*Location, err error) {
	var current uint64
	var bt *bolt.Bucket
	var cr *bolt.Cursor
	var key, value []byte
	var item *Location

	ret = make([]*Location, 0, count)
	err = smm.db.View(func(tx *bolt.Tx) (err error) {
		bt = tx.Bucket([]byte(keyBucketURLSet))
		cr = bt.Cursor()
		for key, value = cr.First(); key != nil; key, value = cr.Next() {
			if current < from {
				current++
				continue
			}
			if item, err = smm.decodeValue(value); err != nil {
				err = fmt.Errorf("decode value from bucket %q error: %s", keyBucketURLSet, err)
				return
			}
			ret = append(ret, item)
			current++
			if len(ret) >= int(count) {
				return
			}
		}
		return
	})

	return
}

// SitemapXMLWriteTo Возвращает sitemap.xml, с данными начиная с позиции from в количестве count
func (smm *impl) SitemapXMLWriteTo(wr io.Writer, url string, from uint64, count uint64) (err error) {
	var rsp *Sitemap
	var items []*Location
	var item *Location
	var part *URLPart
	var xmlEnc *xml.Encoder

	rsp = &Sitemap{
		XMLNS: XMLns,
		URL:   make([]*URLPart, 0, count),
	}
	if items, err = smm.GetByRange(from, count); err != nil {
		err = fmt.Errorf("sitemap GetByRange() error: %s", err)
		return
	}
	if len(items) == 0 {
		err = smm.Errors().ErrNotFound()
		return
	}
	for _, item = range items {
		part = &URLPart{
			Loc:        url + item.URN,
			LastMod:    item.ModTime,
			ChangeFreq: item.Change,
			Priority:   item.Priority,
		}
		rsp.URL = append(rsp.URL, part)
	}
	xmlEnc = xml.NewEncoder(wr)
	if smm.debug {
		xmlEnc.Indent("", "  ")
	}
	if err = xmlEnc.Encode(rsp); err != nil {
		err = fmt.Errorf("marshal xml error: %s", err)
		return
	}

	return
}

// SitemapXML Возвращает sitemap.xml с данными начиная с позиции from в количестве count
func (smm *impl) SitemapXML(url string, from uint64, count uint64) (ret []byte, err error) {
	var buf = &bytes.Buffer{}
	if err = smm.SitemapXMLWriteTo(buf, url, from, count); err != nil {
		return
	}
	ret = buf.Bytes()

	return
}

// SitemapIndexXMLWriteTo Возвращает sitemap-index.xml
func (smm *impl) SitemapIndexXMLWriteTo(wr io.Writer, url string, idx uint64) (err error) {
	var rsp *Index
	var maxItems, fromItems, countItems, from, count, n uint64
	var i float64
	var xmlEnc *xml.Encoder

	// -----------------------------------------------------------------------------------------------------------------
	// Элементы sitemap в проекции на sitemap-index
	// Пример:
	// * [            0 -        49'999] - sitemap-index.xml или sitemap-index-0.xml в котором один sitemap.xml
	// * [       50'000 - 2'499'999'999] - sitemap-index-1.xml в котором 50'000 записей sitemap-0.xml ... sitemap-49999.xml
	// * [2'500'000'000 - 4'999'999'999] - sitemap-index-2.xml в котором 50'000 записей sitemap-50000.xml ... sitemap-99999.xml
	// Максимальное возможное количество URN в БД sitemap: 18'446'744'073'709'551'615 это sitemap-index-7378697629.xml
	// -----------------------------------------------------------------------------------------------------------------
	if idx >= MaxIndex/MaxRecords {
		err = smm.Errors().ErrTooLarge()
		return
	}
	// Получение количества URN в базе данных
	maxItems, fromItems = smm.Count(), idx*MaxRecords*MaxRecords
	if fromItems > maxItems || maxItems == 0 {
		err = smm.Errors().ErrNotFound()
		return
	}
	// Расчёт количества элементов sitemap.xml в sitemap-index.xml
	if countItems = maxItems - fromItems; countItems > (MaxRecords * MaxRecords) {
		countItems = MaxRecords * MaxRecords
	}
	if i = float64(countItems) / float64(MaxRecords); i-float64(uint64(i)) > 0.0 {
		i++
	}
	from, count = idx*MaxRecords, uint64(i)
	rsp = &Index{
		XMLNS:   XMLns,
		Sitemap: make([]*IndexPart, 0, MaxRecords),
	}
	for n = 0; n < count; n++ {
		rsp.Sitemap = append(rsp.Sitemap, &IndexPart{
			Loc:     fmt.Sprintf("%s/sitemap-%d.xml", url, from+n),
			LastMod: smm.URLSetAt(),
		})
	}
	xmlEnc = xml.NewEncoder(wr)
	if smm.debug {
		xmlEnc.Indent("", "  ")
	}
	if err = xmlEnc.Encode(rsp); err != nil {
		err = fmt.Errorf("marshal xml error: %s", err)
		return
	}

	return
}

// SitemapIndexXML Возвращает sitemap-index.xml
func (smm *impl) SitemapIndexXML(url string, idx uint64) (ret []byte, err error) {
	var buf = &bytes.Buffer{}
	if err = smm.SitemapIndexXMLWriteTo(buf, url, idx); err != nil {
		return
	}
	ret = buf.Bytes()

	return
}

// Del Удаление URN из sitemap
func (smm *impl) Del(urn string) (Interface, error) {
	var err error

	if err = smm.db.Update(func(tx *bolt.Tx) (err error) {
		var key []byte

		if key, _, err = smm.makeKeyValue(&Location{URN: urn}); err != nil {
			err = fmt.Errorf("make key/value error: %s", err)
			return
		}
		err = tx.Bucket([]byte(keyBucketURLSet)).
			Delete(key)

		return
	}); err != nil {
		err = fmt.Errorf("Delete item from sitemap error: %s", err)
		return smm, err
	}
	// Обновление Count записи в системной табличке
	if err = smm.updateURLCount(); err != nil {
		err = fmt.Errorf("update URLCount error: %s", err)
	}

	return smm, err
}
