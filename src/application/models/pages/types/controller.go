package types // import "application/models/pages/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"gopkg.in/webnice/web.v1/method"
)

// Controller interface
type Controller interface {
	// Init Инициализация контроллера
	Init(url string) (rider *Manifest, err error)
}

// HandlerFunc Функция обработчик запроса
type HandlerFunc func(rr RequestResponse) (err error)

// Manifest Регистрационные данные для роутинга и потребности контроллера для выполнения работы
type Manifest struct {
	// Handlers Описание URN на которые контроллер подписывается
	Handlers []Handler

	// TemplatesURN Список URN, шаблоны которых требуются контроллеру для работы
	TemplatesURN []string
}

// Handler Структура подписки контроллера на URN
type Handler struct {
	// URN адрес ресурса
	URN string

	// Method Метод HTTP запроса ресурса
	Method method.Method

	// HandlerFunc Функция-обработчик запроса
	HandlerFunc HandlerFunc

	// Private Если указано true - не регистрируется для обработки запросов
	Private bool
}
