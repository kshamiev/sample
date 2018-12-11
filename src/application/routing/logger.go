package routing // import "application/routing"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"net"
	"net/http"
	"time"

	"gopkg.in/webnice/kit.v1/middleware/wrapsrw"
	"gopkg.in/webnice/log.v2/level"
	"gopkg.in/webnice/web.v1/header"
)

type logStruct struct {
	Address  string        // IP адрес клиента
	Code     int           // Код ответа
	Method   string        // Метод запроса
	Size     uint64        // Размер ответа сервера в байтах
	Path     string        // Запрашиваемый путь
	LeadTime time.Duration // Время выполнения запроса
}

// Logger Настройка логгера
func (rt *impl) Logger() Interface {
	rt.Rou.Use(rt.LoggerHandler())
	return rt
}

// IrisLoggerHandlerFunc Кастомный логер
func (rt *impl) LoggerHandler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		var fn = func(wr http.ResponseWriter, rq *http.Request) {
			var beginTime time.Time
			var ld logStruct
			var ll level.Level
			var wrp = wrapsrw.New(wr, rq.ProtoMajor)

			// Request time
			beginTime = time.Now()

			// The result is always written
			defer func() {
				// Production mode
				if !rt.debug && wrp.Status() < 500 {
					return
				}

				// Расчёт времени и вывод в лог
				ld.LeadTime = time.Since(beginTime)
				ld.Address = rt.GetRemoteAddress(rq)
				ld.Method = rq.Method
				ld.Path = rq.URL.Path
				ld.Code = wrp.Status()
				ld.Size = wrp.Len()

				// Projection of http code to error level
				switch {
				case ld.Code >= 500:
					ll = level.New().Error()
				case ld.Code >= 400:
					ll = level.New().Warning()
				default:
					ll = level.New().Informational()
				}

				// Disable normal request in production mode
				if !rt.debug && ll.Int() >= 200 && ll.Int() < 300 {
					return
				}

				log.Keys(
					log.Key{"Address": ld.Address},
					log.Key{"Code": ld.Code},
					log.Key{"Method": ld.Method},
					log.Key{"Size": ld.Size},
					log.Key{"Path": ld.Path},
					log.Key{"LeadTime": ld.LeadTime.String()},
					log.Key{"LeadTimeMillisecond": ld.LeadTime.Nanoseconds() / int64(time.Millisecond)},
					log.Key{"Location": wr.Header().Get(header.Location)},
				).Message(ll, "%-15s %3d %-7s %014d %-60s %s",
					ld.Address,
					ld.Code,
					ld.Method,
					ld.Size,
					ld.Path,
					ld.LeadTime)
			}()

			// Performing the next handle in the queue
			next.ServeHTTP(wrp, rq)
		}
		return http.HandlerFunc(fn)
	}
}

// GetRemoteAddress Получение IP адреса пользователя
func (rt *impl) GetRemoteAddress(rq *http.Request) (remoteAddress string) {
	var err error

	remoteAddress = rq.Header.Get(header.XRealIP)
	if remoteAddress != "" {
		return
	}
	remoteAddress = rq.Header.Get(header.XForwardedFor)
	if remoteAddress != "" {
		return
	}
	remoteAddress, _, err = net.SplitHostPort(rq.RemoteAddr)
	if err != nil {
		remoteAddress = rq.RemoteAddr
	}

	return
}
