package settings // import "application/controllers/apiV1/core/settings"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"time"
)

// Interface is an interface of controller
type Interface interface {
	// Time Check time difference and return server time in UTC
	// PUT /api/v1.0/settings/time
	Time(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of Controller
type impl struct {
}

// RequestTime Request API structure
// Верификация данных через библиотеку "gopkg.in/go-playground/validator.v9"
type RequestTime struct {
	Time time.Time `json:"time" validate:"required"` // Дата и время локального времени агента приведённая в UTC таймзону
}

// ResponseTime Response API structure
type ResponseTime struct {
	Time  time.Time     `json:"time"`  // Дата и время серверного времени приведённая в UTC таймзону
	Delta time.Duration `json:"delta"` // Посчитанная разница между локальным временем клиента и серверным временем в наносекундах. Отрицательное или положительное число. В результате сложения дельты и локального времени клиента в UTC должно получиться серверное время в UTC
}
