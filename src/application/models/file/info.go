package file // import "application/models/file"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

// GetFileName Выделение из полного пути и имени файла, имя файла
func (fl *impl) GetFileName(fileName string) (ret string) {
	var ch []string

	ret, ch = fileName, strings.Split(fileName, string(os.PathSeparator))
	if len(ch) > 0 {
		ret = ch[len(ch)-1]
	}

	return
}

// GetInfoSha512 Считывание информации о файле с контрольной суммой
func (fl *impl) GetInfoSha512(fn string) (inf *InfoSha512, err error) {
	var fh *os.File
	var s512 hash.Hash
	var fi os.FileInfo

	inf = new(InfoSha512)
	if fh, err = os.OpenFile(fn, os.O_RDONLY, os.FileMode(0755)); err != nil {
		err = fmt.Errorf("Opening file %q error: %s", fn, err)
		return
	}
	defer fh.Close() // nolint: errcheck, gosec

	s512 = sha512.New()
	inf.Size, err = io.Copy(s512, fh)
	if err != nil {
		err = fmt.Errorf("Reading file %q error: %s", fn, err)
		return
	}
	fi, err = fh.Stat()
	if err != nil {
		err = fmt.Errorf("Getting information about a file %q error: %s", fn, err)
		return
	}
	if inf.Size != fi.Size() {
		err = fmt.Errorf("SHA512 sum error, file size mismatch")
		return
	}
	inf.Sha512, inf.Name = hex.EncodeToString(s512.Sum(nil)[:]), fl.GetFileName(fn)

	return
}
