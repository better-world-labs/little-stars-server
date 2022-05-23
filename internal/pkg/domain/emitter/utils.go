package emitter

import (
	"fmt"
	"reflect"
)

func GetStructType(i interface{}) string {
	_type := reflect.TypeOf(i)
	if isStructPointer(_type) {
		_type = _type.Elem()
	}

	return fmt.Sprintf("%s/%s", _type.PkgPath(), _type.Name())
}

func isStructPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}
