package types // import "application/models/myitem/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/nul.v1/nul"
import (
	"time"

	"github.com/jinzhu/gorm"
)

// Myitem Какая-то сущность бизнес логики
type Myitem struct {
	ID       uint64    `gorm:"column:id;primary_key;              AUTO_INCREMENT;UNSIGNED;NOT NULL;type:BIGINT(20)"` // Уникальный идентификатор записи
	CreateAt nul.Time  `gorm:"column:createAt;                    NULL;DEFAULT NULL;type:DATETIME"`                  // Дата и время создания записи
	UpdateAt nul.Time  `gorm:"column:updateAt;                    NULL;DEFAULT NULL;type:DATETIME"`                  // Дата и время обновления записи
	DeleteAt nul.Time  `gorm:"column:deleteAt;                    INDEX;NULL;DEFAULT NULL;type:DATETIME"`            // Дата и время удаления записи (пометка на удаление)
	AccessAt nul.Time  `gorm:"column:accessAt;                    NULL;DEFAULT NULL;type:DATETIME"`                  // Дата и время последней авторизации или обновления токена доступа
	Date     time.Time `gorm:"column:date;                        NOT NULL;type:DATETIME"`                           // Любая дата и время
	Number   int64     `gorm:"column:number;                      NOT NULL;type:BIGINT(20)"`                         // Любое число
	Text     string    `gorm:"column:text;                        NOT NULL;type:TINYTEXT"`                           // Любая строка
}

// TableName ORM set default table name
func (mit *Myitem) TableName() string { return "myitem" }

// BeforeCreate Функция вызывается до создания нового объекта в базе данных
func (mit *Myitem) BeforeCreate(scope *gorm.Scope) (err error) {
	var field *gorm.Field
	var ok bool

	if err = mit.BeforeUpdate(scope); err != nil {
		return
	}
	mit.CreateAt.SetValid(mit.UpdateAt.MustValue())
	mit.CreateAt.NullIfDefault()
	if field, ok = scope.FieldByName("CreateAt"); ok {
		if err = scope.SetColumn(field, mit.CreateAt); err != nil {
			return
		}
	}
	mit.AccessAt.SetValid(mit.UpdateAt.MustValue())
	mit.AccessAt.NullIfDefault()
	if field, ok = scope.FieldByName("AccessAt"); ok {
		if err = scope.SetColumn(field, mit.AccessAt); err != nil {
			return
		}
	}

	return
}

// BeforeUpdate Функция вызывается до обновления объекта в базе данных
func (mit *Myitem) BeforeUpdate(scope *gorm.Scope) (err error) {
	var field *gorm.Field
	var ok bool

	mit.UpdateAt.SetValid(time.Now())
	if field, ok = scope.FieldByName("UpdateAt"); ok {
		if err = scope.SetColumn(field, mit.UpdateAt); err != nil {
			return
		}
	}

	return
}
