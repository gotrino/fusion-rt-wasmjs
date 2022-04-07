package reflect

import (
	"fmt"
	"reflect"
)

func Fields(a any) []reflect.StructField {
	var res []reflect.StructField
	typ := reflect.TypeOf(a).Elem()
	for i := 0; i < typ.NumField(); i++ {
		res = append(res, typ.Field(i))
	}

	return res
}

func UnsupportedType(a any) {
	panic(fmt.Errorf("type %s is unsupported", reflect.TypeOf(a).String()))
}
