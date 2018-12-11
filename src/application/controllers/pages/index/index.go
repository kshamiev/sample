package index // import "application/controllers/pages/index"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	modelsPages "application/models/pages"
	pagesTypes "application/models/pages/types"

	"gopkg.in/webnice/web.v1/method"
)

// Регистрация контроллера в модели работы с шаблонами страниц
func init() {
	modelsPages.RegisterController(New())
}

// New creates a new object and return interface
func New() Interface {
	var pix = new(impl)
	return pix
}

// Init Инициализация контроллера
func (pix *impl) Init(url string) (rider *pagesTypes.Manifest, err error) {
	var i int

	pix.serverURL = url
	rider = &pagesTypes.Manifest{
		Handlers: []pagesTypes.Handler{
			pagesTypes.Handler{
				URN:         `/landing`,
				Method:      method.Get,
				HandlerFunc: pix.Main,
			},
			pagesTypes.Handler{
				URN:         `/test`,
				Method:      method.Get,
				HandlerFunc: pix.Main,
			},
		},
		TemplatesURN: []string{`/`, `/landing`},
	}
	for i = range indexURN {
		rider.Handlers = append(rider.Handlers, pagesTypes.Handler{
			URN:         indexURN[i],
			Method:      method.Get,
			HandlerFunc: pix.Main,
			Private:     true,
		})
	}

	return
}

// Main Функция - обработчик запросов
func (pix *impl) Main(rr pagesTypes.RequestResponse) (err error) {

	return
}
