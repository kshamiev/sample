package myitem

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	myitemTypes "application/models/myitem/types"

	"gopkg.in/webnice/kit.v1/modules/db"
)

// Interface is an interface of package
type Interface interface {
	// Create Создание новой сущности в БД
	Create(d time.Time, n int64, t string) (*myitemTypes.Myitem, error)

	// Status Получение состояния данных в БД
	Status() (*StatusInfo, error)

	// Load Загрузка сущности из БД
	Load(id uint64) (*myitemTypes.Myitem, error)

	// Delete Удаление сущности (пометка удалено, без физического удаления)
	Delete(id uint64) error

	// Clean Очистка талицы от старых данных
	Clean() error

	// ERRORS

	// ErrNotFound Not found
	ErrNotFound() error
}

// impl is an implementation of package
type impl struct {
	db.Implementation
}

// StatusInfo Состояние данных в БД
type StatusInfo struct {
	FirstDate time.Time // Дата и время создания первой записи
	LastDate  time.Time // Дата и время создания последней записи
	Size      int       // Кодичество записей
	Ids       []uint64  // Идентификаторы записей
}
