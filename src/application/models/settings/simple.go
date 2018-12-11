package settings // import "application/models/settings"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	nul "gopkg.in/webnice/lin.v1/nl"
)

// StringSet Запись значения string
func (st *impl) StringSet(key string, value string) {
	st.LastError = st.Set(key, &settings{ValueString: nul.NewStringValue(value)})
}

// StringGet Чтение значения string
func (st *impl) StringGet(key string) (ret string) {
	var item = new(settings)

	if item, st.LastError = st.Get(key); st.LastError != nil {
		return
	}
	if !item.ValueString.Valid {
		st.LastError = st.Errors().ErrKeyOrValueNotFound()
		return
	}
	ret = item.ValueString.MustValue()

	return
}

// TimeSet Запись значения time.Time
func (st *impl) TimeSet(key string, value time.Time) {
	st.LastError = st.Set(key, &settings{ValueDate: nul.NewTimeValue(value)})
}

// TimeGet Чтение значения time.Time
func (st *impl) TimeGet(key string) (ret time.Time) {
	var item = new(settings)

	if item, st.LastError = st.Get(key); st.LastError != nil {
		return
	}
	if !item.ValueDate.Valid {
		st.LastError = st.Errors().ErrKeyOrValueNotFound()
		return
	}
	ret = item.ValueDate.MustValue()

	return
}

// UintSet Запись значения uint64
func (st *impl) UintSet(key string, value uint64) {
	st.LastError = st.Set(key, &settings{ValueUint: nul.NewUint64Value(value)})
}

// UintGet Чтение значения uint64
func (st *impl) UintGet(key string) (ret uint64) {
	var item = new(settings)

	if item, st.LastError = st.Get(key); st.LastError != nil {
		return
	}
	if !item.ValueUint.Valid {
		st.LastError = st.Errors().ErrKeyOrValueNotFound()
		return
	}
	ret = item.ValueUint.MustValue()

	return
}

// IntSet Запись значения int64
func (st *impl) IntSet(key string, value int64) {
	st.LastError = st.Set(key, &settings{ValueInt: nul.NewInt64Value(value)})
}

// IntGet Чтение значения int64
func (st *impl) IntGet(key string) (ret int64) {
	var item = new(settings)

	if item, st.LastError = st.Get(key); st.LastError != nil {
		return
	}
	if !item.ValueInt.Valid {
		st.LastError = st.Errors().ErrKeyOrValueNotFound()
		return
	}
	ret = item.ValueInt.MustValue()

	return
}

// DecimalSet Запись значения float64 как decimal
func (st *impl) DecimalSet(key string, value float64) {
	st.LastError = st.Set(key, &settings{ValueDecimal: nul.NewFloat64Value(value)})
}

// DecimalGet Чтение значения float64 как decimal
func (st *impl) DecimalGet(key string) (ret float64) {
	var item = new(settings)

	if item, st.LastError = st.Get(key); st.LastError != nil {
		return
	}
	if !item.ValueDecimal.Valid {
		st.LastError = st.Errors().ErrKeyOrValueNotFound()
		return
	}
	ret = item.ValueDecimal.MustValue()

	return
}

// FloatSet Запись значения float64 как double
func (st *impl) FloatSet(key string, value float64) {
	st.LastError = st.Set(key, &settings{ValueFloat: nul.NewFloat64Value(value)})
}

// FloatGet Чтение значения float64 как double
func (st *impl) FloatGet(key string) (ret float64) {
	var item = new(settings)

	if item, st.LastError = st.Get(key); st.LastError != nil {
		return
	}
	if !item.ValueFloat.Valid {
		st.LastError = st.Errors().ErrKeyOrValueNotFound()
		return
	}
	ret = item.ValueFloat.MustValue()

	return
}

// BooleanSet Запись значения bool
func (st *impl) BooleanSet(key string, value bool) {
	st.LastError = st.Set(key, &settings{ValueBit: nul.NewBoolValue(value)})
}

// BooleanGet Чтение значения bool
func (st *impl) BooleanGet(key string) (ret bool) {
	var item = new(settings)

	if item, st.LastError = st.Get(key); st.LastError != nil {
		return
	}
	if !item.ValueBit.Valid {
		st.LastError = st.Errors().ErrKeyOrValueNotFound()
		return
	}
	ret = item.ValueBit.MustValue()

	return
}

// BlobSet Запись значения []byte
func (st *impl) BlobSet(key string, value []byte) {
	st.LastError = st.Set(key, &settings{ValueBlob: nul.NewBytesValue(value)})
}

// BlobGet Чтение значения []byte
func (st *impl) BlobGet(key string) (ret []byte) {
	var item = new(settings)

	if item, st.LastError = st.Get(key); st.LastError != nil {
		return
	}
	if !item.ValueBlob.Valid {
		st.LastError = st.Errors().ErrKeyOrValueNotFound()
		return
	}
	ret = item.ValueBlob.MustValue()

	return
}
