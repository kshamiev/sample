package goose // import "application/models/goose"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

// New creates a new object and return interface
func New() Interface {
	var gse = new(impl)
	return gse
}

// CurrentVersion Возвращается текущая версия схемы базы данных
func (gse *impl) CurrentVersion() (ver *DbVersion, err error) {
	ver = new(DbVersion)
	err = gse.Gist().
		Where("`is_applied` = 1").
		Order("`tstamp` DESC").
		Last(ver).
		Error

	return
}
