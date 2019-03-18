package types // import "application/models/myitem/types"

//go:generate db2struct create "" "myitem" "types // import \"application/models/myitem/types\"" "Myitem" "myitem_model.go"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	"github.com/jinzhu/gorm"
)

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
