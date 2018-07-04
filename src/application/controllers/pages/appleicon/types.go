package appleicon // import "application/controllers/pages/appleicon"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"regexp"

	"application/configuration"
)

const (
	_Png = `png`
)

var (
	// Паттерн поиска картинки в файловой системе для /apple-touch-icon.png
	rexAppleIconName = regexp.MustCompile(`(?i)^favicon\.(jpg|jpeg|png|gif|ico)$`)

	// Паттерн определения размера для apple-touch-icon.png
	rexAppleIconSizeAndType = regexp.MustCompile(`^/apple-touch-*(icon|startup\-image)*-*(\d+)*x*(\d+)*(-*precomposed)*\.png$`)
)

// Interface is an interface of package
type Interface interface {
	// AppleIcon Specifying a Webpage Icon for Web Clip
	// Поиск в папке DocumentRoot файла favicon.* являющегося картинкой, конвертация в png, ресайз при необходимости и вывод
	AppleIcon(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of package
type impl struct {
	cfg          configuration.Interface // Конфигурация приложения
	DocumentRoot string                  // Путь к папке document root
}
