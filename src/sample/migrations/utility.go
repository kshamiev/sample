package migrations

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"fmt"
	"os/exec"

	"application/modules/launcher"

	//"gopkg.in/webnice/kit.v1/modules/ch"
	"gopkg.in/webnice/kit.v1/modules/db"
)

// Поиск и проверка версии утилиты миграции
func (cgn *impl) migrationsUtility() (ret string) {
	const command = `gsmigrate`
	var err error
	if ret, err = exec.LookPath(command); err != nil {
		log.Warningf("Can't find migrations utility: %s", err.Error())
		return
	}
	return
}

func (cgn *impl) migrationsMysql(command string) (err error) {
	var dsn string

	// Если не указана папка с миграциями, то выход
	if cgn.conf.Configuration().Database.Migrations == "" {
		log.Warningf("Folder with mysql database migration files is not specified. Migrations are not not applied")
		return
	}
	if dsn, err = db.New().Dsn(); err != nil {
		err = fmt.Errorf("Database configuration error: %s", err.Error())
		return
	}
	// Применение миграций
	if err = cgn.migrationsApply(command, cgn.conf.Configuration().Database.Migrations, cgn.conf.Configuration().Database.Driver, dsn); err != nil {
		log.Warningf("Migrations warnings: %s", err.Error())
		return
	}

	return
}

func (cgn *impl) migrationsClickhouse(command string) (err error) {
	//	var dsn string
	//
	//	// Если не указана папка с миграциями, то выход
	//	if cgn.conf.Configuration().Clickhouse.Migrations == "" {
	//		log.Warningf("Folder with clickhouse database migration files is not specified. Migrations are not not applied")
	//		return
	//	}
	//	if dsn, err = ch.New().Dsn(); err != nil {
	//		err = fmt.Errorf("Database configuration error: %s", err.Error())
	//		return
	//	}
	//	// Применение миграций
	//	if err = cgn.migrationsApply(command, cgn.conf.Configuration().Clickhouse.Migrations, "clickhouse", dsn); err != nil {
	//		log.Warningf("Migrations warnings: %s", err.Error())
	//		return
	//	}

	return
}

// Примерение миграций
func (cgn *impl) migrationsApply(command string, dir string, drv string, dsn string) (err error) {
	var lau launcher.Interface
	var cmd, env []string
	var ecode int
	var oBuf, eBuf *bytes.Buffer

	env = []string{`GOOSE_DIR=` + dir, `GOOSE_DRV=` + drv, `GOOSE_DSN=` + dsn}
	cmd = []string{command, `up`}
	oBuf, eBuf = &bytes.Buffer{}, &bytes.Buffer{}
	lau = launcher.New()
	if err = lau.Launch(cmd, env, nil, oBuf, eBuf); err != nil {
		return
	}
	if ecode, err = lau.Wait(); err != nil || ecode != 0 {
		err = fmt.Errorf("Utility %q exit with error code %d: %s", command, ecode, err.Error())
	}
	if oBuf.Len() > 0 {
		log.Noticef("Migration utility (out): %q", oBuf.String())
	}
	if eBuf.Len() > 0 {
		log.Warningf("Migration utility (err): %q", eBuf.String())
	}

	return
}
