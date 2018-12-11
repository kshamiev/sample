package pool // import "application/controllers/resource/pool"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"io"
	"sync"
)

// Interface is an interface of package
type Interface interface {
	// ByteBufferGet Извлечение из sync.Pool элемента для использования
	ByteBufferGet() ByteBufferInterface

	// ByteBufferPut Возврат в sync.Pool использованного элемента
	ByteBufferPut(item ByteBufferInterface)

	// Debug Set debug mode
	Debug(d bool) Interface
}

// impl is an implementation of package
type impl struct {
	debug       bool       // =true - debug mode is on
	ByteBuffers *sync.Pool // Pool of objects
}

// ByteBufferInterface Interface of pool objects
type ByteBufferInterface interface {
	io.Writer
	io.WriterTo
	io.Reader

	// Len Returns the number of bytes of data
	Len() int

	// Выгрузка содержимого в виде среза байт
	Bytes() []byte
}

// ByteBuffer Структура элемента
type ByteBuffer struct {
	Buf *bytes.Buffer
}
