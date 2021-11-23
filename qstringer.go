package qstringer

import (
	"fmt"
	"net/url"
	"reflect"
)

// Q is the type of the query string parameters
type Q map[string]interface{}

// ArrayQ is a type of query string in array format
type ArrayQ []interface{}

// MapQ is a type of query string in map format
type MapQ map[string]interface{}

// Encode returns the URL-encoded query string
//
// add "?" to the beginning and return
func (q *Q) Encode() (string, error) {
	if len(*q) == 0 {
		return "", nil
	}

	v := values{
		v: url.Values{},
	}
	for key, value := range *q {
		r := reflect.ValueOf(value)
		if err := v.add(key, r); err != nil {
			return "", err
		}
	}

	return "?" + v.v.Encode(), nil
}

type values struct {
	v url.Values
}

func (v *values) add(key string, r reflect.Value) (err error) {
	switch r.Kind() {
	case reflect.Bool:
		v.v.Add(key, fmt.Sprintf("%v", r.Bool()))
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		v.v.Add(key, fmt.Sprintf("%d", r.Int()))
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		v.v.Add(key, fmt.Sprintf("%d", r.Uint()))
	case reflect.Float32:
		v.v.Add(key, fmt.Sprintf("%v", float32(r.Float())))
	case reflect.Float64:
		v.v.Add(key, fmt.Sprintf("%v", r.Float()))
	case reflect.Map:
		err = v.addMap(key, r)
	case reflect.Array, reflect.Slice:
		v.addArray(key, r)
	case reflect.String:
		v.v.Add(key, r.String())
	default:
		err = fmt.Errorf("type %s is not available (key: %s)", r.Kind().String(), key)
	}
	return err
}

func (v *values) addMap(key string, r reflect.Value) error {
	iter := r.MapRange()
	for iter.Next() {
		// map key must be a string
		if iter.Key().Kind() != reflect.String {
			return fmt.Errorf("the key of the map type must be a string (key: %s)", key)
		}

		mapKey := key + "[" + iter.Key().String() + "]"
		m := iter.Value()
		if m.Kind() == reflect.Interface {
			v.add(mapKey, reflect.ValueOf((m.Interface())))
			continue
		}
		v.add(mapKey, m)
	}
	return nil
}

func (v *values) addArray(key string, r reflect.Value) {
	arrKey := key + "[]"
	for i := 0; i < r.Len(); i++ {
		s := r.Index(i)
		if s.Kind() == reflect.Interface {
			v.add(arrKey, reflect.ValueOf((s.Interface())))
			continue
		}
		v.add(arrKey, s)
	}
}
