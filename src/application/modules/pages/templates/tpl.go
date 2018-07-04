package templates // include "application/modules/pages/templates"

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"bytes"
)

// Len Количество шаблонов содержащихся в объекте
func (tp *tpls) Len() int { return len(tp.Tpl) }

// ListUrn Список URN всех шаблонов
func (tp *tpls) ListUrn() (ret []string) {
	for i := range tp.Tpl {
		ret = append(ret, tp.Tpl[i].UrnAddress)
	}
	return
}

// FirstUrn Первый URN в списке, как правило это URN к которому пришел запрос
func (tp *tpls) FirstUrn() (ret string) {
	for i := range tp.Tpl {
		ret = tp.Tpl[i].UrnAddress
		break
	}
	return
}

// HasUrn (urn) True если существует шаблон с указанным URN
func (tp *tpls) HasUrn(urn string) (ret bool) {
	for i := range tp.Tpl {
		if tp.Tpl[i].UrnAddress == urn {
			ret = true
		}
	}
	return
}

// Map Данные шаблона указанного URN, представленные в виде map.
// Ключём является имя файла без расширения. Например для файла layout.tpl.html ключ будет layout
func (tp *tpls) Map(urn string) (ret map[string]*TplBody) {
	for i := range tp.Tpl {
		if tp.Tpl[i].UrnAddress == urn {
			ret = tp.Tpl[i].MapData
		}
	}
	return
}

// Возвращает структуру тела шаблона с указанным URN и ключём, если шаблона нет или он пустой то вернётся nil
func (tp *tpls) keyBody(urn string, keyName string, orFirst bool) (ret *TplBody) {
	var key string
	for i := range tp.Tpl {
		if tp.Tpl[i].UrnAddress != urn {
			continue
		}
		for key = range tp.Tpl[i].MapData {
			if key == keyName || len(tp.Tpl[i].MapData) == 1 && orFirst {
				ret = tp.Tpl[i].MapData[key]
			}
		}
	}
	return
}

// Возвращает данные шаблона с указанным URN и ключём, если шаблона нет или он пустой то вернётся nil
func (tp *tpls) keyData(urn string, keyName string, orFirst bool) (ret *bytes.Buffer) {
	var key string
	for i := range tp.Tpl {
		if tp.Tpl[i].UrnAddress != urn {
			continue
		}
		for key = range tp.Tpl[i].MapData {
			if key == keyName || len(tp.Tpl[i].MapData) == 1 && orFirst {
				ret = tp.Tpl[i].MapData[key].Data
			}
		}
	}
	return
}

// Index Возвращает главный шаблон для указанного URN
// если шаблон является многофайловым (папка), то ищется шаблон с именем index,
// если шаблон является однофайловым, то возвращается единственный существующий шаблон
// если шаблона нет, то возвращается nil
func (tp *tpls) Index(urn string) *TplBody { return tp.keyBody(urn, _KeyIndex, true) }

// Layout Возвращает шаблон-макет (layout) для указанного URN
// если шаблона нет, то возвращается nil
func (tp *tpls) Layout(urn string) *TplBody { return tp.keyBody(urn, _KeyLayout, false) }

// KeyBody Возвращает структуру тела шаблона с указанным URN и ключём, если шаблона нет или он пустой то вернётся nil
func (tp *tpls) KeyBody(urn string, key string) *TplBody { return tp.keyBody(urn, key, false) }

// KeyData Возвращает данные шаблона с указанным URN и ключём, если шаблона нет или он пустой то вернётся nil
func (tp *tpls) KeyData(urn string, key string) *bytes.Buffer { return tp.keyData(urn, key, false) }

// Keys Получение списка всех ключей шаблонов для указанного URN
func (tp *tpls) Keys(urn string) (ret []string) {
	for i := range tp.Tpl {
		if tp.Tpl[i].UrnAddress != urn {
			continue
		}
		for key := range tp.Tpl[i].MapData {
			ret = append(ret, key)
		}
	}
	return
}

// HasKey (urn, key) True если для указанного URN существует шаблон с указанным ключём
func (tp *tpls) HasKey(urn string, keyName string) (ret bool) {
	for i := range tp.Tpl {
		if tp.Tpl[i].UrnAddress != urn {
			continue
		}
		for key := range tp.Tpl[i].MapData {
			if key == keyName {
				ret = true
			}
		}
	}
	return
}

// AllData Все данные шаблонов отсортированные в порядке: [все inc, layout, index] для указанного URN
// Шаблоны с пустыми данными пропускаются
func (tp *tpls) AllData(urn string) (ret []*bytes.Buffer) {
	var keys []string
	var key string

	keys = tp.Keys(urn)
	for _, key = range keys {
		if key == _KeyIndex || key == _KeyLayout {
			continue
		}
		if buf := tp.keyData(urn, key, false); buf != nil {
			ret = append(ret, buf)
		}
	}
	if buf := tp.keyData(urn, _KeyLayout, false); buf != nil {
		ret = append(ret, buf)
	}
	if buf := tp.keyData(urn, _KeyIndex, false); buf != nil {
		ret = append(ret, buf)
	}
	return
}

// Reload Перечитываем все шаблоны объекта и возвращаем интерфейс
// шаблоны перечитываются в соответствии с общими правилами, то есть только в режиме дебага либо если они не загружены в память
func (tp *tpls) Reload() (ret Template) {
	ret = tp
	for i := range tp.Tpl {
		for key := range tp.Tpl[i].MapData {
			if err := singleton.LoadFileBodyData(tp.Tpl[i].MapData[key]); err != nil {
				log.Errorf("Error LoadFileBodyData: %s", err.Error())
				return
			}
		}

	}
	return
}
