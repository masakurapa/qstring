package qstring

import (
	"reflect"
	"strings"
)

const (
	tagName      = "qstring"
	optSeparator = ","
	omitempty    = "omitempty"
)

type tagOption struct {
	omitempty bool
}

func parseTag(tag reflect.StructTag) (string, tagOption) {
	s := tag.Get(tagName)
	if idx := strings.Index(s, optSeparator); idx != -1 {
		return s[:idx], tagOption{
			omitempty: strings.Contains(s[idx:], omitempty),
		}
	}
	return s, tagOption{}
}

func isEmptyValue(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	}
	return false
}
