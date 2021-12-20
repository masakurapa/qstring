package qstringer

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Decode returns the URL-encoded query string
//
// add "?" to the beginning and return
func Decode(s string, v interface{}) error {
	if v == nil {
		return fmt.Errorf("nil is not available")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("not pointer")
	}

	riv := reflect.Indirect(rv)
	if !riv.IsValid() {
		return fmt.Errorf("nil is not available")
	}
	d := decoder{query: s, rv: riv}

	switch riv.Kind() {
	case reflect.String:
		return d.decodeString()
	case reflect.Array:
		return d.decodeArray()
	case reflect.Slice:
		return d.decodeSlice()
	}

	return fmt.Errorf("type %s is not available", riv.Kind().String())
}

type decoder struct {
	query string
	rv    reflect.Value
}

func (d *decoder) decodeMap(values url.Values) error {
	for key, vs := range values {
		// not slice or map
		if !strings.Contains(key, "[") {
			if len(vs) == 0 {
				// it shouldn't pass, but add a skip process just in case.
				continue
			}
			if len(vs) == 1 {
				// use only first value
				d.rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(vs[0]))
				continue
			}
			d.rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(vs))
			continue
		}

		d.rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(vs))
	}
	return nil
}

func (d *decoder) decodeString() error {
	q, err := url.QueryUnescape(d.query)
	if err == nil {
		d.rv.SetString(q)
	}
	return err
}

func (d *decoder) decodeArray() error {
	if d.rv.Type().Elem().Kind() != reflect.String {
		return fmt.Errorf("allocation type must be [n]stirng")
	}

	values, err := url.ParseQuery(d.query)
	if err != nil {
		return err
	}

	arrVals, ok := d.getArrayValues(values)
	if !ok {
		return nil
	}

	if len(arrVals) > d.rv.Len() {
		return fmt.Errorf("array capacity exceeded")
	}

	arr := reflect.Indirect(reflect.New(reflect.ArrayOf(d.rv.Len(), d.rv.Type().Elem())))
	for i, v := range arrVals {
		fmt.Println(arr.Index(i))
		arr.Index(i).Set(reflect.ValueOf(v))
	}

	d.rv.Set(arr)
	return nil
}

func (d *decoder) decodeSlice() error {
	if d.rv.Type().Elem().Kind() != reflect.String {
		return fmt.Errorf("allocation type must be []stirng")
	}

	values, err := url.ParseQuery(d.query)
	if err != nil {
		return err
	}

	arrVals, ok := d.getArrayValues(values)
	if !ok {
		return nil
	}

	d.rv.Set(reflect.AppendSlice(d.rv, reflect.ValueOf(arrVals)))
	return nil
}

// returns true if the key can be converted to an array.
//
// valid values
//   key[]
//   key[0]
func (d *decoder) getArrayValues(values url.Values) ([]string, bool) {
	keys := make(map[string]struct{})
	arrValues := make([]string, 0)

	for k, vals := range values {
		startIdx := strings.Index(k, "[")
		closeIdx := strings.Index(k, "]")

		// not array key
		if startIdx == -1 || closeIdx == -1 {
			return nil, false
		}
		// nested array
		if strings.Count(k, "[") > 1 {
			return nil, false
		}

		key := k[:startIdx]
		arrayKey := k[startIdx+1 : closeIdx]

		// no index array
		if arrayKey == "" {
			_, ok := keys[key]
			if len(keys) > 0 && !ok {
				return nil, false
			}
			keys[key] = struct{}{}
			arrValues = append(arrValues, vals...)
			continue
		}

		// map
		if _, err := strconv.Atoi(arrayKey); err != nil {
			return nil, false
		}

		_, ok := keys[key]
		if len(keys) > 0 && !ok {
			return nil, false
		}
		keys[key] = struct{}{}
		arrValues = append(arrValues, vals...)
	}

	return arrValues, true
}
