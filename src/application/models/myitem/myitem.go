package myitem

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	myitemTypes "application/models/myitem/types"
)

// New creates a new object and return interface
func New() Interface {
	var mim = new(impl)
	return mim
}

// Create Создание новой сущности в БД
func (mim *impl) Create(d time.Time, n int64, t string) (ret *myitemTypes.Myitem, err error) {
	ret = new(myitemTypes.Myitem)
	ret.Date, ret.Number, ret.Text = d, n, t
	err = mim.Gist().
		Save(ret).
		Error
	return
}

// Status Получение состояния данных в БД
func (mim *impl) Status() (ret *StatusInfo, err error) {
	var items []*myitemTypes.Myitem
	var i int

	if err = mim.Gist().
		Select("`id`, `createAt`").
		Where("`deleteAt` IS NULL").
		Order("`createAt` ASC").
		Find(&items).
		Error; err != nil {
		return
	}
	ret = &StatusInfo{
		Ids: make([]uint64, 0, len(items)),
	}
	for i = range items {
		ret.Ids = append(ret.Ids, items[i].ID)
		ret.Size++
		if i == 0 {
			ret.FirstDate = items[i].CreateAt.MustValue()
		}
		ret.LastDate = items[i].CreateAt.MustValue()
	}

	return
}

// Load Загрузка сущности из БД
func (mim *impl) Load(id uint64) (ret *myitemTypes.Myitem, err error) {
	ret = new(myitemTypes.Myitem)
	if mim.Gist().
		Where("`deleteAt` IS NULL").
		First(ret, id).
		RecordNotFound() {
		ret, err = nil, mim.ErrNotFound()
		return
	}
	return
}

// Delete Удаление сущности (пометка удалено, без физического удаления)
func (mim *impl) Delete(id uint64) (err error) {
	var item = &myitemTypes.Myitem{}

	if mim.Gist().
		Where("`deleteAt` IS NULL").
		First(item, id).
		RecordNotFound() {
		err = mim.ErrNotFound()
		return
	}
	item.DeleteAt.SetValid(time.Now().In(time.Local))
	err = mim.Gist().
		Model(item).
		UpdateColumn("deleteAt", item.DeleteAt).
		Error

	return
}

// Clean Очистка талицы от старых данных
func (mim *impl) Clean() (err error) {
	var dd time.Time
	var item = &myitemTypes.Myitem{}

	// Удаление всех записей помеченных на удаление больше 1 недели назад
	dd = time.Now().In(time.Local).Add(0 - (time.Hour * 24 * 7))
	err = mim.Gist().
		Model(item).
		Where("NOT `deleteAt` IS NULL").
		Where("`deleteAt` < ?", dd).
		Delete(item).
		Error

	return
}
