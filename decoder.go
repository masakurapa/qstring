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

	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}

	if len(valueMap) != 1 {
		return nil
	}

	arrVals := valueMap.firstValue()

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

	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}

	if len(valueMap) != 1 {
		return nil
	}

	d.rv.Set(reflect.AppendSlice(d.rv, reflect.ValueOf(valueMap.firstValue())))
	return nil
}

func (d *decoder) createIntermediateStruct() (urlValueMap, error) {
	urlValues, err := url.ParseQuery(d.query)
	if err != nil {
		return nil, err
	}

	valueMap := make(urlValueMap)

	for key, values := range urlValues {
		// convert `key[a][b]` to `[]string{"key", "a", "b"}`
		splitKeys := strings.Split(key, "[")
		for i, v := range splitKeys {
			splitKeys[i] = strings.TrimSuffix(v, "]")
		}

		k := splitKeys[0]
		if _, ok := valueMap[k]; ok {
			valueMap[k] = d.toUrlValue(valueMap[k], splitKeys, values)
			continue
		}
		valueMap[k] = d.toUrlValue(urlValue{}, splitKeys, values)
	}

	return d.conpact(valueMap), nil
}

func (d *decoder) toUrlValue(uv urlValue, keys []string, values []string) urlValue {
	key := keys[0]
	uv.key = key

	if len(keys) == 1 {
		uv.values = append(uv.values, values...)
		uv.isString = true
		return uv
	}

	if uv.child == nil {
		uv.child = make(urlValueMap)
	}

	nextKey := keys[1]
	if _, ok := uv.child[nextKey]; ok {
		uv.child[nextKey] = d.toUrlValue(uv.child[nextKey], keys[1:], values)
		return uv
	}

	uv.child[nextKey] = d.toUrlValue(urlValue{}, keys[1:], values)
	return uv
}

func (d *decoder) conpact(valueMap urlValueMap) urlValueMap {
	newMap := make(urlValueMap)

	for _, uv := range valueMap {
		uv := uv

		// not has child value
		if uv.child == nil || len(uv.child) == 0 {
			newMap[uv.key] = uv
			continue
		}

		// simple array value
		if arrVals, ok := d.toArray(uv.child); ok {
			uv.values = arrVals
			uv.child = nil
			newMap[uv.key] = uv
			continue
		}

		// nested array, map
		uv.child = d.conpact(uv.child)
		newMap[uv.key] = uv
	}

	return newMap
}

func (d *decoder) toArray(valueMap urlValueMap) ([]string, bool) {
	ret := make([]string, 0, len(valueMap))

	for _, uv := range valueMap {
		// nested array or map
		if uv.child != nil && len(uv.child) > 0 {
			return nil, false
		}

		if uv.key == "" {
			ret = append(ret, uv.values...)
			continue
		}

		// if not number, map type
		if _, err := strconv.Atoi(uv.key); err != nil {
			return nil, false
		}

		ret = append(ret, uv.values...)
	}
	return ret, true
}

type urlValue struct {
	key      string
	values   []string
	isString bool
	child    urlValueMap
}

type urlValueMap map[string]urlValue

func (vm urlValueMap) firstValue() []string {
	var val []string
	for _, v := range vm {
		val = v.values
	}
	return val
}
