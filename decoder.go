package qstringer

import (
	"fmt"
	"net/url"
	"reflect"
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
	// d := decoder{rv: riv}

	switch riv.Kind() {
	default:
		return fmt.Errorf("type %s is not available", riv.Kind().String())
	}

	return nil
}

type decoder struct {
	rv reflect.Value
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

func (d *decoder) decodeString(s string) error {
	q, err := url.QueryUnescape(s)
	if err == nil {
		d.rv.SetString(q)
	}
	return err
}
