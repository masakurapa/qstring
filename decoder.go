package qstring

import (
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type decoder struct {
	query string
}

func (d *decoder) decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &invalidDecodeError{reflect.TypeOf(v)}
	}

	rv = rv.Elem()
	switch rv.Kind() {
	case reflect.String:
		return d.decodeString(rv)
	case reflect.Array:
		return d.decodeArray(rv)
	case reflect.Slice:
		return d.decodeSlice(rv)
	case reflect.Map:
		return d.decodeMap(rv)
	case reflect.Struct:
		return d.decodeStruct(rv)
	}

	return &unsupportedTypeError{rv.Type()}
}

func (d *decoder) decodeString(rv reflect.Value) error {
	q, err := url.QueryUnescape(d.query)
	if err == nil {
		rv.SetString(q)
	}
	return err
}

func (d *decoder) decodeArray(rv reflect.Value) error {
	if rv.Type().Elem().Kind() != reflect.String {
		return &unsupportedTypeError{rv.Type()}
	}

	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}

	if len(valueMap) != 1 {
		return &multipleKeysError{}
	}

	arrVals := valueMap.firstValue()

	if len(arrVals) > rv.Len() {
		return &arrayIndexOutOfRangeDecodeError{rv.Type(), len(arrVals)}
	}

	arr := reflect.Indirect(reflect.New(reflect.ArrayOf(rv.Len(), rv.Type().Elem())))
	for i, v := range arrVals {
		arr.Index(i).Set(reflect.ValueOf(v))
	}

	rv.Set(arr)
	return nil
}

func (d *decoder) decodeSlice(rv reflect.Value) error {
	if rv.Type().Elem().Kind() != reflect.String {
		return &unsupportedTypeError{rv.Type()}
	}

	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}

	if len(valueMap) != 1 {
		return &multipleKeysError{}
	}

	rv.Set(reflect.AppendSlice(rv, reflect.ValueOf(valueMap.firstValue())))
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
	tmp := make([]urlValue, 0, len(valueMap))

	for _, uv := range valueMap {
		// nested array or map
		if uv.child != nil && len(uv.child) > 0 {
			return nil, false
		}

		if uv.key == "" {
			tmp = append(tmp, uv)
			continue
		}

		// if not number, map type
		if _, err := strconv.Atoi(uv.key); err != nil {
			return nil, false
		}

		tmp = append(tmp, uv)
	}

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].key < tmp[j].key
	})

	ret := make([]string, 0, len(valueMap))
	for _, v := range tmp {
		ret = append(ret, v.values...)
	}

	return ret, true
}
