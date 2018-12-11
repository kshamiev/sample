package settings // import "application/models/settings"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"time"

	nul "gopkg.in/webnice/lin.v1/nl"

	"github.com/jinzhu/gorm"
)

// Хранение настроек с типизированными значениями
type settings struct {
	ID           uint64      `gorm:"column:id;primary_key;        AUTO_INCREMENT;UNSIGNED;NOT NULL;type:BIGINT(20)"` // Уникальный идентификатор записи
	CreateAt     nul.Time    `gorm:"column:createAt;              INDEX;NULL;DEFAULT NULL;type:DATETIME"`            // Дата и время создания записи
	UpdateAt     nul.Time    `gorm:"column:updateAt;              INDEX;NULL;DEFAULT NULL;type:DATETIME"`            // Дата и время обновления записи
	AccessAt     nul.Time    `gorm:"column:accessAt;              INDEX;NULL;DEFAULT NULL;type:DATETIME"`            // Дата и время последнего доступа к записи
	Key          string      `gorm:"column:key;                   INDEX;NOT NULL;type:VARCHAR(255)"`                 // Ключ
	ValueString  nul.String  `gorm:"column:valueString;           NULL;DEFAULT NULL;type:LONGTEXT"`                  // Строковое значение
	ValueDate    nul.Time    `gorm:"column:valueDate;             NULL;DEFAULT NULL;type:DATETIME"`                  // Значение даты и времени
	ValueUint    nul.Uint64  `gorm:"column:valueUint;             NULL;DEFAULT NULL;UNSIGNED;type:BIGINT(20)"`       // Числовое unsigned значение
	ValueInt     nul.Int64   `gorm:"column:valueInt;              NULL;DEFAULT NULL;type:BIGINT(20)"`                // Числовое значение
	ValueDecimal nul.Float64 `gorm:"column:valueDecimal;          NULL;DEFAULT NULL;type:DECIMAL(16,8)"`             // Десятичное значение
	ValueFloat   nul.Float64 `gorm:"column:valueFloat;            NULL;DEFAULT NULL;type:DOUBLE"`                    // IEEE-754 64-bit floating-point number
	ValueBit     nul.Bool    `gorm:"column:valueBit;              NULL;DEFAULT NULL;type:BIT(1)"`                    // Boolean value
	ValueBlob    nul.Bytes   `gorm:"column:valueBlob;             NULL;DEFAULT NULL;type:LONGBLOB"`                  // Blob value
}

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
	var field *gorm.Field
	var ok bool

	stg.UpdateAt = nul.NewTimeValue(time.Now())
	if field, ok = scope.FieldByName("UpdateAt"); ok {
		if err = scope.SetColumn(field, stg.UpdateAt); err != nil {
			return
		}
	}

	return
}

// AfterFind Функция вызывается после чтения объекта из базы данных
func (stg *settings) AfterFind(scope *gorm.Scope) {
	var err error
	var tmp *settings
	var tm nul.Time

	if stg.ID == 0 {
		return
	}
	tm, tmp = nul.NewTimeValue(time.Now()), new(settings)
	tmp.ID, tmp.Key = stg.ID, stg.Key
	if err = scope.DB().
		Model(tmp).
		UpdateColumn("accessAt", tm).
		Error; err != nil {
		log.Error(err.Error())
		return
	}
	stg.AccessAt = tm
}
