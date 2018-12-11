package filestore // import "application/models/filestore"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"
	"path"
	"time"

	"application/models/file"
	filestoreTypes "application/models/filestore/types"
)

// CleanOldData Очистка filestore от устаревших данных
func (ufm *impl) CleanOldData() (err error) {
	// Пометка на удаление файлов с возрастом больше 1 часа
	if err = ufm.cleanMarkDataTodelete(); err != nil {
		return
	}
	// Физическое удаление файлов и записей в БД с возрастом больше 1 суток
	// Удаляются только файлы о которых есть запись в БД
	if err = ufm.cleanData(); err != nil {
		return
	}
	// Удаление файлов во временной папке о которых нет записей в БД
	if err = ufm.cleanOrphanedTemporaryFiles(); err != nil {
		return
	}
	// Очистка пустых подпапок всего файлового хранилища
	if err = file.New().CleanEmptyFolder(ufm.storagePath); err != nil {
		return
	}

	return
}

// Пометка на удаление файлов с возрастом больше 1 часа
func (ufm *impl) cleanMarkDataTodelete() (err error) {
	var oneHour = time.Now().Add(0 - time.Hour)
	err = ufm.Gist().
		Model(&filestoreTypes.FilesTemporary{}).
		Where("`deleteAt` IS NULL").
		Where("`createAt` < ?", oneHour).
		Update("deleteAt", time.Now()).
		Error
	return
}

// Физическое удаление файлов и записей в БД с возрастом больше 1 суток
// Удаляются только файлы о которых есть запись в БД
func (ufm *impl) cleanData() (err error) {
	var oneDay time.Time
	var fl []*filestoreTypes.FilesTemporary
	var i int
	var fln string

	oneDay = time.Now().Add(0 - (time.Hour * 24))
	if err = ufm.Gist().
		Where("`deleteAt` IS NOT NULL").
		Where("`deleteAt` < ?", oneDay).
		Find(&fl).
		Error; err != nil {
		err = fmt.Errorf("database model error: %s", err)
		return
	}
	for i = range fl {
		fln = path.Join(ufm.storagePath, fl[i].LocalPath.MustValue())
		if err = os.RemoveAll(fln); err != nil {
			log.Warningf("remove %q error: %s", fln, err)
		}
		log.Noticef("removed temporary file: %q", fln)
	}
	err = ufm.Gist().
		Model(&filestoreTypes.FilesTemporary{}).
		Where("`deleteAt` IS NOT NULL").
		Where("`deleteAt` < ?", oneDay).
		Delete(&filestoreTypes.FilesTemporary{}).
		Error

	return
}

// Удаление файлов во временной папке о которых нет записей в БД
func (ufm *impl) cleanOrphanedTemporaryFiles() (err error) {
	var fListAll []string
	var fln string
	var count uint64
	var i int

	// Загрузка полного списока файлов во временной папке
	fListAll, err = file.New().RecursiveFileList(path.Join(ufm.storagePath, pathStorageTemporary))
	if err != nil {
		err = nil
		return
	}
	// Поиск загруженных файлов в базе данных
	for i = range fListAll {
		fln = path.Join(pathStorageTemporary, fListAll[i])
		if err = ufm.Gist().
			Model(&filestoreTypes.FilesTemporary{}).
			Where("`localPath` = ?", fln).
			Count(&count).
			Error; err != nil {
			log.Errorf("database model error: %s", err)
			return
		}
		if count > 0 {
			continue
		}
		// Файл потерян
		fln = path.Join(ufm.storagePath, fln)
		if err = os.RemoveAll(fln); err != nil {
			log.Warningf("can't remove file %q from filestore: %s", fln, err)
		} else {
			log.Noticef("deleted orphaned file %q from filestore", path.Join(pathStorageTemporary, fListAll[i]))
		}
	}

	return
}
