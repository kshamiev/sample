package pages // import "application/models/pages"

//import "gopkg.in/webnice/log.v2"
//import "gopkg.in/webnice/debug.v1"
import (
	"reflect"
)

// Получение уникального имени пакета
func packageName(obj interface{}) (ret string) {
	var rt = reflect.TypeOf(obj)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	ret = rt.PkgPath()

	return
}
