package qstringer

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

const (
	defaultNilValue = ""
	que             = "?"
)

// Q is the type of the query string parameters.
type Q map[string]interface{}

// ArrayQ is a type of query string in array format.
type ArrayQ []interface{}

// Encode returns the URL-encoded query string.
//
// Support struct, map type where the key is a string.
//
// Add "?" to the beginning and return.
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

type encoder struct {
	v url.Values
}

func (e *encoder) encode(key string, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Bool:
		e.v.Add(key, fmt.Sprintf("%v", rv.Bool()))
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		e.v.Add(key, fmt.Sprintf("%d", rv.Int()))
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		e.v.Add(key, fmt.Sprintf("%d", rv.Uint()))
	case reflect.Float32:
		e.v.Add(key, fmt.Sprintf("%v", float32(rv.Float())))
	case reflect.Float64:
		e.v.Add(key, fmt.Sprintf("%v", rv.Float()))
	case reflect.Map:
		return e.encodeMap(key, rv)
	case reflect.Array:
		return e.encodeArray(key, rv)
	case reflect.Slice:
		return e.encodeSlice(key, rv)
	case reflect.Struct:
		return e.encodeStruct(key, rv)
	case reflect.String:
		e.v.Add(key, rv.String())
	case reflect.Interface:
		if rv.IsNil() {
			e.v.Add(key, defaultNilValue)
			return nil
		}
		return e.encode(key, reflect.ValueOf(rv.Interface()))
	case reflect.Ptr:
		if rv.IsNil() {
			e.v.Add(key, defaultNilValue)
			return nil
		}
		return e.encode(key, reflect.Indirect(rv))
	default:
		return fmt.Errorf("type %s is not available (key: %s)", rv.Kind().String(), key)
	}

	return nil
}

func (e *encoder) encodeMap(key string, rv reflect.Value) error {
	if rv.IsNil() {
		e.v.Add(key, defaultNilValue)
		return nil
	}

	iter := rv.MapRange()
	for iter.Next() {
		// map key must be a string
		if iter.Key().Kind() != reflect.String {
			if key == "" {
				return fmt.Errorf("the key of the map type must be a string")
			} else {
				return fmt.Errorf("the key of the map type must be a string (key: %s)", key)
			}
		}

		if err := e.encode(e.makeMapKey(key, iter.Key().String()), iter.Value()); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeArray(key string, rv reflect.Value) error {
	for i := 0; i < rv.Len(); i++ {
		if err := e.encode(fmt.Sprintf("%s[%d]", key, i), rv.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeSlice(key string, rv reflect.Value) error {
	if rv.IsNil() {
		e.v.Add(key, defaultNilValue)
		return nil
	}
	return e.encodeArray(key, rv)
}

func (e *encoder) encodeStruct(key string, rv reflect.Value) error {
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		if !f.IsExported() {
			continue
		}

		tag, opt := parseTag(f.Tag)
		if tag == "" {
			continue
		}

		frv := rv.FieldByName(f.Name)
		if opt.omitempty && isEmptyValue(frv) {
			continue
		}

		if err := e.encode(e.makeMapKey(key, tag), frv); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) makeMapKey(key, ch string) string {
	if key == "" {
		return ch
	}
	return fmt.Sprintf("%s[%s]", key, ch)
}
