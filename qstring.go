package qstring

import (
	"errors"
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

// Decode returns the URL-encoded query string
//
// add "?" to the beginning and return
func Decode(s string, v interface{}) error {
	if v == nil {
		return errors.New("nil is not available")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("not pointer")
	}

	riv := reflect.Indirect(rv)
	if !riv.IsValid() {
		return errors.New("nil is not available")
	}
	d := decoder{query: s, rv: riv}

	switch riv.Kind() {
	case reflect.String:
		return d.decodeString()
	case reflect.Array:
		return d.decodeArray()
	case reflect.Slice:
		return d.decodeSlice()
	case reflect.Map:
		return d.decodeMap()
	case reflect.Struct:
		return d.decodeStruct()
	}

	return errors.New("type " + riv.Kind().String() + " is not available")
}

func DecodeToMap(s string) (Q, error) {
	var q Q
	err := Decode(s, &q)
	if err != nil {
		return nil, err
	}
	return q, nil
}
