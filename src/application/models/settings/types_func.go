package settings // import "application/models/settings"

//go:generate db2struct create "" "settings" "settings // import \"application/models/settings\"" "settings" "types_model.go"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	nul "gopkg.in/webnice/lin.v1/nl"

	"github.com/jinzhu/gorm"
)

// TableName ORM set default table name
func (stg *settings) TableName() string { return "settings" }

// BeforeCreate Функция вызывается до создания нового объекта в базе данных
func (stg *settings) BeforeCreate(scope *gorm.Scope) (err error) {
	var field *gorm.Field
	var ok bool

	if err = stg.BeforeUpdate(scope); err != nil {
		return
	}
	stg.CreateAt = nul.NewTimeValue(stg.UpdateAt.MustValue())
	stg.CreateAt.NullIfDefault()
	if field, ok = scope.FieldByName("CreateAt"); ok {
		if err = scope.SetColumn(field, stg.CreateAt); err != nil {
			return
		}
	}

	return
}

// BeforeUpdate Функция вызывается до обновления объекта в базе данных
func (stg *settings) BeforeUpdate(scope *gorm.Scope) (err error) {
	const column = `UpdateAt`
	var field *gorm.Field
	var ok bool

	stg.UpdateAt = nul.NewTimeValue(time.Now())
	if field, ok = scope.FieldByName(column); ok {
		if err = scope.SetColumn(field, stg.UpdateAt); err != nil {
			return
		}
	}

	return
}
