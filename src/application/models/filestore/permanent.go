package filestore // import "application/models/filestore"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io"
	"os"
	"path"

	"application/models/file"

	filestoreTypes "application/models/filestore/types"

	nul "gopkg.in/webnice/lin.v1/nl"
)

// NewPermanentFileFromTemporaryFile Создание постоянного файла из временного путём копирования
func (ufm *impl) NewPermanentFileFromTemporaryFile(fileID uint64) (ret *filestoreTypes.Filestore, err error) {
	var ifh io.ReadCloser
	var ofh *os.File
	var ift *filestoreTypes.FilesTemporary
	var pathFull, pathRelative, sha512sum string
	var size int64

	if ift, ifh, err = ufm.TemporaryFileOpen(fileID); err != nil {
		return
	}
	defer ifh.Close() // nolint: errcheck
	// Создание имени файла в постоянном хранилище
	pathFull, pathRelative = ufm.makePermanentFilename()
	// Создание директории к файлу
	if err = os.MkdirAll(path.Dir(pathFull), os.FileMode(0750)); err == os.ErrPermission {
		err = fmt.Errorf("make directory %q error: %s", path.Dir(pathFull), err)
		return
	}
	// Создание файла
	if ofh, err = os.OpenFile(pathFull, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0640)); err != nil {
		err = fmt.Errorf("create file %q error: %s", pathFull, err)
		return
	}
	defer ofh.Close() // nolint: errcheck
	// Копирование файла и параллельное вычисление контрольной суммы
	if size, sha512sum, err = file.New().CopyWithSha512Sum(ofh, ifh); err != nil {
		// Ошибка копирования
		err = fmt.Errorf("copy file to permanent storage %q error: %s", pathFull, err)
		return
	} else if ift.Size != uint64(size) {
		// Ошибка не совпадения размеров
		err = fmt.Errorf("write to file %q is not full size. Expected %d byte, writed %d byte", pathFull, ift.Size, size)
		return
	} else if ift.Sha512.MustValue() != sha512sum {
		// Ошибка контрольной суммы
		err = fmt.Errorf("SHA512 of file %q is wrong. Expected %q, calculated %q", pathFull, ift.Sha512.MustValue(), sha512sum)
		return
	}
	// Ошибка сброса буфера файла на диск
	if err = ofh.Sync(); err != nil {
		err = fmt.Errorf("sync(%q) error: %s", pathFull, err)
		return
	}
	// Сохранение информации в базу данных
	ret = &filestoreTypes.Filestore{
		Filename:    nul.NewStringValue(ift.Filename.MustValue()),
		FileExt:     nul.NewStringValue(ift.FileExt.MustValue()),
		Size:        uint64(size),
		Sha512:      nul.NewStringValue(sha512sum),
		LocalPath:   nul.NewStringValue(pathRelative),
		ContentType: nul.NewStringValue(ift.ContentType.MustValue()),
	}
	if err = ufm.Gist().
		Create(ret).
		Error; err != nil {
		err = fmt.Errorf("database model error: %s", err)
		return
	}

	return
}

// PermanentFileInfo Загрузка информации о файле по ID
func (ufm *impl) PermanentFileInfo(fileID uint64) (info *filestoreTypes.Filestore, err error) {
	info = new(filestoreTypes.Filestore)
	if ufm.Gist().
		Where("`deleteAt` IS NULL").
		Where("`id` = ?", fileID).
		First(info).
		RecordNotFound() {
		info, err = nil, ufm.Errors().ErrNotFound()
		return
	}
	return
}

// PermanentFileOpen Открытие для чтения постоянного файла по его ID
func (ufm *impl) PermanentFileOpen(fileID uint64) (info *filestoreTypes.Filestore, fh filestoreTypes.File, err error) {
	var pathFull string

	info = new(filestoreTypes.Filestore)
	if ufm.Gist().
		Where("`deleteAt` IS NULL").
		Where("`id` = ?", fileID).
		First(info).
		RecordNotFound() {
		info, err = nil, ufm.Errors().ErrNotFound()
		return
	}
	pathFull = path.Join(ufm.storagePath, info.LocalPath.MustValue())
	if fh, err = os.OpenFile(pathFull, os.O_RDONLY, os.FileMode(0640)); err != nil {
		info, fh, err = nil, nil, fmt.Errorf("open file %q error: %s", pathFull, err)
		return
	}

	return
}
