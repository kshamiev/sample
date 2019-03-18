package filestore // import "application/models/filestore"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"application/models/file"
	filestoreTypes "application/models/filestore/types"

	nul "gopkg.in/webnice/lin.v1/nl"
)

// NewTemporaryFile Создание нового временного файла
func (ufm *impl) NewTemporaryFile(filename string, size uint64, contentType string, inpFh io.Reader) (id uint64, err error) {
	var ext string
	var pathFull, pathRelative string
	var fh *os.File
	var l int64
	var sha512sum string
	var ft *filestoreTypes.FilesTemporary

	// Извлечение расширения из имени файла
	if tmp := rexFileExt.FindStringSubmatch(filename); len(tmp) == 2 {
		ext = strings.ToLower(tmp[1])
	}
	// Создание имени файла во временном хранилище
	pathFull, pathRelative = ufm.makeTemporaryFilename()
	// Создание директории к файлу
	if err = os.MkdirAll(path.Dir(pathFull), os.FileMode(0750)); err != nil {
		err = fmt.Errorf("make directory error: %s", err)
		return
	}
	// Создание файла
	if fh, err = os.OpenFile(pathFull, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0640)); err != nil {
		err = fmt.Errorf("create file %q error: %s", pathFull, err)
		return
	}
	defer fh.Close() // nolint: errcheck
	// Запись файла и параллельное вычисление контрольной суммы
	if l, sha512sum, err = file.New().CopyWithSha512Sum(fh, inpFh); err != nil {
		err = fmt.Errorf("copy data to file %q error: %s", pathFull, err)
		return
	} else if uint64(l) != size {
		err = fmt.Errorf("write to file %q is not full size. Expected %d byte, Writed %d byte", pathFull, size, uint64(l))
		return
	}
	if err = fh.Sync(); err != nil {
		err = fmt.Errorf("sync(%q) error: %s", pathFull, err)
		return
	}
	// Сохранение информации в базу данных
	ft = &filestoreTypes.FilesTemporary{
		Filename:    nul.NewStringValue(filename),
		FileExt:     nul.NewStringValue(ext),
		Size:        size,
		Sha512:      nul.NewStringValue(sha512sum),
		LocalPath:   nul.NewStringValue(pathRelative),
		ContentType: nul.NewStringValue(contentType),
	}
	if err = ufm.Gist().
		Create(ft).
		Error; err != nil {
		err = fmt.Errorf("database error: %s", err)
		return
	}
	id = ft.ID

	return
}

// TemporaryFileOpen Открытие для чтения временного файла по его ID
func (ufm *impl) TemporaryFileOpen(fileID uint64) (info *filestoreTypes.FilesTemporary, fh filestoreTypes.File, err error) {
	var pathFull string

	info = new(filestoreTypes.FilesTemporary)
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
