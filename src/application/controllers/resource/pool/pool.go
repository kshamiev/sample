package pool // import "application/controllers/resource/pool"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"io"
	"sync"
)

// New creates a new object and return interface
func New() Interface {
	var poo = new(impl)
	poo.ByteBuffers = &sync.Pool{
		New: poo.NewByteBuffersItem,
	}

	return poo
}

// Debug Set debug mode
func (poo *impl) Debug(d bool) Interface { poo.debug = d; return poo }

// NewByteBuffersItem Конструктор sync.Pool объекта для ByteBuffers
func (poo *impl) NewByteBuffersItem() interface{} {
	return &ByteBuffer{
		Buf: &bytes.Buffer{},
	}
}

// ByteBufferGet Извлечение из sync.Pool элемента для использования
func (poo *impl) ByteBufferGet() ByteBufferInterface {
	return poo.ByteBuffers.Get().(*ByteBuffer)
}

// ByteBufferPut Возврат в sync.Pool использованного элемента
func (poo *impl) ByteBufferPut(item ByteBufferInterface) {
	var elm = item.(*ByteBuffer)
	elm.ByteBufferClean()
	poo.ByteBuffers.Put(elm)
}

// ByteBufferClean Очистка элемента для переиспользования
func (bbi *ByteBuffer) ByteBufferClean() {
	bbi.Buf.Reset()
}

// Write Реализация интерфейса io.Writer
func (bbi *ByteBuffer) Write(p []byte) (n int, err error) { return bbi.Buf.Write(p) }

// WriteTo Реализация интерфейса io.WriteTo
func (bbi *ByteBuffer) WriteTo(w io.Writer) (n int64, err error) { return bbi.Buf.WriteTo(w) }

// Read Реализация интерфейса io.Reader
func (bbi *ByteBuffer) Read(p []byte) (n int, err error) { return bbi.Buf.Read(p) }

// Len Returns the number of bytes of data
func (bbi *ByteBuffer) Len() int { return bbi.Buf.Len() }

// Bytes returns a slice of length b.Len() holding the unread portion
func (bbi *ByteBuffer) Bytes() []byte { return bbi.Buf.Bytes() }
