package icons // import "application/controllers/resource/icons"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"regexp"
	"sync"

	"application/controllers/resource/pool"
	"application/models/filecache"
	modulesImg "application/modules/img"

	"gopkg.in/webnice/web.v1/route"
)

const (
	keySource                 = `icons`                         // Название шаблона
	keyFavicon                = `favicon`                       // Название иконки favicon
	keyAppleTouchIcon         = `apple-touch-icon`              // Название иконки apple-touch-icon
	keyAndroidChrome          = `android-chrome`                // Название иконки android-chrome
	keyMstile                 = `mstile`                        // Название иконки mstile
	keyIconsPreviewPage       = `/icons_preview_page`           // Страница просмотра всех иконок
	keySafariPinnedTabVector  = `/safari-pinned-tab.svg`        // Название safari-pinned-tab SVG картинки
	keyManifest               = `/manifest.json`                // Название файла манифеста
	keyBrowserconfig          = `/browserconfig.xml`            // Название browserconfig
	ifModifiedSinceTimeFormat = `Mon, 02 Jan 2006 15:04:05 GMT` // Формат даты и времени для заголовка IfModifiedSince
	htmlPreviewTemplate       = keySource + `.tpl.html`         // Шаблоны текстовых документов
	imgRaster                 = keySource + `.tpl.png`          // Исходный растровый шаблон
	imgVector                 = keySource + `.tpl.svg`          // Исходный векторный шаблон
)

var (
	// Паттерн разбора запроса иконки
	rexIconSizeAndType = regexp.MustCompile(`(` +
		keyFavicon +
		`|` + keyAppleTouchIcon +
		`|` + keyAndroidChrome +
		`|` + keyMstile +
		`)-*(\d+)*x*(\d+)*\.(.*?)$`)
	// Паттерн отделения тела JSON от самых первых скобок {}
	rexJSONInside = regexp.MustCompile(`(?ms)^(\s*{\s*)(.+)(\s*}\s*)$`)
)

// Interface is an interface of package
type Interface interface {
	// Debug Set debug mode
	Debug(d bool) Interface

	// DocumentRoot Устанавливает путь к корню веб сервера
	DocumentRoot(path string) Interface

	// Установка роутинга к статическим файлам
	SetRouting(rou route.Interface) Interface

	// ERRORS

	// ErrNotFound Not found
	ErrNotFound() error
}

// impl is an implementation of package
type impl struct {
	debug            bool                 // =true - debug mode is on
	rootPath         string               // Путь к корню веб сервера
	groups           []*group             // Группы иконок
	Mfc              filecache.Interface  // Интерфейс кеша в памяти
	Pool             pool.Interface       // Интерфейс пула переменных для переиспользования памяти
	Img              modulesImg.Interface // Интерфейс работы с графическими изображениями
	ImgSrcRaster     modulesImg.Image     // Загруженный в виде графического изображения исходный растровый шаблон
	ImgSrcRasterSync *sync.RWMutex        // Защита от race
}

type size struct {
	Width       uint
	Height      uint
	Ext         string
	ContentType string
	Wide        bool
}

type group struct {
	name   string             // Название группы, совпадает с константой keyFavicon, keyAppleTouchIcon, keyAndroidChrome, keyMstile
	items  []size             // Размеры иконок группы
	min    uint               // Минимальный размер иконки группы
	max    uint               // Максимальный размер иконки группы
	nameFn func(*size) string // Функция создания имени иконки
}

type htmlPreviewTemplateVars struct {
	Groups []*htmlPreviewTemplateVarsGroup
}

type htmlPreviewTemplateVarsGroup struct {
	Group string
	Items []size
}

// Переменные шаблонизатора
type manifestTemplateVars struct {
	Icons string
}

// Структура иконок для публикации в манифесте
type manifestIcons struct {
	Icons []*manifestIconsItems `json:"icons"`
}

// Структура иконок для публикации в манифесте
type manifestIconsItems struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes"`
	Type  string `json:"type"`
}

// Переменные шаблонизатора
type browserconfigTemplateVars struct {
	Icons []*browserconfigIconItem
}

// Структура иконок для публикации в browserconfig
type browserconfigIconItem struct {
	Width       uint
	Height      uint
	Ext         string
	ContentType string
	Wide        bool
}
