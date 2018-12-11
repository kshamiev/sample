package types // import "application/models/filestore/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	nul "gopkg.in/webnice/lin.v1/nl"

	"github.com/jinzhu/gorm"
)

// Filestore Реестр файлов постоянного хранения
type Filestore struct {
	ID          uint64     `gorm:"column:id;primary_key;              AUTO_INCREMENT;UNSIGNED;NOT NULL;type:BIGINT(20)"` // Уникальный идентификатор записи
	CreateAt    nul.Time   `gorm:"column:createAt;                    INDEX;NULL;DEFAULT NULL;type:DATETIME"`            // Дата и время создания записи
	UpdateAt    nul.Time   `gorm:"column:updateAt;                    NULL;DEFAULT NULL;type:DATETIME"`                  // Дата и время обновления записи
	DeleteAt    nul.Time   `gorm:"column:deleteAt;                    INDEX;NULL;DEFAULT NULL;type:DATETIME"`            // Дата и время удаления записи (пометка на удаление)
	AccessAt    nul.Time   `gorm:"column:accessAt;                    NULL;DEFAULT NULL;type:DATETIME"`                  // Дата и время последнего доступа к записи
	FileName    nul.String `gorm:"column:filename;                    NULL;DEFAULT NULL;type:VARCHAR(4096)"`             // Оригинальное имя файла
	FileExt     nul.String `gorm:"column:fileExt;                     INDEX;NULL;DEFAULT NULL;type:VARCHAR(256)"`        // Расширение имени файла без точки
	Size        uint64     `gorm:"column:size;                        UNSIGNED;NOT NULL;DEFAULT '0';type:BIGINT(20)"`    // Размер файла в байтах
	Sha512      nul.String `gorm:"column:sha512;                      NULL;DEFAULT NULL;type:VARCHAR(128)"`              // SHA512 контрольная сумма файла в HEX формате
	LocalPath   nul.String `gorm:"column:localPath;                   NULL;DEFAULT NULL;type:VARCHAR(4096)"`             // Относительный путь и имя файла
	ContentType nul.String `gorm:"column:contentType;                 NULL;DEFAULT NULL;type:TEXT"`                      // MIME Content-Type загруженного файла
}

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
