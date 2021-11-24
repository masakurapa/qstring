package qstringer

import (
	"fmt"
	"net/url"
	"reflect"
)

// Q is the type of the query string parameters.
type Q map[string]interface{}

// ArrayQ is a type of query string in array format.
type ArrayQ []interface{}

// MapQ is a type of query string in map format.
type MapQ map[string]interface{}

// Encode returns the URL-encoded query string.
//
// supports arguments of type map.
//
// add "?" to the beginning and return.
func Encode(v interface{}) (string, error) {
	e := encoder{v: url.Values{}}
	r := reflect.ValueOf(v)

	switch r.Kind() {
	case reflect.Map:
		iter := r.MapRange()
		for iter.Next() {
			// map key must be a string
			if iter.Key().Kind() != reflect.String {
				return "", fmt.Errorf("the key of the map type must be a string")
			}

			mapKey := iter.Key().String()
			m := iter.Value()
			if m.Kind() == reflect.Interface {
				if err := e.encode(mapKey, reflect.ValueOf((m.Interface()))); err != nil {
					return "", err
				}
				continue
			}
			if err := e.encode(mapKey, m); err != nil {
				return "", err
			}
		}
	default:
		return "", fmt.Errorf("type %s is not available", r.Kind().String())
	}

	if len(e.v) == 0 {
		return "", nil
	}
	return "?" + e.v.Encode(), nil
}

type encoder struct {
	v url.Values
}

func (e *encoder) encode(key string, r reflect.Value) (err error) {
	switch r.Kind() {
	case reflect.Bool:
		e.v.Add(key, fmt.Sprintf("%v", r.Bool()))
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		e.v.Add(key, fmt.Sprintf("%d", r.Int()))
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		e.v.Add(key, fmt.Sprintf("%d", r.Uint()))
	case reflect.Float32:
		e.v.Add(key, fmt.Sprintf("%v", float32(r.Float())))
	case reflect.Float64:
		e.v.Add(key, fmt.Sprintf("%v", r.Float()))
	case reflect.Map:
		err = e.encodeMap(key, r)
	case reflect.Array, reflect.Slice:
		err = e.encodeArray(key, r)
	case reflect.String:
		e.v.Add(key, r.String())
	default:
		err = fmt.Errorf("type %s is not available (key: %s)", r.Kind().String(), key)
	}
	return err
}

func (e *encoder) encodeMap(key string, r reflect.Value) error {
	iter := r.MapRange()
	for iter.Next() {
		// map key must be a string
		if iter.Key().Kind() != reflect.String {
			return fmt.Errorf("the key of the map type must be a string (key: %s)", key)
		}

		mapKey := key + "[" + iter.Key().String() + "]"
		m := iter.Value()
		if m.Kind() == reflect.Interface {
			if err := e.encode(mapKey, reflect.ValueOf((m.Interface()))); err != nil {
				return err
			}
			continue
		}
		if err := e.encode(mapKey, m); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeArray(key string, r reflect.Value) error {
	for i := 0; i < r.Len(); i++ {
		arrKey := fmt.Sprintf("%s[%d]", key, i)
		s := r.Index(i)
		if s.Kind() == reflect.Interface {
			if err := e.encode(arrKey, reflect.ValueOf((s.Interface()))); err != nil {
				return err
			}
			continue
		}
		if err := e.encode(arrKey, s); err != nil {
			return err
		}
	}
	return nil
}
