package goose // import "application/models/goose"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"gopkg.in/webnice/kit.v1/modules/db"
	nul "gopkg.in/webnice/lin.v1/nl"
)

// Interface is an interface of package
type Interface interface {
	// CurrentVersion Возвращается текущая версия схемы базы данных
	CurrentVersion() (ver *DbVersion, err error)
}

// impl is an implementation of package
type impl struct {
	db.Implementation
}

// DbVersion Лог миграций проекта
type DbVersion struct {
	ID        uint64   `gorm:"column:id;primary_key" sql:"UNSIGNED;NOT NULL;AUTO_INCREMENT;type:BIGINT(20)"` // Уникальный идентификатор записи
	VersionID uint64   `gorm:"column:version_id"     sql:"NOT NULL;type:BIGINT(20)"`                         // Уникальный идентификатор миграции
	IsApplied bool     `gorm:"column:is_applied"     sql:"NOT NULL;type:TINYINT(1)"`                         // =true - миграция применена, =false - миграция не применена
	TimeStamp nul.Time `gorm:"column:tstamp"         sql:"NULL;DEFAULT CURRENT_TIMESTAMP;type:TIMESTAMP"`    // Дата и время применения миграции
}

// TableName ORM set default table name
func (cpt *DbVersion) TableName() string { return "goose_db_version" }
