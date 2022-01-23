package qstring

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	defaultNilValue = ""
	que             = "?"
)

type encoder struct {
	v url.Values
}

func (e *encoder) encode(v interface{}) (string, error) {
	rv := reflect.ValueOf(v)
	if v == nil || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
		return "", &InvalidEncodeError{reflect.TypeOf(v)}
	}

	e.v = make(url.Values)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		if err := e.encodeMap("", rv); err != nil {
			return "", err
		}
	case reflect.Struct:
		if err := e.encodeStruct("", rv); err != nil {
			return "", err
		}
	case reflect.String:
		return e.encodeString(rv), nil
	default:
		return "", &UnsupportedTypeError{rv.Type()}
	}

	if len(e.v) == 0 {
		return "", nil
	}
	return e.v.Encode(), nil
}

func (e *encoder) encodeByType(key string, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Bool:
		bv := "true"
		if !rv.Bool() {
			bv = "false"
		}
		e.v.Add(key, bv)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		e.v.Add(key, strconv.FormatInt(rv.Int(), 10))
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		e.v.Add(key, strconv.FormatUint(rv.Uint(), 10))
	case reflect.Float32:
		e.v.Add(key, strconv.FormatFloat(rv.Float(), 'f', -1, 32))
	case reflect.Float64:
		e.v.Add(key, strconv.FormatFloat(rv.Float(), 'f', -1, 64))
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
		return e.encodeByType(key, reflect.ValueOf(rv.Interface()))
	case reflect.Ptr:
		if rv.IsNil() {
			e.v.Add(key, defaultNilValue)
			return nil
		}
		return e.encodeByType(key, reflect.Indirect(rv))
	default:
		return &UnsupportedTypeError{rv.Type()}
	}

	return nil
}

func (e *encoder) encodeString(rv reflect.Value) string {
	q := ""
	if strings.HasPrefix(rv.String(), que) {
		q = que
	}

	queries := strings.Split(strings.TrimPrefix(rv.String(), que), "&")
	encoded := make([]string, 0, len(queries))

	for _, s := range queries {
		encoded = append(encoded, strings.ReplaceAll(url.QueryEscape(s), "%3D", "="))
	}
	return q + strings.Join(encoded, "&")
}

func (e *encoder) encodeMap(key string, rv reflect.Value) error {
	if rv.IsNil() || rv.Len() == 0 {
		if key != "" {
			e.v.Add(key, defaultNilValue)
		}
		return nil
	}

	rt := rv.Type()
	if rt.Key().Kind() != reflect.String {
		return &UnsupportedTypeError{rt}
	}

	iter := rv.MapRange()
	for iter.Next() {
		if err := e.encodeByType(e.makeMapKey(key, iter.Key().String()), iter.Value()); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeArray(key string, rv reflect.Value) error {
	if rv.Len() == 0 {
		e.v.Add(key, defaultNilValue)
		return nil
	}

	for i := 0; i < rv.Len(); i++ {
		k := key + "[" + strconv.Itoa(i) + "]"
		if err := e.encodeByType(k, rv.Index(i)); err != nil {
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

		if err := e.encodeByType(e.makeMapKey(key, tag), frv); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) makeMapKey(key, ch string) string {
	if key == "" {
		return ch
	}
	return key + "[" + ch + "]"
}
