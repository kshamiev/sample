package icons // import "application/controllers/resource/icons"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"net/http"
	"path"
	"sync"

	"application/controllers/resource/pool"
	"application/models/filecache"
	"application/models/fileinfo"
	modulesImg "application/modules/img"

	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/route"
)

// New creates a new object and return interface
func New() Interface {
	var ici = &impl{
		Mfc:              filecache.Get(),
		Pool:             pool.New(),
		Img:              modulesImg.New(),
		ImgSrcRasterSync: new(sync.RWMutex),
	}
	return ici
}

// Debug Set debug mode
func (ici *impl) Debug(d bool) Interface {
	ici.debug = d
	ici.Pool.Debug(ici.debug)
	ici.Mfc.Debug(ici.debug)
	return ici
}

// DocumentRoot Устанавливает путь к корню веб сервера
func (ici *impl) DocumentRoot(path string) Interface { ici.rootPath = path; return ici }

// SetRouting Установка роутинга к статическим файлам
func (ici *impl) SetRouting(rou route.Interface) Interface {
	var err error
	var iName string
	var newHFn func(gi int) http.HandlerFunc
	var newVfn func(req string, gi int) filecache.CreateFn
	var i, j int

	// Preview page
	rou.Get(keyIconsPreviewPage, ici.iconsPreviewPageHandlerFunc)
	// favicon
	ici.groups = append(ici.groups, &group{
		name: keyFavicon, items: sizeFavicon, nameFn: ici.faviconImageName,
		min: ici.getMin(sizeFavicon), max: ici.getMax(sizeFavicon),
	})
	// apple-touch-icon
	ici.groups = append(ici.groups, &group{
		name: keyAppleTouchIcon, items: sizeAppleTouchIcon, nameFn: ici.appleTouchIconImageName,
		min: ici.getMin(sizeAppleTouchIcon), max: ici.getMax(sizeAppleTouchIcon),
	})
	// android-chrome
	ici.groups = append(ici.groups, &group{
		name: keyAndroidChrome, items: sizeAndroidChrome, nameFn: ici.androidChromeImageName,
		min: ici.getMin(sizeAndroidChrome), max: ici.getMax(sizeAndroidChrome),
	})
	// mstile
	ici.groups = append(ici.groups, &group{
		name: keyMstile, items: sizeMstile, nameFn: ici.mstileImageName,
		min: ici.getMin(sizeMstile), max: ici.getMax(sizeMstile),
	})
	// Дополнительные файлы
	// apple-touch-icon: GET /safari-pinned-tab.svg
	rou.Get(keySafariPinnedTabVector, ici.appleSafariPinnedTabVectorHandlerFunc)
	// android-chrome: GET /manifest.json
	rou.Get(keyManifest, ici.iconsManifestHandlerFunc)
	// mstile: GET /browserconfig.xml
	rou.Get(keyBrowserconfig, ici.iconsBrowserconfigHandlerFunc)
	// Создание роутинга на все группы иконок
	// используя замыкание с передачей группы иконок в контроллер http.HandlerFunc при каждом запросе
	newHFn = func(gi int) http.HandlerFunc {
		return func(wr http.ResponseWriter, rq *http.Request) { ici.iconHandlerFunc(wr, rq, ici.groups[gi]) }
	}
	newVfn = func(req string, gi int) filecache.CreateFn {
		return func(name string) (*filecache.MemoryObject, error) {
			return ici.makeRasterVirtual(req, ici.groups[gi])
		}
	}
	for i = range ici.groups {
		for j = range ici.groups[i].items {
			iName = ici.groups[i].nameFn(&ici.groups[i].items[j])
			if err = ici.Mfc.Virtual(iName, 0, newVfn(iName, i)); err != nil {
				log.Warningf("Register virtual file %q in filecache, error: %s", iName, err)
				continue
			}
			rou.Get(iName, newHFn(i))
		}
	}

	return ici
}

// Получение минимального не нулевого значения размера иконки группы
func (ici *impl) getMin(items []size) (min uint) {
	for i := range items {
		if min == 0 && items[i].Width > 0 || min == 0 && items[i].Height > 0 {
			if min = items[i].Width; min == 0 {
				min = items[i].Height
			}
		}
		if min > items[i].Width && items[i].Width > 0 {
			min = items[i].Width
		}
		if min > items[i].Height && items[i].Height > 0 {
			min = items[i].Height
		}
	}
	return
}

// Получение максимального значения размера иконки группы
func (ici *impl) getMax(items []size) (max uint) {
	for i := range items {
		if max < items[i].Width {
			max = items[i].Width
		}
		if max < items[i].Height {
			max = items[i].Height
		}
	}
	return
}

// Загрузка исходного растрового шаблона из файла в объект графического изображения
func (ici *impl) loadImgRaster(name string) (err error) {
	var mdi filecache.Data

	ici.ImgSrcRasterSync.Lock()
	defer ici.ImgSrcRasterSync.Unlock()

	// Загрузка исходного файла
	if mdi, err = ici.Mfc.Load(name); err != nil {
		err = fmt.Errorf("source file error: %s", err)
		return
	}
	// Загрузка графического объекта из файла
	ici.ImgSrcRaster = ici.Img.New()
	if _, err = ici.ImgSrcRaster.ReadFrom(mdi.Reader()); err != nil {
		ici.ImgSrcRaster, err = nil, fmt.Errorf("loading image file %q error: %s", name, err)
		return
	}
	ici.ImgSrcRaster.
		SetFileInfo(mdi).
		SetFilename(name)

	return
}

// Функция создания графического изображения через конвертацию или изменение размеров исходного графического изображения
func (ici *impl) makeRasterVirtual(req string, grp *group) (ret *filecache.MemoryObject, err error) {
	var name, srcName string
	var inf *size
	var im modulesImg.Image
	var content pool.ByteBufferInterface
	var ffi fileinfo.Interface

	log.Infof(" - create virtual file: %q", req)
	srcName, name, inf = path.Join(ici.rootPath, imgRaster), path.Join(ici.rootPath, req), ici.parseRequest(req)
	// Загрузка исходного растрового шаблона из файла в объект графического изображения
	if ici.ImgSrcRaster == nil {
		if err = ici.loadImgRaster(srcName); err != nil {
			return
		}
	}
	ici.ImgSrcRasterSync.RLock()
	defer ici.ImgSrcRasterSync.RUnlock()

	// Изменение размеров исходного растрового изображения
	if inf.Width != uint(ici.ImgSrcRaster.Config().Width) || inf.Height != uint(ici.ImgSrcRaster.Config().Height) {
		if inf.Width < grp.min || inf.Height < grp.min {
			im = ici.Img.Resize(ici.ImgSrcRaster, grp.min, grp.min)
		} else if inf.Width > grp.max || inf.Height > grp.max {
			im = ici.Img.Resize(ici.ImgSrcRaster, grp.max, grp.max)
		} else {
			im = ici.Img.Resize(ici.ImgSrcRaster, inf.Width, inf.Height)
		}
	} else {
		im = ici.ImgSrcRaster
	}
	im.SetFilename(name)
	// Получение и возврат интерфейса ByteBuffer
	content = ici.Pool.ByteBufferGet()
	defer ici.Pool.ByteBufferPut(content)
	// Выгрузка нового графического изобращения в виртуальный файл
	ici.Pool.ByteBufferGet()
	switch inf.ContentType {
	case mime.ImageICO:
		_, err = im.SetType(modulesImg.TypeICO).WriteTo(content)
	case mime.ImagePNG:
		_, err = im.SetType(modulesImg.TypePNG).WriteTo(content)
	default:
		err = fmt.Errorf("Content-Type %q is incorrect", inf.ContentType)
	}
	if err != nil {
		return
	}
	ffi = fileinfo.New().
		CopyFrom(im.FileInfo()).
		SetName(req).
		SetSize(int64(content.Len()))
	ret = &filecache.MemoryObject{
		ContentType: inf.ContentType,
		Info:        ffi,
		Body:        content.Bytes(),
	}

	return
}
