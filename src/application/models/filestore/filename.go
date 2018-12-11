package filestore

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"
	"path"
	"time"
)

// Создание пути и имени файла
func (ufm *impl) makeFilename() (ret string) {
	var tm = time.Now().In(time.Local)

	ret = path.Join(
		fmt.Sprintf("%04d%c%02d%c%02d", tm.Year(), os.PathSeparator, tm.Month(), os.PathSeparator, tm.Day()),
		fmt.Sprintf("%020d", tm.UnixNano()),
	)

	return
}

// Создание пути и имени файла к временному хранилищу
func (ufm *impl) makeTemporaryFilename() (pathFull string, pathRelative string) {
	pathRelative = path.Join(pathStorageTemporary, ufm.makeFilename())
	pathFull = path.Join(ufm.storagePath, pathRelative)
	return
}

// Создание пути и имени файла к постоянному хранилищу
func (ufm *impl) makePermanentFilename() (pathFull string, pathRelative string) {
	pathRelative = path.Join(pathStoragePermanent, ufm.makeFilename())
	pathFull = path.Join(ufm.storagePath, pathRelative)
	return
}
