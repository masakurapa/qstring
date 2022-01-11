package qstring

import (
	"fmt"
	"reflect"
	"strconv"
)

func (d *decoder) decodeStruct(rv reflect.Value) error {
	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}

	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		if !f.IsExported() {
			continue
		}

		tag, _ := parseTag(f.Tag)
		if tag == "" {
			continue
		}

		val, ok := valueMap[tag]
		if !ok {
			continue
		}

		// if opt.omitempty && isEmptyValue(frv) {
		// 	continue
		// }

		var err error
		frv := rv.FieldByName(f.Name)
		if frv.Kind() == reflect.Ptr {
			err = d.setTypeVlaue(frv.Type().Elem(), frv, val, true)
		} else {
			err = d.setTypeVlaue(frv.Type(), frv, val, false)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) setTypeVlaue(rt reflect.Type, rv reflect.Value, uv urlValue, isPtr bool) error {
	switch rt.Kind() {
	case reflect.Bool:
		return d.setBool(rv, uv, isPtr)
	case reflect.Int:
		return d.setInt(rv, uv, isPtr)
	case reflect.Int8:
		return d.setInt8(rv, uv, isPtr)
	case reflect.Int16:
		return d.setInt16(rv, uv, isPtr)
	case reflect.Int32:
		return d.setInt32(rv, uv, isPtr)
	case reflect.Int64:
		return d.setInt64(rv, uv, isPtr)
	case reflect.Uint:
		return d.setUint(rv, uv, isPtr)
	case reflect.Uint8:
		return d.setUint8(rv, uv, isPtr)
	case reflect.Uint16:
		return d.setUint16(rv, uv, isPtr)
	case reflect.Uint32:
		return d.setUint32(rv, uv, isPtr)
	case reflect.Uint64:
		return d.setUint64(rv, uv, isPtr)
	case reflect.String:
		return d.setString(rv, uv, isPtr)
	case reflect.Array:
		return d.setArray(rv, uv, isPtr)
	case reflect.Slice:
		return d.setSlice(rv, uv, isPtr)
	}

	fmt.Println(rt.String())

	return &UnsupportedTypeError{rt}
}

func (d *decoder) setBool(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := false
	switch v := uv.values[0]; v {
	case "0", "false":
		// skip
	case "1", "true":
		val = true
	default:
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseInt(uv.values[0], 10, 64)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := int(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt8(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseInt(uv.values[0], 10, 8)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := int8(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt16(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseInt(uv.values[0], 10, 16)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := int16(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt32(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseInt(uv.values[0], 10, 32)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := int32(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt64(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val, err := strconv.ParseInt(uv.values[0], 10, 64)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 64)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint8(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 8)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint8(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint16(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 16)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint16(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint32(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 32)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint32(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint64(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 64)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint64(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setString(rv reflect.Value, uv urlValue, isPtr bool) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uv.values[0]
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setArray(rv reflect.Value, uv urlValue, isPtr bool) error {
	val := rv
	if isPtr {
		if !rv.Elem().IsValid() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		val = rv.Elem()
	}

	if uv.hasChild() {
		crt := val.Type().Elem()
		cPtr := crt.Kind() == reflect.Ptr
		if cPtr {
			crt = crt.Elem()
		}

		if val.Len() < len(uv.child) {
			return &ArrayIndexOutOfRangeDecodeError{val.Type(), val.Len()}
		}

		for i, cuv := range uv.child.sortedChild() {
			crv := reflect.New(crt).Elem()
			err := d.setTypeVlaue(crt, crv, cuv, cPtr)
			if err != nil {
				return err
			}
			val.Index(i).Set(crv)
		}

		return nil
	}

	if val.Len() < len(uv.values) {
		return &ArrayIndexOutOfRangeDecodeError{val.Type(), val.Len()}
	}

	if val.Index(0).Type().Kind() != reflect.String {
		return &UnsupportedTypeError{val.Type()}
	}

	for i, v := range uv.values {
		val.Index(i).Set(reflect.ValueOf(v))
	}
	return nil
}

func (d *decoder) setSlice(rv reflect.Value, uv urlValue, isPtr bool) error {
	val := rv
	if isPtr {
		if !rv.Elem().IsValid() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		val = rv.Elem()
	}

	if uv.hasChild() {
		crt := val.Type().Elem()
		cPtr := crt.Kind() == reflect.Ptr
		if cPtr {
			crt = crt.Elem()
		}

		for _, cuv := range uv.child.sortedChild() {
			crv := reflect.New(crt).Elem()
			err := d.setTypeVlaue(crt, crv, cuv, cPtr)
			if err != nil {
				return err
			}
			val.Set(reflect.Append(val, crv))
		}

		return nil
	}

	if !val.Type().AssignableTo(reflect.TypeOf(uv.values)) {
		return &UnsupportedTypeError{val.Type()}
	}

	val.Set(reflect.AppendSlice(val, reflect.ValueOf(uv.values)))
	return nil
}
