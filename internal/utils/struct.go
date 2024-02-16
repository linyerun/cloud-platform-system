package utils

import (
	"reflect"
)

func IsNotNilPointer(value any) bool {
	valueOf := reflect.ValueOf(value)
	return valueOf.Kind() == reflect.Pointer && !valueOf.IsNil()
}
