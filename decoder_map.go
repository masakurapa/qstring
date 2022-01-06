package qstringer

import (
	"errors"
	"reflect"
	"sort"
	"strconv"
)

func DecodeToMap(s string) (Q, error) {
	var q Q
	err := Decode(s, &q)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (d *decoder) decodeMap() error {
	if !d.rv.Type().Key().AssignableTo(reflect.TypeOf("")) {
		return errors.New("map key must be assignable to a string")
	}
	if d.rv.Type().Elem().Kind() != reflect.Interface {
		return errors.New("map value must be assignable to interface{}")
	}

	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}

	if d.rv.IsNil() {
		d.rv.Set(reflect.MakeMap(d.rv.Type()))
	}

	for _, uv := range valueMap {
		if uv.isString && len(uv.values) == 1 {
			d.rv.SetMapIndex(reflect.ValueOf(uv.key), reflect.ValueOf(uv.values[0]))
			continue
		}

		if uv.child == nil || len(uv.child) == 0 {
			d.rv.SetMapIndex(reflect.ValueOf(uv.key), reflect.ValueOf(uv.values))
			continue
		}

		// nested array or map
		val := d.makeMapValueRecursive(uv.child)
		if q, ok := val.(Q); ok {
			if aq, ok := d.toArrayQ(q); ok {
				d.rv.SetMapIndex(reflect.ValueOf(uv.key), reflect.ValueOf(aq))
				continue
			}
		}
		d.rv.SetMapIndex(reflect.ValueOf(uv.key), reflect.ValueOf(val))
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

	if aq, ok := d.toArrayQ(q); ok {
		return aq
	}
	return q
}

func (d *decoder) toArrayQ(q Q) (ArrayQ, bool) {
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

	aq := make(ArrayQ, 0, len(q))
	for _, v := range tmp {
		aq = append(aq, v.value)
	}

	return aq, true
}
