package qstringer

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

const (
	// KeyTypeCamel uses struct fields in camel-case
	KeyTypeCamel KeyType = iota
	// KeyTypeSnake uses struct fields in snake-case
	KeyTypeSnake
	// KeyTypeKebab uses struct fields in kebab-case
	KeyTypeKebab
)

var (
	matchFirstCap   = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap     = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchSnake      = regexp.MustCompile("_+([A-Za-z])")
	keyType         = KeyTypeCamel
	outputNilValue  = false
	defaultNilValue = ""
)

// KeyType is the type of the key format
// when using the public field of the struct as the key
type KeyType int

// Q is the type of the query string parameters.
type Q map[string]interface{}

// ArrayQ is a type of query string in array format.
type ArrayQ []interface{}

// MapQ is a type of query string in map format.
type MapQ map[string]interface{}

// Encode returns the URL-encoded query string.
//
// support struct, map type where the key is a string.
//
// add "?" to the beginning and return.
//
// when a struct is used as an argument, the key format defaults to camel-case.
// if you want to change the format, use qstringer.SetKeyType().
func Encode(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("nil is not available")
	}

	e := encoder{v: url.Values{}}
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			// map key must be a string
			if iter.Key().Kind() != reflect.String {
				return "", fmt.Errorf("the key of the map type must be a string")
			}
			if err := e.encode(iter.Key().String(), iter.Value()); err != nil {
				return "", err
			}
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			f := rv.Type().Field(i)
			if !f.IsExported() {
				continue
			}
			if err := e.encode(e.structKey(f.Name), rv.FieldByName(f.Name)); err != nil {
				return "", err
			}
		}
	default:
		return "", fmt.Errorf("type %s is not available", rv.Kind().String())
	}

	if len(e.v) == 0 {
		return "", nil
	}
	return "?" + e.v.Encode(), nil
}

// SetKeyType sets the key format of the struct field.
func SetKeyType(t KeyType) {
	keyType = t
}

// OutputNilValue output nil value if set to true.
//
// default is false
func OutputNilValue(b bool) {
	outputNilValue = b
}

type encoder struct {
	v url.Values
}

func (e *encoder) encode(key string, rv reflect.Value) (err error) {
	if !rv.IsValid() {
		if outputNilValue {
			e.v.Add(key, defaultNilValue)
		}
		return
	}

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
		err = e.encodeMap(key, rv)
	case reflect.Array:
		err = e.encodeArray(key, rv)
	case reflect.Slice:
		err = e.encodeSlice(key, rv)
	case reflect.Struct:
		err = e.encodeStruct(key, rv)
	case reflect.String:
		e.v.Add(key, rv.String())
	case reflect.Interface:
		return e.encode(key, reflect.ValueOf(rv.Interface()))
	default:
		err = fmt.Errorf("type %s is not available (key: %s)", rv.Kind().String(), key)
	}

	return
}

func (e *encoder) encodeMap(key string, rv reflect.Value) error {
	if rv.IsNil() {
		if outputNilValue {
			e.v.Add(key, defaultNilValue)
		}
		return nil
	}

	iter := rv.MapRange()
	for iter.Next() {
		// map key must be a string
		if iter.Key().Kind() != reflect.String {
			return fmt.Errorf("the key of the map type must be a string (key: %s)", key)
		}
		if err := e.encode(fmt.Sprintf("%s[%s]", key, iter.Key().String()), iter.Value()); err != nil {
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
	fmt.Println(rv.IsNil())
	if rv.IsNil() {
		if outputNilValue {
			e.v.Add(key, defaultNilValue)
		}
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
		if err := e.encode(fmt.Sprintf("%s[%s]", key, e.structKey(f.Name)), rv.FieldByName(f.Name)); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) structKey(str string) string {
	switch keyType {
	case KeyTypeSnake:
		return e.toSnakeCase(str)
	case KeyTypeKebab:
		return e.toKebabCase(str)
	default:
		return e.toCamelCase(str)
	}
}

func (e *encoder) toCamelCase(str string) string {
	s := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	s = matchAllCap.ReplaceAllString(s, "${1}_${2}")
	s = strings.ToLower(s)
	return matchSnake.ReplaceAllStringFunc(s, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}

func (e *encoder) toSnakeCase(str string) string {
	s := strings.ReplaceAll(str, "_", "-")
	s = matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	s = matchAllCap.ReplaceAllString(s, "${1}_${2}")
	s = strings.ReplaceAll(s, "-", "")
	return strings.ToLower(s)
}

func (e *encoder) toKebabCase(str string) string {
	s := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	s = matchAllCap.ReplaceAllString(s, "${1}-${2}")
	s = strings.ReplaceAll(s, "_", "")
	return strings.ToLower(s)
}
