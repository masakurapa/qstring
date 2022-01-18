package qstring

import (
	"reflect"
	"strconv"
)

func (d *decoder) decodeStruct(rv reflect.Value) error {
	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}
	return d.setStruct(rv, valueMap)
}

func (d *decoder) setTypeVlaue(rt reflect.Type, rv reflect.Value, uv urlValue) error {
	switch rt.Kind() {
	case reflect.Struct:
		return d.setStruct(rv, uv.child)
	case reflect.Bool:
		return d.setBool(rv, uv)
	case reflect.Int:
		return d.setInt(rv, uv)
	case reflect.Int8:
		return d.setInt8(rv, uv)
	case reflect.Int16:
		return d.setInt16(rv, uv)
	case reflect.Int32:
		return d.setInt32(rv, uv)
	case reflect.Int64:
		return d.setInt64(rv, uv)
	case reflect.Uint:
		return d.setUint(rv, uv)
	case reflect.Uint8:
		return d.setUint8(rv, uv)
	case reflect.Uint16:
		return d.setUint16(rv, uv)
	case reflect.Uint32:
		return d.setUint32(rv, uv)
	case reflect.Uint64:
		return d.setUint64(rv, uv)
	case reflect.String:
		return d.setString(rv, uv)
	case reflect.Array:
		return d.setArray(rv, uv)
	case reflect.Slice:
		return d.setSlice(rv, uv)
	case reflect.Map:
		if !uv.hasChild() {
			return &NoAssignableValueError{rt, uv.String()}
		}
		return d.setMap(rv, uv.child)
	}

	return &UnsupportedTypeError{rt}
}

func (d *decoder) isPtr(rv reflect.Value) bool {
	return rv.Kind() == reflect.Ptr
}

func (d *decoder) setStruct(rv reflect.Value, uvm urlValueMap) error {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
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

		val, ok := uvm[tag]
		if !ok {
			continue
		}

		var err error
		frv := rv.FieldByName(f.Name)
		if frv.Kind() == reflect.Ptr {
			err = d.setTypeVlaue(frv.Type().Elem(), frv, val)
		} else {
			err = d.setTypeVlaue(frv.Type(), frv, val)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) setBool(rv reflect.Value, uv urlValue) error {
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

	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseInt(uv.values[0], 10, 64)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := int(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt8(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseInt(uv.values[0], 10, 8)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := int8(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt16(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseInt(uv.values[0], 10, 16)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := int16(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt32(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseInt(uv.values[0], 10, 32)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := int32(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setInt64(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val, err := strconv.ParseInt(uv.values[0], 10, 64)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 64)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint8(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 8)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint8(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint16(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 16)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint16(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint32(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 32)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint32(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setUint64(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 64)
	if err != nil {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uint64(i)
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setString(rv reflect.Value, uv urlValue) error {
	if !uv.hasSingleValue() {
		return &NoAssignableValueError{rv.Type(), uv.String()}
	}

	val := uv.values[0]
	if d.isPtr(rv) {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
	return nil
}

func (d *decoder) setArray(rv reflect.Value, uv urlValue) error {
	val := rv
	if d.isPtr(rv) {
		if !rv.Elem().IsValid() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		val = rv.Elem()
	}

	if uv.hasChild() {
		crt := val.Type().Elem()
		if crt.Kind() == reflect.Ptr {
			crt = crt.Elem()
		}

		if val.Len() < len(uv.child) {
			return &ArrayIndexOutOfRangeDecodeError{val.Type(), val.Len()}
		}

		for i, cuv := range uv.child.sortedChild() {
			crv := reflect.New(crt).Elem()
			err := d.setTypeVlaue(crt, crv, cuv)
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

func (d *decoder) setSlice(rv reflect.Value, uv urlValue) error {
	val := rv
	if d.isPtr(rv) {
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
			err := d.setTypeVlaue(crt, crv, cuv)
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
