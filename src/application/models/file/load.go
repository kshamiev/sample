package file // import "application/models/file"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// LoadFile Загрузка файла в память и возврат в виде *bytes.Buffer
func (fl *impl) LoadFile(name string) (data *bytes.Buffer, info os.FileInfo, err error) {
	var fh *os.File

	fh, err = os.OpenFile(name, os.O_RDONLY, os.FileMode(0755))
	if err != nil {
		return
	}
	defer fh.Close() // nolint: errcheck, gosec
	if info, err = fh.Stat(); err != nil {
		err = fmt.Errorf("Getting stat() of file %q ended with error: %s", name, err)
		return
	}
	data = &bytes.Buffer{}
	if _, err = io.Copy(data, fh); err != nil {
		err = fmt.Errorf("Reading data from the file %q ended with error: %s", name, err)
		return
	}

	return
}
