package uploadfile // import "application/controllers/apiV1/core/uploadfile"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"

	"application/models/filestore"
)

// Interface is an interface of package
type Interface interface {
	// UploadFile Загрузка файла на сервер
	UploadFile(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of package
type impl struct {
	ufm filestore.Interface // Интерфейс загрузки и хранения файлов
}

// Response Ответ. Загрузка файла на сервер
type Response struct {
	ID          uint64 `json:"id"`          // Уникальный идентификатор присвоенный файлу
	Field       string `json:"field"`       // Имя поля формы в котором был найден файл
	Filename    string `json:"filename"`    // Имя загруженного файла, которое было указано в форме
	Size        uint64 `json:"size"`        // Размер загруженного файла в байтах
	ContentType string `json:"contentType"` // MIME Content-Type содержимого файла
}
