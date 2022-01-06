package qstringer

import (
	"math"
	"reflect"
	"strconv"
)

func (d *decoder) decodeStruct() error {
	valueMap, err := d.createIntermediateStruct()
	if err != nil {
		return err
	}

	for i := 0; i < d.rv.NumField(); i++ {
		f := d.rv.Type().Field(i)
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

		frv := d.rv.FieldByName(f.Name)
		if frv.Kind() == reflect.Ptr {
			d.setTypeVlaue(frv.Type().Elem(), frv, val, true)
		} else {
			d.setTypeVlaue(frv.Type(), frv, val, false)
		}
	}
	return nil
}

func (d *decoder) setTypeVlaue(rt reflect.Type, rv reflect.Value, uv urlValue, isPtr bool) {
	switch rt.Kind() {
	case reflect.Bool:
		d.setBool(rv, uv, isPtr)
	case reflect.Int:
		d.setInt(rv, uv, isPtr)
	case reflect.Int8:
		d.setInt8(rv, uv, isPtr)
	case reflect.Int16:
		d.setInt16(rv, uv, isPtr)
	case reflect.Int32:
		d.setInt32(rv, uv, isPtr)
	case reflect.Int64:
		d.setInt64(rv, uv, isPtr)
	case reflect.Uint:
		d.setUint(rv, uv, isPtr)
	case reflect.Uint8:
		d.setUint8(rv, uv, isPtr)
	case reflect.Uint16:
		d.setUint16(rv, uv, isPtr)
	case reflect.Uint32:
		d.setUint32(rv, uv, isPtr)
	case reflect.Uint64:
		d.setUint64(rv, uv, isPtr)
	case reflect.String:
		d.setString(rv, uv, isPtr)
	}
}

func (d *decoder) setBool(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	val := false
	switch v := uv.values[0]; v {
	case "0", "false":
		// skip
	case "1", "true":
		val = true
	default:
		return
	}

	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setInt(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	val, err := strconv.Atoi(uv.values[0])
	if err != nil {
		return
	}

	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setInt8(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.Atoi(uv.values[0])
	if err != nil {
		return
	}
	if math.MinInt8 > i || math.MaxInt8 < i {
		return
	}

	val := int8(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setInt16(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.Atoi(uv.values[0])
	if err != nil {
		return
	}
	if math.MinInt16 > i || math.MaxInt16 < i {
		return
	}

	val := int16(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setInt32(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.Atoi(uv.values[0])
	if err != nil {
		return
	}
	if math.MinInt32 > i || math.MaxInt32 < i {
		return
	}

	val := int32(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setInt64(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.Atoi(uv.values[0])
	if err != nil {
		return
	}

	val := int64(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setUint(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 64)
	if err != nil {
		return
	}

	val := uint(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setUint8(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 8)
	if err != nil {
		return
	}

	val := uint8(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setUint16(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 16)
	if err != nil {
		return
	}

	val := uint16(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setUint32(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 32)
	if err != nil {
		return
	}

	val := uint32(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setUint64(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	i, err := strconv.ParseUint(uv.values[0], 10, 64)
	if err != nil {
		return
	}

	val := uint64(i)
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func (d *decoder) setString(rv reflect.Value, uv urlValue, isPtr bool) {
	if !uv.hasSingleValue() {
		return
	}

	val := uv.values[0]
	if isPtr {
		rv.Set(reflect.ValueOf(&val))
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}
