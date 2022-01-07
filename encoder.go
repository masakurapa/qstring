package qstring

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
		q := ""
		if strings.HasPrefix(rv.String(), que) {
			q = que
		}
		s := strings.TrimPrefix(rv.String(), que)
		return q + strings.ReplaceAll(url.QueryEscape(s), "%3D", "="), nil
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

func (e *encoder) encodeMap(key string, rv reflect.Value) error {
	if rv.IsNil() {
		e.v.Add(key, defaultNilValue)
		return nil
	}

	iter := rv.MapRange()
	for iter.Next() {
		// map key must be a string
		if iter.Key().Kind() != reflect.String {
			return &UnsupportedTypeError{rv.Type()}
		}

		if err := e.encodeByType(e.makeMapKey(key, iter.Key().String()), iter.Value()); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeArray(key string, rv reflect.Value) error {
	for i := 0; i < rv.Len(); i++ {
		if err := e.encodeByType(fmt.Sprintf("%s[%d]", key, i), rv.Index(i)); err != nil {
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
	return fmt.Sprintf("%s[%s]", key, ch)
}
