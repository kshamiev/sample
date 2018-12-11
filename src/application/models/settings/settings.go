package settings // import "application/models/settings"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

// New creates new implementation
func New() Interface {
	var st = new(impl)
	return st
}

// Errors Ошибки известного состояни, которые могут вернуть функции пакета
func (st *impl) Errors() *Error { return Errors() }

// Error Ошибка возникшая в результате последней операции
func (st *impl) Error() error { return st.LastError }

// Set Запись значения
func (st *impl) Set(key string, value *settings) (err error) {
	var tmp *settings

	if key == "" {
		err = st.Errors().ErrKeyIsNotUnique()
		return
	}
	tmp, value.Key = new(settings), key
	if !st.Gist().
		Where("`key` = ?", key).
		First(tmp).
		RecordNotFound() {
		value.ID, value.CreateAt, value.AccessAt = tmp.ID, tmp.CreateAt, tmp.AccessAt
	}
	err = st.Gist().
		Save(value).
		Error

	return
}

// Get Чтение значения
func (st *impl) Get(key string) (value *settings, err error) {
	if key == "" {
		err = st.Errors().ErrKeyIsNotUnique()
		return
	}
	value = new(settings)
	if st.Gist().
		Where("`key` = ?", key).
		First(value).
		RecordNotFound() {
		err = st.Errors().ErrKeyOrValueNotFound()
	}

	return
}
