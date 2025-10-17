package openapi

import (
	"fmt"
	"reflect"
)

func typeName(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Name()
}

func componentRef(category, elem string) string {
	return fmt.Sprintf("#/components/%s/%s", category, elem)
}
