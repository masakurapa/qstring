package qstring

import (
	"reflect"
	"sort"
	"strconv"
)

func (d *decoder) decodeMap(rv reflect.Value) error {
	if rv.Type().Key().Kind() != reflect.String || rv.Type().Elem().Kind() != reflect.Interface {
		return &unsupportedTypeError{rv.Type()}
	}

	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}
	return d.setMap(rv, valueMap)
}

func (d *decoder) setMap(rv reflect.Value, uvm urlValueMap) error {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rv.Type()))
	}

	for _, uv := range uvm {
		if uv.isString && len(uv.values) == 1 {
			rv.SetMapIndex(reflect.ValueOf(uv.key), reflect.ValueOf(uv.values[0]))
			continue
		}

		if uv.child == nil || len(uv.child) == 0 {
			rv.SetMapIndex(reflect.ValueOf(uv.key), reflect.ValueOf(uv.values))
			continue
		}

		// nested array or map
		val := d.makeMapValueRecursive(uv.child)
		if q, ok := val.(Q); ok {
			if aq, ok := d.toSlice(q); ok {
				rv.SetMapIndex(reflect.ValueOf(uv.key), reflect.ValueOf(aq))
				continue
			}
		}
		rv.SetMapIndex(reflect.ValueOf(uv.key), reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) makeMapValueRecursive(valueMap urlValueMap) interface{} {
	q := make(Q)
	for _, uv := range valueMap {
		if uv.isString && len(uv.values) == 1 {
			q[uv.key] = uv.values[0]
			continue
		}

		if uv.child == nil || len(uv.child) == 0 {
			q[uv.key] = uv.values
			continue
		}

		// nested array or map
		q[uv.key] = d.makeMapValueRecursive(uv.child)
	}

	if aq, ok := d.toSlice(q); ok {
		return aq
	}
	return q
}

func (d *decoder) toSlice(q Q) (S, bool) {
	type tmpAq struct {
		key   string
		value interface{}
	}
	tmp := make([]tmpAq, 0, len(q))

	for key, value := range q {
		if _, err := strconv.Atoi(key); err != nil {
			return nil, false
		}
		tmp = append(tmp, tmpAq{key: key, value: value})
	}

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].key < tmp[j].key
	})

	aq := make(S, 0, len(q))
	for _, v := range tmp {
		aq = append(aq, v.value)
	}

	return aq, true
}
