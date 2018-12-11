package file // import "application/models/file"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// Copy Копирует один файл в другой
func (fl *impl) Copy(dst, src string) (size int64, err error) {
	var fhIn, fhOu *os.File

	fhIn, err = os.Open(src) // nolint: gosec
	if err != nil {
		return
	}
	defer fhIn.Close() // nolint: errcheck, gosec

	fhOu, err = os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644) // nolint: gosec
	if err != nil {
		return
	}
	defer func() {
		if lerr := fhOu.Sync(); err == nil {
			err = lerr
		}
		if lerr := fhOu.Close(); err == nil {
			err = lerr
		}
	}()
	size, err = io.Copy(fhOu, fhIn)

	return
}

// CopyWithSha512Sum Копирование контента с параллельным вычислением контрольной суммы алгоритмом SHA512
func (fl *impl) CopyWithSha512Sum(dst io.Writer, src io.Reader) (written int64, sha512sum string, err error) {
	var er, ew error
	var size, nr, nw int
	var buf []byte
	var l *io.LimitedReader
	var ok bool
	var sha hash.Hash

	size = 32 * 1024
	if l, ok = src.(*io.LimitedReader); ok && int64(size) > l.N {
		if l.N < 1 {
			size = 1
		} else {
			size = int(l.N)
		}
	}
	buf = make([]byte, size)
	sha = sha512.New()
	for {
		nr, er = src.Read(buf)
		if nr > 0 {
			_, _ = sha.Write(buf[0:nr]) // nolint: errcheck, gosec
			nw, ew = dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	sha512sum = hex.EncodeToString(sha.Sum(nil)[:])

	return
}
