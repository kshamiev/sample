package goose // import "application/models/goose"

//go:generate db2struct create "" "goose_db_version" "goose // import \"application/models/goose\"" "DbVersion" "types_model.go"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"gopkg.in/webnice/kit.v1/modules/db"
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

// TableName ORM set default table name
func (cpt *DbVersion) TableName() string { return "goose_db_version" }
