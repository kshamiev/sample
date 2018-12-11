package sitemap

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

// Вероятная частота изменения страницы web сайта
// Это значение предоставляет общую информацию для поисковых систем и может не соответствовать точно частоте
// сканирования страницы web сайта
const (
	// CfNever Описания документов, которые не изменяются никогда
	CfNever ChangeFreq = iota

	// CfAlways Описания документов, которые изменяются при каждом доступе к этим документам
	CfAlways

	// CfHourly Описания документов, которые изменяются раз в час
	CfHourly

	// CfDaily Описания документов, которые изменяются раз в день
	CfDaily

	// CfWeekly Описания документов, которые изменяются раз в неделю
	CfWeekly

	// CfMonthly Описания документов, которые изменяются раз в месяц
	CfMonthly

	// CfYearly Описания документов, которые изменяются раз в год
	CfYearly
)

var cfString = map[ChangeFreq]string{
	CfNever:   `never`,
	CfAlways:  `always`,
	CfHourly:  `hourly`,
	CfDaily:   `daily`,
	CfWeekly:  `weekly`,
	CfMonthly: `monthly`,
	CfYearly:  `yearly`,
}

// ChangeFreq type of <changefreq>
type ChangeFreq uint8

// String convert ChangeFreq to string
func (cf ChangeFreq) String() string { return cfString[cf] }

// Uint8 convert ChangeFreq to uint8
func (cf ChangeFreq) Uint8() uint8 { return uint8(cf) }

// Uint64 convert ChangeFreq to uint64
func (cf ChangeFreq) Uint64() uint64 { return uint64(cf) }

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (cf ChangeFreq) MarshalBinary() (ret []byte, err error) {
	ret = make([]byte, 1)
	ret[0] = byte(cf)
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (cf *ChangeFreq) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 1 {
		err = errors.New("invalid size of change freq")
		return
	}
	switch ChangeFreq(data[0]) {
	case CfNever, CfAlways, CfHourly, CfDaily, CfWeekly, CfMonthly, CfYearly:
		*cf = ChangeFreq(data[0])
	default:
		err = errors.New("invalid change freq constant")
		return
	}

	return
}

// MarshalXML implements the xml.Marshaler interface
func (cf ChangeFreq) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	return e.EncodeElement(cf.String(), start)
}

// UnmarshalXML implements the xml.Unmarshaler interface
func (cf *ChangeFreq) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var value string

	if err = d.DecodeElement(&value, &start); err != nil {
		err = fmt.Errorf("unmarshal XML error: %s", err)
		return
	}
	switch strings.ToLower(value) {
	case CfNever.String():
		*cf = CfNever
	case CfAlways.String():
		*cf = CfAlways
	case CfHourly.String():
		*cf = CfHourly
	case CfDaily.String():
		*cf = CfDaily
	case CfWeekly.String():
		*cf = CfWeekly
	case CfMonthly.String():
		*cf = CfMonthly
	case CfYearly.String():
		*cf = CfYearly
	default:
		err = errors.New("invalid change freq constant")
		return
	}

	return
}
