package sitemap // import "application/controllers/resource/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"

	modelsSitemap "application/models/sitemap"
)

// Links Получение текущих линков sitemap.xml или sitemap-index.xml
func (smi *impl) Links() []string { return smi.generatorSitemapURN() }

// Расчёт ссылок на sitemap.xml или sitemap-index.xml
// Максимально возможное количество страниц публикуемых sitemap: 18'446'744'073'709'551'615
func (smi *impl) generatorSitemapURN() (ret []string) {
	var count uint64
	var i, n int
	var j, l float64

	// Получение количества URN в базе данных
	if count = smi.Smm.Count(); count == 0 {
		return
	}
	// до 50'000 - один sitemap.xml
	if count <= modelsSitemap.MaxRecords {
		ret = append(ret, fmt.Sprintf("/sitemap.xml"))
		return
	}
	// больше 50'000 но меньше 2'500'000'000 - один sitemap-index.xml
	j = float64(count) / float64(modelsSitemap.MaxRecords)
	if j <= float64(modelsSitemap.MaxRecords) {
		ret = append(ret, fmt.Sprintf("/sitemap-index.xml"))
		return
	}
	// больше 2'500'000'000 - несколько sitemap-index.xml
	j = j / float64(modelsSitemap.MaxRecords)
	if l = j - float64(int(j)); l != 0.0 {
		i = int(j) + 1
	} else {
		i = int(j)
	}
	// но не более 368'934'881'474'191 индексов
	// или не более 18'446'744'073'709'551'615 URN адресов
	for n = 0; n < i && n < int(modelsSitemap.MaxIndex); n++ {
		ret = append(ret, fmt.Sprintf("/sitemap-index-%d.xml", n))
	}

	return
}
