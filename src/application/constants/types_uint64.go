package constants // import "application/constants"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"strconv"
	"time"
)

// TypeConstantUint64 Константа типа Uint64
type TypeConstantUint64 uint64

// Second Возвращает значение как uint64 конвертируя наносекунды в секунды
func (ctu64 TypeConstantUint64) Second() uint64 {
	return uint64(time.Duration(ctu64) / time.Second)
}

// Duration Возвращает значение как time.Duration
func (ctu64 TypeConstantUint64) Duration() time.Duration {
	return time.Duration(ctu64)
}

// String Возвращает значение как строку
func (ctu64 TypeConstantUint64) String() string {
	return strconv.FormatUint(uint64(ctu64), 10)
}
