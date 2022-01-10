package qstring_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/masakurapa/qstring"
)

type ds struct {
	FieldB       bool         `qstring:"field_b"`
	FieldI       int          `qstring:"fieldI"`
	FieldI8      int8         `qstring:"fieldI8"`
	FieldI16     int16        `qstring:"fieldI16"`
	FieldI32     int32        `qstring:"fieldI32"`
	FieldI64     int64        `qstring:"fieldI64"`
	FieldUI8     uint8        `qstring:"fieldUI8"`
	FieldUI      uint         `qstring:"fieldUI"`
	FieldUI16    uint16       `qstring:"fieldUI16"`
	FieldUI32    uint32       `qstring:"fieldUI32"`
	FieldUI64    uint64       `qstring:"fieldUI64"`
	FieldFloat32 float32      `qstring:"fieldFloat32"`
	FieldFloat64 float64      `qstring:"fieldFloat64"`
	JSONStr      string       `qstring:"json_str"`
	Array        [3]string    `qstring:"array"`
	ArrayNest    [3][3]string `qstring:"arrayNest"`
	ArrayI       [3]int       `qstring:"arrayI"`
	Slice        []string     `qstring:"slice"`
	SliceNest    [][]string   `qstring:"sliceNest"`
	SliceI       []int        `qstring:"sliceI"`
}

type dsp struct {
	FieldB       *bool          `qstring:"field_b"`
	FieldI       *int           `qstring:"fieldI"`
	FieldI8      *int8          `qstring:"fieldI8"`
	FieldI16     *int16         `qstring:"fieldI16"`
	FieldI32     *int32         `qstring:"fieldI32"`
	FieldI64     *int64         `qstring:"fieldI64"`
	FieldUI8     *uint8         `qstring:"fieldUI8"`
	FieldUI      *uint          `qstring:"fieldUI"`
	FieldUI16    *uint16        `qstring:"fieldUI16"`
	FieldUI32    *uint32        `qstring:"fieldUI32"`
	FieldUI64    *uint64        `qstring:"fieldUI64"`
	FieldFloat32 *float32       `qstring:"fieldFloat32"`
	FieldFloat64 *float64       `qstring:"fieldFloat64"`
	JSONStr      *string        `qstring:"json_str"`
	Array        *[3]string     `qstring:"array"`
	ArrayNest    *[3]*[3]string `qstring:"arrayNest"`
	ArrayI       *[3]int        `qstring:"arrayI"`
	Slice        *[]string      `qstring:"slice"`
	SliceNest    *[]*[]string   `qstring:"sliceNest"`
	SliceI       *[]int         `qstring:"sliceI"`
}

// type ds2 struct {
// 	Field    string `qstring:"field"`
// 	NoTag    string
// 	privateS string `qstring:"private-S"`
// }

func TestDecode(t *testing.T) {
	q := "test=1"

	testCases := []struct {
		name     string
		q        string
		v        interface{}
		err      error
		expected interface{}
	}{
		{name: "nil", q: q, v: nil, err: fmt.Errorf("nil is not supported")},
		{name: "not pointer", q: q, v: func() string { return "" }(), err: fmt.Errorf("non-pointer is not supported")},

		// unsupported type
		{name: "bool type", q: q, v: boolP(false), err: fmt.Errorf("bool is not supported")},
		{name: "int type", q: q, v: intP(0), err: fmt.Errorf("int is not supported")},
		{name: "int64 type", q: q, v: int64P(0), err: fmt.Errorf("int64 is not supported")},
		{name: "int32 type", q: q, v: int32P(0), err: fmt.Errorf("int32 is not supported")},
		{name: "int16 type", q: q, v: int16P(0), err: fmt.Errorf("int16 is not supported")},
		{name: "int8 type", q: q, v: int8P(0), err: fmt.Errorf("int8 is not supported")},
		{name: "uint type", q: q, v: uintP(0), err: fmt.Errorf("uint is not supported")},
		{name: "uint64 type", q: q, v: uint64P(0), err: fmt.Errorf("uint64 is not supported")},
		{name: "uint32 type", q: q, v: uint32P(0), err: fmt.Errorf("uint32 is not supported")},
		{name: "uint16 type", q: q, v: uint16P(0), err: fmt.Errorf("uint16 is not supported")},
		{name: "uint8 type", q: q, v: uint8P(0), err: fmt.Errorf("uint8 is not supported")},
		{name: "float64 type", q: q, v: float64P(0), err: fmt.Errorf("float64 is not supported")},
		{name: "float32 type", q: q, v: float32P(0), err: fmt.Errorf("float32 is not supported")},
		// {name: "array type", q: q, v: func() *bool { var a bool; return &a }(), err: fmt.Errorf("array is not supported")},
		// {name: "slice type", q: q, v: func() *bool { var a bool; return &a }(), err: fmt.Errorf("slice is not supported")},
		{name: "uintptr type", q: q, v: func() *uintptr { var a uintptr; return &a }(), err: fmt.Errorf("uintptr is not supported")},
		{name: "complex64 type", q: q, v: func() *complex64 { var a complex64; return &a }(), err: fmt.Errorf("complex64 is not supported")},
		{name: "complex128 type", q: q, v: func() *complex128 { var a complex128; return &a }(), err: fmt.Errorf("complex128 is not supported")},
		// {name: "chan type", q: q, v: func() *chan { var a chan; return &a }(), err: fmt.Errorf("chan is not supported")},
		{name: "func type", q: q, v: func() *string { return nil }, err: fmt.Errorf("non-pointer is not supported")},
		{name: "nil ptr type", q: q, v: func() *bool { return nil }(), err: fmt.Errorf("nil ptr is not supported")},
		{name: "unsafe pointer type", q: q, v: func() *unsafe.Pointer { var a unsafe.Pointer; return &a }(), err: fmt.Errorf("unsafe.Pointer is not supported")},

		// string type
		{name: "string type with quote", q: "?hoge[key]=fuga", v: stringP(""), expected: "?hoge[key]=fuga"},
		{name: "string type without quote", q: "hoge[key]=fuga", v: stringP(""), expected: "hoge[key]=fuga"},

		// array type
		{name: "array of string type", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *[3]string { var a [3]string; return &a }(), expected: [3]string{"a", "2", "3"}},
		{name: "array of string type - capacity exceeded", q: "hoge[]=a&hoge[]=2&hoge[]=3&hoge[]=4", v: func() *[3]string { var a [3]string; return &a }(), err: fmt.Errorf("index out of range [4] with [3]string")},
		{name: "array of string type - multiple key name", q: "hoge[]=a&fuga[]=2", v: func() *[3]string { var a [3]string; return &a }(), expected: [3]string{}},
		{name: "array of int type", q: "hoge[]=1", v: func() *[3]int { var a [3]int; return &a }(), err: fmt.Errorf("[3]int is not supported")},

		// slice type
		{name: "slice of string type", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *[]string { var a []string; return &a }(), expected: []string{"a", "2", "3"}},
		{name: "slice of int type", q: "hoge[]=1", v: func() *[]int { var a []int; return &a }(), err: fmt.Errorf("[]int is not supported")},
		{name: "slice of string type - multiple key name", q: "hoge[]=a&fuga[]=2", v: func() *[]string { var a []string; return &a }(), expected: func() []string { var a []string; return a }()},

		// map type
		{name: "map type", q: "hoge=a", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": "a"}},
		{name: "map type - multiple key name", q: "hoge=a&fuga=1", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": "a", "fuga": "1"}},
		{name: "map type - duplicate key name", q: "hoge=a&hoge=1", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": []string{"a", "1"}}},
		{name: "map type - array key", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": []string{"a", "2", "3"}}},
		{name: "map type - index array key", q: "hoge[0]=a&hoge[1]=2&hoge[2]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": []string{"a", "2", "3"}}},

		{name: "map type - nested array and array", q: "hoge[0][0]=a&hoge[0][1]=2&hoge[0][2]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": qstring.ArrayQ{[]string{"a", "2", "3"}}}},
		{name: "map type - nested map and array", q: "hoge[a][0]=a&hoge[a][1]=2&hoge[a][2]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": qstring.Q{"a": []string{"a", "2", "3"}}}},
		{name: "map type - nested map and map", q: "hoge[a][b]=a&hoge[a][1]=2&hoge[0][a]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": qstring.Q{"0": qstring.Q{"a": "3"}, "a": qstring.Q{"1": "2", "b": "a"}}}},

		// struct type
		{name: "struct type - bool", q: "field_b=1", v: dsP(), expected: ds{FieldB: true}},
		{name: "struct type - bool - not bool", q: "field_b=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to bool`)},
		{name: "struct type - bool pointer", q: "field_b=1", v: dspP(), expected: dsp{FieldB: boolP(true)}},
		{name: "struct type - bool pointer - not bool", q: "field_b=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *bool`)},

		{name: "struct type - int - min", q: "fieldI=-9223372036854775808", v: dsP(), expected: ds{FieldI: -9223372036854775808}},
		{name: "struct type - int - min out of range", q: "fieldI=-9223372036854775809", v: dsP(), err: fmt.Errorf(`"-9223372036854775809" can not be assign to int`)},
		{name: "struct type - int - max", q: "fieldI=9223372036854775807", v: dsP(), expected: ds{FieldI: 9223372036854775807}},
		{name: "struct type - int - max out of range", q: "fieldI=9223372036854775808", v: dsP(), err: fmt.Errorf(`"9223372036854775808" can not be assign to int`)},
		{name: "struct type - int - not int", q: "fieldI=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to int`)},
		{name: "struct type - int pointer - min", q: "fieldI=-9223372036854775808", v: dspP(), expected: dsp{FieldI: intP(-9223372036854775808)}},
		{name: "struct type - int pointer - min out of range", q: "fieldI=-9223372036854775809", v: dspP(), err: fmt.Errorf(`"-9223372036854775809" can not be assign to *int`)},
		{name: "struct type - int pointer - max", q: "fieldI=9223372036854775807", v: dspP(), expected: dsp{FieldI: intP(9223372036854775807)}},
		{name: "struct type - int pointer - max out of range", q: "fieldI=9223372036854775808", v: dspP(), err: fmt.Errorf(`"9223372036854775808" can not be assign to *int`)},
		{name: "struct type - int pointer - not int", q: "fieldI=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *int`)},

		{name: "struct type - int8 - min", q: "fieldI8=-128", v: dsP(), expected: ds{FieldI8: -128}},
		{name: "struct type - int8 - min out of range", q: "fieldI8=-129", v: dsP(), err: fmt.Errorf(`"-129" can not be assign to int8`)},
		{name: "struct type - int8 - max", q: "fieldI8=127", v: dsP(), expected: ds{FieldI8: 127}},
		{name: "struct type - int8 - max out of range", q: "fieldI8=128", v: dsP(), err: fmt.Errorf(`"128" can not be assign to int8`)},
		{name: "struct type - int8 - not int8", q: "fieldI8=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to int8`)},
		{name: "struct type - int8 pointer - min", q: "fieldI8=-128", v: dspP(), expected: dsp{FieldI8: int8P(-128)}},
		{name: "struct type - int8 pointer - min out of range", q: "fieldI8=-129", v: dspP(), err: fmt.Errorf(`"-129" can not be assign to *int8`)},
		{name: "struct type - int8 pointer - max", q: "fieldI8=127", v: dspP(), expected: dsp{FieldI8: int8P(127)}},
		{name: "struct type - int8 pointer - max out of range", q: "fieldI8=128", v: dspP(), err: fmt.Errorf(`"128" can not be assign to *int8`)},
		{name: "struct type - int8 pointer - not int8", q: "fieldI8=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *int8`)},

		{name: "struct type - int16 - min", q: "fieldI16=-32768", v: dsP(), expected: ds{FieldI16: -32768}},
		{name: "struct type - int16 - min out of range", q: "fieldI16=-32769", v: dsP(), err: fmt.Errorf(`"-32769" can not be assign to int16`)},
		{name: "struct type - int16 - max", q: "fieldI16=32767", v: dsP(), expected: ds{FieldI16: 32767}},
		{name: "struct type - int16 - max out of range", q: "fieldI16=32768", v: dsP(), err: fmt.Errorf(`"32768" can not be assign to int16`)},
		{name: "struct type - int16 - not int16", q: "fieldI16=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to int16`)},
		{name: "struct type - int16 pointer - min", q: "fieldI16=-32768", v: dspP(), expected: dsp{FieldI16: int16P(-32768)}},
		{name: "struct type - int16 pointer - min out of range", q: "fieldI16=-32769", v: dspP(), err: fmt.Errorf(`"-32769" can not be assign to *int16`)},
		{name: "struct type - int16 pointer - max", q: "fieldI16=32767", v: dspP(), expected: dsp{FieldI16: int16P(32767)}},
		{name: "struct type - int16 pointer - max out of range", q: "fieldI16=32768", v: dspP(), err: fmt.Errorf(`"32768" can not be assign to *int16`)},
		{name: "struct type - int16 pointer - not int8", q: "fieldI16=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *int16`)},

		{name: "struct type - int32 - min", q: "fieldI32=-2147483648", v: dsP(), expected: ds{FieldI32: -2147483648}},
		{name: "struct type - int32 - min out of range", q: "fieldI32=-2147483649", v: dsP(), err: fmt.Errorf(`"-2147483649" can not be assign to int32`)},
		{name: "struct type - int32 - max", q: "fieldI32=2147483647", v: dsP(), expected: ds{FieldI32: 2147483647}},
		{name: "struct type - int32 - max out of range", q: "fieldI32=2147483648", v: dsP(), err: fmt.Errorf(`"2147483648" can not be assign to int32`)},
		{name: "struct type - int32 - not int32", q: "fieldI32=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to int32`)},
		{name: "struct type - int32 pointer - min", q: "fieldI32=-2147483648", v: dspP(), expected: dsp{FieldI32: int32P(-2147483648)}},
		{name: "struct type - int32 pointer - min out of range", q: "fieldI32=-2147483649", v: dspP(), err: fmt.Errorf(`"-2147483649" can not be assign to *int32`)},
		{name: "struct type - int32 pointer - max", q: "fieldI32=2147483647", v: dspP(), expected: dsp{FieldI32: int32P(2147483647)}},
		{name: "struct type - int32 pointer - max out of range", q: "fieldI32=2147483648", v: dspP(), err: fmt.Errorf(`"2147483648" can not be assign to *int32`)},
		{name: "struct type - int32 pointer - not int32", q: "fieldI32=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *int32`)},

		{name: "struct type - int64 - min", q: "fieldI64=-9223372036854775808", v: dsP(), expected: ds{FieldI64: -9223372036854775808}},
		{name: "struct type - int64 - min out of range", q: "fieldI64=-9223372036854775809", v: dsP(), err: fmt.Errorf(`"-9223372036854775809" can not be assign to int64`)},
		{name: "struct type - int64 - max", q: "fieldI64=9223372036854775807", v: dsP(), expected: ds{FieldI64: 9223372036854775807}},
		{name: "struct type - int64 - max out of range", q: "fieldI64=9223372036854775808", v: dsP(), err: fmt.Errorf(`"9223372036854775808" can not be assign to int64`)},
		{name: "struct type - int64 - not int64", q: "fieldI64=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to int64`)},
		{name: "struct type - int64 pointer - min", q: "fieldI64=-9223372036854775808", v: dspP(), expected: dsp{FieldI64: int64P(-9223372036854775808)}},
		{name: "struct type - int64 pointer - min out of range", q: "fieldI64=-9223372036854775809", v: dspP(), err: fmt.Errorf(`"-9223372036854775809" can not be assign to *int64`)},
		{name: "struct type - int64 pointer - max", q: "fieldI64=9223372036854775807", v: dspP(), expected: dsp{FieldI64: int64P(9223372036854775807)}},
		{name: "struct type - int64 pointer - max out of range", q: "fieldI64=9223372036854775808", v: dspP(), err: fmt.Errorf(`"9223372036854775808" can not be assign to *int64`)},
		{name: "struct type - int64 pointer - not int64", q: "fieldI64=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *int64`)},

		{name: "struct type - uint - min", q: "fieldUI=0", v: dsP(), expected: ds{FieldUI: 0}},
		{name: "struct type - uint - min out of range", q: "fieldUI=-1", v: dsP(), err: fmt.Errorf(`"-1" can not be assign to uint`)},
		{name: "struct type - uint - max", q: "fieldUI=18446744073709551615", v: dsP(), expected: ds{FieldUI: 18446744073709551615}},
		{name: "struct type - uint - max out of range", q: "fieldUI=18446744073709551616", v: dsP(), err: fmt.Errorf(`"18446744073709551616" can not be assign to uint`)},
		{name: "struct type - uint - not uint", q: "fieldUI=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to uint`)},
		{name: "struct type - uint pointer - min", q: "fieldUI=0", v: dspP(), expected: dsp{FieldUI: uintP(0)}},
		{name: "struct type - uint pointer - min out of range", q: "fieldUI=-1", v: dspP(), err: fmt.Errorf(`"-1" can not be assign to *uint`)},
		{name: "struct type - uint pointer - max", q: "fieldUI=18446744073709551615", v: dspP(), expected: dsp{FieldUI: uintP(18446744073709551615)}},
		{name: "struct type - uint pointer - max out of range", q: "fieldUI=18446744073709551616", v: dspP(), err: fmt.Errorf(`"18446744073709551616" can not be assign to *uint`)},
		{name: "struct type - uint pointer - not uint", q: "fieldUI=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *uint`)},

		{name: "struct type - uint8 - min", q: "fieldUI8=0", v: dsP(), expected: ds{FieldUI8: 0}},
		{name: "struct type - uint8 - min out of range", q: "fieldUI8=-1", v: dsP(), err: fmt.Errorf(`"-1" can not be assign to uint8`)},
		{name: "struct type - uint8 - max", q: "fieldUI8=128", v: dsP(), expected: ds{FieldUI8: 128}},
		{name: "struct type - uint8 - max out of range", q: "fieldUI8=256", v: dsP(), err: fmt.Errorf(`"256" can not be assign to uint8`)},
		{name: "struct type - uint8 - not uint8", q: "fieldUI8=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to uint8`)},
		{name: "struct type - uint8 pointer - min", q: "fieldUI8=0", v: dspP(), expected: dsp{FieldUI8: uint8P(0)}},
		{name: "struct type - uint8 pointer - min out of range", q: "fieldUI8=-1", v: dspP(), err: fmt.Errorf(`"-1" can not be assign to *uint8`)},
		{name: "struct type - uint8 pointer - max", q: "fieldUI8=128", v: dspP(), expected: dsp{FieldUI8: uint8P(128)}},
		{name: "struct type - uint8 pointer - max out of range", q: "fieldUI8=256", v: dspP(), err: fmt.Errorf(`"256" can not be assign to *uint8`)},
		{name: "struct type - uint8 pointer - not uint8", q: "fieldUI8=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *uint8`)},

		{name: "struct type - uint16 - min", q: "fieldUI16=0", v: dsP(), expected: ds{FieldUI16: 0}},
		{name: "struct type - uint16 - min out of range", q: "fieldUI16=-1", v: dsP(), err: fmt.Errorf(`"-1" can not be assign to uint16`)},
		{name: "struct type - uint16 - max", q: "fieldUI16=65535", v: dsP(), expected: ds{FieldUI16: 65535}},
		{name: "struct type - uint16 - max out of range", q: "fieldUI16=65536", v: dsP(), err: fmt.Errorf(`"65536" can not be assign to uint16`)},
		{name: "struct type - uint16 - not uint16", q: "fieldUI16=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to uint16`)},
		{name: "struct type - uint16 pointer - min", q: "fieldUI16=0", v: dspP(), expected: dsp{FieldUI16: uint16P(0)}},
		{name: "struct type - uint16 pointer - min out of range", q: "fieldUI16=-1", v: dspP(), err: fmt.Errorf(`"-1" can not be assign to *uint16`)},
		{name: "struct type - uint16 pointer - max", q: "fieldUI16=65535", v: dspP(), expected: dsp{FieldUI16: uint16P(65535)}},
		{name: "struct type - uint16 pointer - max out of range", q: "fieldUI16=65536", v: dspP(), err: fmt.Errorf(`"65536" can not be assign to *uint16`)},
		{name: "struct type - uint16 pointer - not uint16", q: "fieldUI16=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *uint16`)},

		{name: "struct type - uint32 - min", q: "fieldUI32=0", v: dsP(), expected: ds{FieldUI32: 0}},
		{name: "struct type - uint32 - min out of range", q: "fieldUI32=-1", v: dsP(), err: fmt.Errorf(`"-1" can not be assign to uint32`)},
		{name: "struct type - uint32 - max", q: "fieldUI32=4294967295", v: dsP(), expected: ds{FieldUI32: 4294967295}},
		{name: "struct type - uint32 - max out of range", q: "fieldUI32=4294967296", v: dsP(), err: fmt.Errorf(`"4294967296" can not be assign to uint32`)},
		{name: "struct type - uint32 - not uint32", q: "fieldUI32=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to uint32`)},
		{name: "struct type - uint32 pointer - min", q: "fieldUI32=0", v: dspP(), expected: dsp{FieldUI32: uint32P(0)}},
		{name: "struct type - uint32 pointer - min out of range", q: "fieldUI32=-1", v: dspP(), err: fmt.Errorf(`"-1" can not be assign to *uint32`)},
		{name: "struct type - uint32 pointer - max", q: "fieldUI32=4294967295", v: dspP(), expected: dsp{FieldUI32: uint32P(4294967295)}},
		{name: "struct type - uint32 pointer - max out of range", q: "fieldUI32=4294967296", v: dspP(), err: fmt.Errorf(`"4294967296" can not be assign to *uint32`)},
		{name: "struct type - uint32 pointer - not uint32", q: "fieldUI32=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *uint32`)},

		{name: "struct type - uint64 - min", q: "fieldUI64=0", v: dsP(), expected: ds{FieldUI64: 0}},
		{name: "struct type - uint64 - min out of range", q: "fieldUI64=-1", v: dsP(), err: fmt.Errorf(`"-1" can not be assign to uint64`)},
		{name: "struct type - uint64 - max", q: "fieldUI64=18446744073709551615", v: dsP(), expected: ds{FieldUI64: 18446744073709551615}},
		{name: "struct type - uint64 - max out of range", q: "fieldUI64=18446744073709551616", v: dsP(), err: fmt.Errorf(`"18446744073709551616" can not be assign to uint64`)},
		{name: "struct type - uint64 - not uint64", q: "fieldUI64=a", v: dsP(), err: fmt.Errorf(`"a" can not be assign to uint64`)},
		{name: "struct type - uint64 pointer - min", q: "fieldUI64=0", v: dspP(), expected: dsp{FieldUI64: uint64P(0)}},
		{name: "struct type - uint64 pointer - min out of range", q: "fieldUI64=-1", v: dspP(), err: fmt.Errorf(`"-1" can not be assign to *uint64`)},
		{name: "struct type - uint64 pointer - max", q: "fieldUI64=18446744073709551615", v: dspP(), expected: dsp{FieldUI64: uint64P(18446744073709551615)}},
		{name: "struct type - uint64 pointer - max out of range", q: "fieldUI64=18446744073709551616", v: dspP(), err: fmt.Errorf(`"18446744073709551616" can not be assign to *uint64`)},
		{name: "struct type - uint64 pointer - not uint64", q: "fieldUI64=a", v: dspP(), err: fmt.Errorf(`"a" can not be assign to *uint64`)},

		{name: "struct type - string", q: "json_str=1", v: dsP(), expected: ds{JSONStr: "1"}},
		{name: "struct type - string pointer", q: "json_str=1", v: dspP(), expected: dsp{JSONStr: stringP("1")}},

		{name: "struct type - string array", q: "array[0]=1&array[1]=a&array[2]=true", v: dsP(), expected: ds{Array: [3]string{"1", "a", "true"}}},
		{name: "struct type - string array - out of range", q: "array[0]=1&array[1]=a&array[2]=true&array[3]=b", v: dsP(), err: fmt.Errorf("index out of range [3] with [3]string")},
		{name: "struct type - int array", q: "arrayI[0]=1&arrayI[1]=a&arrayI[2]=true", v: dsP(), err: fmt.Errorf("[3]int is not supported")},
		{name: "struct type - string array pointer", q: "array[0]=1&array[1]=a&array[2]=true", v: dspP(), expected: dsp{Array: &[3]string{"1", "a", "true"}}},
		{name: "struct type - string array pointer - out of range", q: "array[0]=1&array[1]=a&array[2]=true&array[3]=b", v: dspP(), err: fmt.Errorf("index out of range [3] with [3]string")},
		{name: "struct type - int array pointer", q: "arrayI[0]=1&arrayI[1]=a&arrayI[2]=true", v: dspP(), err: fmt.Errorf("[3]int is not supported")},

		{name: "struct type - string slice", q: "slice[0]=1&slice[1]=a&slice[2]=true", v: dsP(), expected: ds{Slice: []string{"1", "a", "true"}}},
		{name: "struct type - nested string slice", q: "sliceNest[0][0]=1&sliceNest[0][1]=a&sliceNest[1][0]=true", v: dsP(), expected: ds{SliceNest: [][]string{{"1", "a"}, {"true"}}}},

		{name: "struct type - int slice", q: "sliceI[0]=1&sliceI[1]=a&sliceI[2]=true", v: dsP(), err: fmt.Errorf("[]int is not supported")},
		{name: "struct type - string slice pointer", q: "slice[0]=1&slice[1]=a&slice[2]=true", v: dspP(), expected: dsp{Slice: &[]string{"1", "a", "true"}}},
		{name: "struct type - int slice pointer", q: "sliceI[0]=1&sliceI[1]=a&sliceI[2]=true", v: dspP(), err: fmt.Errorf("[]int is not supported")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := strings.ReplaceAll(tc.q, "%5B", "[")
			q = strings.ReplaceAll(q, "%5D", "]")
			err := qstring.Decode(q, tc.v)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("Decode() should not returns error, got %q", err)
				}
				if err.Error() != tc.err.Error() {
					t.Fatalf("Decode() error returns %q, want %q", err, tc.err)
				}
			}

			if err == nil && tc.err != nil {
				t.Errorf("Decode() should returns error, want %q", tc.err)
			}
			if err == nil && tc.err == nil {
				if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(tc.v)).Interface(), tc.expected) {
					t.Errorf("Decode() returns \n%#v\nwant \n%#v", tc.v, tc.expected)
				}
			}
		})
	}
}

// helpers
func boolP(v bool) *bool          { return &v }
func intP(v int) *int             { return &v }
func int8P(v int8) *int8          { return &v }
func int16P(v int16) *int16       { return &v }
func int32P(v int32) *int32       { return &v }
func int64P(v int64) *int64       { return &v }
func uint8P(v uint8) *uint8       { return &v }
func uintP(v uint) *uint          { return &v }
func uint16P(v uint16) *uint16    { return &v }
func uint32P(v uint32) *uint32    { return &v }
func uint64P(v uint64) *uint64    { return &v }
func float64P(v float64) *float64 { return &v }
func float32P(v float32) *float32 { return &v }
func stringP(v string) *string    { return &v }
func dsP() *ds                    { var a ds; return &a }
func dspP() *dsp                  { var a dsp; return &a }
