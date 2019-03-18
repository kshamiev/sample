package types // import "application/models/filestore/types"

//go:generate db2struct create "" "filestore" "types // import \"application/models/filestore/types\"" "Filestore" "filestore_model.go"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	"github.com/jinzhu/gorm"
)

// TableName ORM set default table name
func (fse *Filestore) TableName() string { return "filestore" }

// BeforeCreate Функция вызывается до создания нового объекта в базе данных
func (fse *Filestore) BeforeCreate(scope *gorm.Scope) (err error) {
	var field *gorm.Field
	var ok bool

	if err = fse.BeforeUpdate(scope); err != nil {
		return
	}
	fse.CreateAt.SetValid(fse.UpdateAt.MustValue())
	fse.CreateAt.NullIfDefault()
	if field, ok = scope.FieldByName("CreateAt"); ok {
		if err = scope.SetColumn(field, fse.CreateAt); err != nil {
			return
		}
	}
	fse.AccessAt.SetValid(fse.UpdateAt.MustValue())
	fse.AccessAt.NullIfDefault()
	if field, ok = scope.FieldByName("AccessAt"); ok {
		if err = scope.SetColumn(field, fse.AccessAt); err != nil {
			return
		}
	}

	return
}

// BeforeUpdate Функция вызывается до обновления объекта в базе данных
func (fse *Filestore) BeforeUpdate(scope *gorm.Scope) (err error) {
	var field *gorm.Field
	var ok bool

	fse.UpdateAt.SetValid(time.Now())
	if field, ok = scope.FieldByName("UpdateAt"); ok {
		if err = scope.SetColumn(field, fse.UpdateAt); err != nil {
			return
		}
	}

	return
}
