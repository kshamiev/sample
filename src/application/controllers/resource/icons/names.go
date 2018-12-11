package icons // import "application/controllers/resource/icons"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
)

// Создание имени файла для иконок
func (ici *impl) makeImageName(inf *size, preffix string) (ret string) {
	const (
		tplNoSize = `/%s.%s`
		tplSize   = `/%s-%dx%d.%s`
	)
	if inf.Width == 0 || inf.Height == 0 {
		ret = fmt.Sprintf(tplNoSize, preffix, inf.Ext)
	} else {
		ret = fmt.Sprintf(tplSize, preffix, inf.Width, inf.Height, inf.Ext)
	}

	return
}

// Создание имени файла для иконок favicon
func (ici *impl) faviconImageName(inf *size) string {
	return ici.makeImageName(inf, keyFavicon)
}

// Создание имени файла для иконок apple-touch-icon
func (ici *impl) appleTouchIconImageName(inf *size) string {
	return ici.makeImageName(inf, keyAppleTouchIcon)
}

// Создание имени файла для иконок android-chrome
func (ici *impl) androidChromeImageName(inf *size) string {
	return ici.makeImageName(inf, keyAndroidChrome)
}

// Создание имени файла для иконок mstile
func (ici *impl) mstileImageName(inf *size) string {
	return ici.makeImageName(inf, keyMstile)
}
