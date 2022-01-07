package qstring

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Q is the type of the query string parameters.
type Q map[string]interface{}

// ArrayQ is a type of query string in array format.
type ArrayQ []interface{}

// Encode returns the URL-encoded query string.
//
// Support struct, map type where the key is a string.
//
// By default, a nil value will be output.
//
// If you don't want to output nil values,
// Please specify option "omitempty" in the tag.
func Encode(v interface{}) (string, error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.IsValid() {
		return "", fmt.Errorf("nil is not available")
	}

	en := encoder{
		v: url.Values{},
	}

	switch rv.Kind() {
	case reflect.Map:
		if err := en.encodeMap("", rv); err != nil {
			return "", err
		}
	case reflect.Struct:
		if err := en.encodeStruct("", rv); err != nil {
			return "", err
		}
	case reflect.String:
		q := ""
		if strings.HasPrefix(rv.String(), que) {
			q = que
		}
		s := strings.TrimPrefix(rv.String(), que)
		return q + strings.ReplaceAll(url.QueryEscape(s), "%3D", "="), nil
	default:
		return "", fmt.Errorf("type %s is not available", rv.Kind().String())
	}

	if len(en.v) == 0 {
		return "", nil
	}
	return en.v.Encode(), nil
}
