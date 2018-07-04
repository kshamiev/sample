package favicon // import "application/controllers/pages/favicon"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"regexp"

	"application/configuration"
)

const (
	_Png = `png`
	_Ico = `ico`
)

var (
	// Паттерн поиска картинки в файловой системе для /favicon.ico
	rexFaviconName = regexp.MustCompile(`(?i)^favicon\.(jpg|jpeg|png|gif|ico)$`)

	// Паттерн определения размера и типа для favicon.ico
	rexFaviconSizeAndType = regexp.MustCompile(`favicon-*(\d+)*x*(\d+)*\.(.*?)$`)
)

// Interface is an interface of package
type Interface interface {
	// FavIcon Cоздания на лету favicon.ico и всех возможных вариантов размеров и форматов favicon
	// Поиск в папке DocumentRoot файла favicon.* являющегося картинкой, конвертация в icon, ресайз при необходимости и вывод
	FavIcon(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of package
type impl struct {
	cfg          configuration.Interface // Конфигурация приложения
	DocumentRoot string                  // Путь к папке document root
}
