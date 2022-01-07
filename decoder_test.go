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
	FieldB       bool    `qstring:"field_b"`
	FieldI       int     `qstring:"fieldI"`
	FieldI8      int8    `qstring:"fieldI8"`
	FieldI16     int16   `qstring:"fieldI16"`
	FieldI32     int32   `qstring:"fieldI32"`
	FieldI64     int64   `qstring:"fieldI64"`
	FieldUI8     uint8   `qstring:"fieldUI8"`
	FieldUI      uint    `qstring:"fieldUI"`
	FieldUI16    uint16  `qstring:"fieldUI16"`
	FieldUI32    uint32  `qstring:"fieldUI32"`
	FieldUI64    uint64  `qstring:"fieldUI64"`
	FieldFloat32 float32 `qstring:"fieldFloat32"`
	FieldFloat64 float64 `qstring:"fieldFloat64"`
	JSONStr      string  `qstring:"json_str"`
}

type dsp struct {
	FieldB       *bool    `qstring:"field_b"`
	FieldI       *int     `qstring:"fieldI"`
	FieldI8      *int8    `qstring:"fieldI8"`
	FieldI16     *int16   `qstring:"fieldI16"`
	FieldI32     *int32   `qstring:"fieldI32"`
	FieldI64     *int64   `qstring:"fieldI64"`
	FieldUI8     *uint8   `qstring:"fieldUI8"`
	FieldUI      *uint    `qstring:"fieldUI"`
	FieldUI16    *uint16  `qstring:"fieldUI16"`
	FieldUI32    *uint32  `qstring:"fieldUI32"`
	FieldUI64    *uint64  `qstring:"fieldUI64"`
	FieldFloat32 *float32 `qstring:"fieldFloat32"`
	FieldFloat64 *float64 `qstring:"fieldFloat64"`
	JSONStr      *string  `qstring:"json_str"`
}

type ds2 struct {
	Field    string `qstring:"field"`
	NoTag    string
	privateS string `qstring:"private-S"`
}

func TestDecode(t *testing.T) {
	q := "test=1"

	testCases := []struct {
		name     string
		q        string
		v        interface{}
		err      error
		expected interface{}
	}{
		{name: "nil", q: q, v: nil, err: fmt.Errorf("nil is not available")},
		{name: "not pointer", q: q, v: func() string { return "" }(), err: fmt.Errorf("not pointer")},

		// unavailable types
		{name: "bool type", q: q, v: boolP(false), err: fmt.Errorf("type bool is not available")},
		{name: "int type", q: q, v: intP(0), err: fmt.Errorf("type int is not available")},
		{name: "int64 type", q: q, v: int64P(0), err: fmt.Errorf("type int64 is not available")},
		{name: "int32 type", q: q, v: int32P(0), err: fmt.Errorf("type int32 is not available")},
		{name: "int16 type", q: q, v: int16P(0), err: fmt.Errorf("type int16 is not available")},
		{name: "int8 type", q: q, v: int8P(0), err: fmt.Errorf("type int8 is not available")},
		{name: "uint type", q: q, v: uintP(0), err: fmt.Errorf("type uint is not available")},
		{name: "uint64 type", q: q, v: uint64P(0), err: fmt.Errorf("type uint64 is not available")},
		{name: "uint32 type", q: q, v: uint32P(0), err: fmt.Errorf("type uint32 is not available")},
		{name: "uint16 type", q: q, v: uint16P(0), err: fmt.Errorf("type uint16 is not available")},
		{name: "uint8 type", q: q, v: uint8P(0), err: fmt.Errorf("type uint8 is not available")},
		{name: "float64 type", q: q, v: float64P(0), err: fmt.Errorf("type float64 is not available")},
		{name: "float32 type", q: q, v: float32P(0), err: fmt.Errorf("type float32 is not available")},
		// {name: "array type", q: q, v: func() *bool { var a bool; return &a }(), err: fmt.Errorf("type array is not available")},
		// {name: "slice type", q: q, v: func() *bool { var a bool; return &a }(), err: fmt.Errorf("type slice is not available")},
		{name: "uintptr type", q: q, v: func() *uintptr { var a uintptr; return &a }(), err: fmt.Errorf("type uintptr is not available")},
		{name: "complex64 type", q: q, v: func() *complex64 { var a complex64; return &a }(), err: fmt.Errorf("type complex64 is not available")},
		{name: "complex128 type", q: q, v: func() *complex128 { var a complex128; return &a }(), err: fmt.Errorf("type complex128 is not available")},
		// {name: "chan type", q: q, v: func() *chan { var a chan; return &a }(), err: fmt.Errorf("type chan is not available")},
		{name: "func type", q: q, v: func() *string { return nil }, err: fmt.Errorf("not pointer")},
		{name: "nil ptr type", q: q, v: func() *bool { return nil }(), err: fmt.Errorf("nil is not available")},
		{name: "unsafe pointer type", q: q, v: func() *unsafe.Pointer { var a unsafe.Pointer; return &a }(), err: fmt.Errorf("type unsafe.Pointer is not available")},

		// string type
		{name: "string type with quote", q: "?hoge[key]=fuga", v: stringP(""), expected: "?hoge[key]=fuga"},
		{name: "string type without quote", q: "hoge[key]=fuga", v: stringP(""), expected: "hoge[key]=fuga"},

		// array type
		{name: "array of string type", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *[3]string { var a [3]string; return &a }(), expected: [3]string{"a", "2", "3"}},
		{name: "array of string type - capacity exceeded", q: "hoge[]=a&hoge[]=2&hoge[]=3&hoge[]=4", v: func() *[3]string { var a [3]string; return &a }(), err: fmt.Errorf("array capacity exceeded")},
		{name: "array of string type - multiple key name", q: "hoge[]=a&fuga[]=2", v: func() *[3]string { var a [3]string; return &a }(), expected: [3]string{}},
		{name: "array of int type", q: "hoge[]=1", v: func() *[3]int { var a [3]int; return &a }(), err: fmt.Errorf("allocation type must be [n]stirng")},

		// slice type
		{name: "slice of string type", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *[]string { var a []string; return &a }(), expected: []string{"a", "2", "3"}},
		{name: "slice of int type", q: "hoge[]=1", v: func() *[]int { var a []int; return &a }(), err: fmt.Errorf("allocation type must be []stirng")},
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
		{name: "struct type - bool - not bool", q: "field_b=a", v: dsP(), expected: ds{FieldB: false}},
		{name: "struct type - bool pointer", q: "field_b=1", v: dspP(), expected: dsp{FieldB: boolP(true)}},
		{name: "struct type - bool pointer - not bool", q: "field_b=a", v: dspP(), expected: dsp{FieldB: nil}},

		{name: "struct type - int - min", q: "fieldI=-9223372036854775808", v: dsP(), expected: ds{FieldI: -9223372036854775808}},
		{name: "struct type - int - min out of range", q: "fieldI=-9223372036854775809", v: dsP(), expected: ds{FieldI: 0}},
		{name: "struct type - int - max", q: "fieldI=9223372036854775807", v: dsP(), expected: ds{FieldI: 9223372036854775807}},
		{name: "struct type - int - max out of range", q: "fieldI=9223372036854775808", v: dsP(), expected: ds{FieldI: 0}},
		{name: "struct type - int - not int", q: "fieldI=a", v: dsP(), expected: ds{FieldI: 0}},
		{name: "struct type - int pointer - min", q: "fieldI=-9223372036854775808", v: dspP(), expected: dsp{FieldI: intP(-9223372036854775808)}},
		{name: "struct type - int pointer - min out of range", q: "fieldI=-9223372036854775809", v: dspP(), expected: dsp{FieldI: nil}},
		{name: "struct type - int pointer - max", q: "fieldI=9223372036854775807", v: dspP(), expected: dsp{FieldI: intP(9223372036854775807)}},
		{name: "struct type - int pointer - max out of range", q: "fieldI=9223372036854775808", v: dspP(), expected: dsp{FieldI: nil}},
		{name: "struct type - int pointer - not int", q: "fieldI=a", v: dspP(), expected: dsp{FieldI: nil}},

		{name: "struct type - int8 - min", q: "fieldI8=-128", v: dsP(), expected: ds{FieldI8: -128}},
		{name: "struct type - int8 - min out of range", q: "fieldI8=-129", v: dsP(), expected: ds{FieldI8: 0}},
		{name: "struct type - int8 - max", q: "fieldI8=127", v: dsP(), expected: ds{FieldI8: 127}},
		{name: "struct type - int8 - max out of range", q: "fieldI8=128", v: dsP(), expected: ds{FieldI8: 0}},
		{name: "struct type - int8 - not int8", q: "fieldI8=a", v: dsP(), expected: ds{FieldI8: 0}},
		{name: "struct type - int8 pointer - min", q: "fieldI8=-128", v: dspP(), expected: dsp{FieldI8: int8P(-128)}},
		{name: "struct type - int8 pointer - min out of range", q: "fieldI8=-129", v: dspP(), expected: dsp{FieldI8: nil}},
		{name: "struct type - int8 pointer - max", q: "fieldI8=127", v: dspP(), expected: dsp{FieldI8: int8P(127)}},
		{name: "struct type - int8 pointer - max out of range", q: "fieldI8=128", v: dspP(), expected: dsp{FieldI8: nil}},
		{name: "struct type - int8 pointer - not int8", q: "fieldI8=a", v: dspP(), expected: dsp{FieldI8: nil}},

		{name: "struct type - int16 - min", q: "fieldI16=-32768", v: dsP(), expected: ds{FieldI16: -32768}},
		{name: "struct type - int16 - min out of range", q: "fieldI16=-32769", v: dsP(), expected: ds{FieldI16: 0}},
		{name: "struct type - int16 - max", q: "fieldI16=32767", v: dsP(), expected: ds{FieldI16: 32767}},
		{name: "struct type - int16 - max out of range", q: "fieldI16=32768", v: dsP(), expected: ds{FieldI16: 0}},
		{name: "struct type - int16 - not int16", q: "fieldI16=a", v: dsP(), expected: ds{FieldI16: 0}},
		{name: "struct type - int16 pointer - min", q: "fieldI16=-32768", v: dspP(), expected: dsp{FieldI16: int16P(-32768)}},
		{name: "struct type - int16 pointer - min out of range", q: "fieldI16=-32769", v: dspP(), expected: dsp{FieldI16: nil}},
		{name: "struct type - int16 pointer - max", q: "fieldI16=32767", v: dspP(), expected: dsp{FieldI16: int16P(32767)}},
		{name: "struct type - int16 pointer - max out of range", q: "fieldI16=32768", v: dspP(), expected: dsp{FieldI16: nil}},
		{name: "struct type - int16 pointer - not int8", q: "fieldI16=a", v: dspP(), expected: dsp{FieldI16: nil}},

		{name: "struct type - int32 - min", q: "fieldI32=-2147483648", v: dsP(), expected: ds{FieldI32: -2147483648}},
		{name: "struct type - int32 - min out of range", q: "fieldI32=-2147483649", v: dsP(), expected: ds{FieldI32: 0}},
		{name: "struct type - int32 - max", q: "fieldI32=2147483647", v: dsP(), expected: ds{FieldI32: 2147483647}},
		{name: "struct type - int32 - max out of range", q: "fieldI32=2147483648", v: dsP(), expected: ds{FieldI32: 0}},
		{name: "struct type - int32 - not int32", q: "fieldI32=a", v: dsP(), expected: ds{FieldI32: 0}},
		{name: "struct type - int32 pointer - min", q: "fieldI32=-2147483648", v: dspP(), expected: dsp{FieldI32: int32P(-2147483648)}},
		{name: "struct type - int32 pointer - min out of range", q: "fieldI32=-2147483649", v: dspP(), expected: dsp{FieldI32: nil}},
		{name: "struct type - int32 pointer - max", q: "fieldI32=2147483647", v: dspP(), expected: dsp{FieldI32: int32P(2147483647)}},
		{name: "struct type - int32 pointer - max out of range", q: "fieldI32=2147483648", v: dspP(), expected: dsp{FieldI32: nil}},
		{name: "struct type - int32 pointer - not int32", q: "fieldI32=a", v: dspP(), expected: dsp{FieldI32: nil}},

		{name: "struct type - int64 - min", q: "fieldI64=-9223372036854775808", v: dsP(), expected: ds{FieldI64: -9223372036854775808}},
		{name: "struct type - int64 - min out of range", q: "fieldI64=-9223372036854775809", v: dsP(), expected: ds{FieldI64: 0}},
		{name: "struct type - int64 - max", q: "fieldI64=9223372036854775807", v: dsP(), expected: ds{FieldI64: 9223372036854775807}},
		{name: "struct type - int64 - max out of range", q: "fieldI64=9223372036854775808", v: dsP(), expected: ds{FieldI64: 0}},
		{name: "struct type - int64 - not int64", q: "fieldI64=a", v: dsP(), expected: ds{FieldI64: 0}},
		{name: "struct type - int64 pointer - min", q: "fieldI64=-9223372036854775808", v: dspP(), expected: dsp{FieldI64: int64P(-9223372036854775808)}},
		{name: "struct type - int64 pointer - min out of range", q: "fieldI64=-9223372036854775809", v: dspP(), expected: dsp{FieldI64: nil}},
		{name: "struct type - int64 pointer - max", q: "fieldI64=9223372036854775807", v: dspP(), expected: dsp{FieldI64: int64P(9223372036854775807)}},
		{name: "struct type - int64 pointer - max out of range", q: "fieldI64=9223372036854775808", v: dspP(), expected: dsp{FieldI64: nil}},
		{name: "struct type - int64 pointer - not int64", q: "fieldI64=a", v: dspP(), expected: dsp{FieldI64: nil}},

		{name: "struct type - uint - min", q: "fieldUI=0", v: dsP(), expected: ds{FieldUI: 0}},
		{name: "struct type - uint - min out of range", q: "fieldUI=-1", v: dsP(), expected: ds{FieldUI: 0}},
		{name: "struct type - uint - max", q: "fieldUI=18446744073709551615", v: dsP(), expected: ds{FieldUI: 18446744073709551615}},
		{name: "struct type - uint - max out of range", q: "fieldUI=18446744073709551616", v: dsP(), expected: ds{FieldUI: 0}},
		{name: "struct type - uint - not uint", q: "fieldUI=a", v: dsP(), expected: ds{FieldUI: 0}},
		{name: "struct type - uint pointer - min", q: "fieldUI=0", v: dspP(), expected: dsp{FieldUI: uintP(0)}},
		{name: "struct type - uint pointer - min out of range", q: "fieldUI=-1", v: dspP(), expected: dsp{FieldUI: nil}},
		{name: "struct type - uint pointer - max", q: "fieldUI=18446744073709551615", v: dspP(), expected: dsp{FieldUI: uintP(18446744073709551615)}},
		{name: "struct type - uint pointer - max out of range", q: "fieldUI=18446744073709551616", v: dspP(), expected: dsp{FieldUI: nil}},
		{name: "struct type - uint pointer - not uint", q: "fieldUI=a", v: dspP(), expected: dsp{FieldUI: nil}},

		{name: "struct type - uint8 - min", q: "fieldUI8=0", v: dsP(), expected: ds{FieldUI8: 0}},
		{name: "struct type - uint8 - min out of range", q: "fieldUI8=-1", v: dsP(), expected: ds{FieldUI8: 0}},
		{name: "struct type - uint8 - max", q: "fieldUI8=128", v: dsP(), expected: ds{FieldUI8: 128}},
		{name: "struct type - uint8 - max out of range", q: "fieldUI8=256", v: dsP(), expected: ds{FieldUI8: 0}},
		{name: "struct type - uint8 - not uint8", q: "fieldUI8=a", v: dsP(), expected: ds{FieldUI8: 0}},
		{name: "struct type - uint8 pointer - min", q: "fieldUI8=0", v: dspP(), expected: dsp{FieldUI8: uint8P(0)}},
		{name: "struct type - uint8 pointer - min out of range", q: "fieldUI8=-1", v: dspP(), expected: dsp{FieldUI8: nil}},
		{name: "struct type - uint8 pointer - max", q: "fieldUI8=128", v: dspP(), expected: dsp{FieldUI8: uint8P(128)}},
		{name: "struct type - uint8 pointer - max out of range", q: "fieldUI8=256", v: dspP(), expected: dsp{FieldUI8: nil}},
		{name: "struct type - uint8 pointer - not uint8", q: "fieldUI8=a", v: dspP(), expected: dsp{FieldUI8: nil}},

		{name: "struct type - uint16 - min", q: "fieldUI16=0", v: dsP(), expected: ds{FieldUI16: 0}},
		{name: "struct type - uint16 - min out of range", q: "fieldUI16=-1", v: dsP(), expected: ds{FieldUI16: 0}},
		{name: "struct type - uint16 - max", q: "fieldUI16=65535", v: dsP(), expected: ds{FieldUI16: 65535}},
		{name: "struct type - uint16 - max out of range", q: "fieldUI16=65536", v: dsP(), expected: ds{FieldUI16: 0}},
		{name: "struct type - uint16 - not uint16", q: "fieldUI16=a", v: dsP(), expected: ds{FieldUI16: 0}},
		{name: "struct type - uint16 pointer - min", q: "fieldUI16=0", v: dspP(), expected: dsp{FieldUI16: uint16P(0)}},
		{name: "struct type - uint16 pointer - min out of range", q: "fieldUI16=-1", v: dspP(), expected: dsp{FieldUI16: nil}},
		{name: "struct type - uint16 pointer - max", q: "fieldUI16=65535", v: dspP(), expected: dsp{FieldUI16: uint16P(65535)}},
		{name: "struct type - uint16 pointer - max out of range", q: "fieldUI16=65536", v: dspP(), expected: dsp{FieldUI16: nil}},
		{name: "struct type - uint16 pointer - not uint16", q: "fieldUI16=a", v: dspP(), expected: dsp{FieldUI16: nil}},

		{name: "struct type - uint32 - min", q: "fieldUI32=0", v: dsP(), expected: ds{FieldUI32: 0}},
		{name: "struct type - uint32 - min out of range", q: "fieldUI32=-1", v: dsP(), expected: ds{FieldUI32: 0}},
		{name: "struct type - uint32 - max", q: "fieldUI32=4294967295", v: dsP(), expected: ds{FieldUI32: 4294967295}},
		{name: "struct type - uint32 - max out of range", q: "fieldUI32=4294967296", v: dsP(), expected: ds{FieldUI32: 0}},
		{name: "struct type - uint32 - not uint32", q: "fieldUI32=a", v: dsP(), expected: ds{FieldUI32: 0}},
		{name: "struct type - uint32 pointer - min", q: "fieldUI32=0", v: dspP(), expected: dsp{FieldUI32: uint32P(0)}},
		{name: "struct type - uint32 pointer - min out of range", q: "fieldUI32=-1", v: dspP(), expected: dsp{FieldUI32: nil}},
		{name: "struct type - uint32 pointer - max", q: "fieldUI32=4294967295", v: dspP(), expected: dsp{FieldUI32: uint32P(4294967295)}},
		{name: "struct type - uint32 pointer - max out of range", q: "fieldUI32=4294967296", v: dspP(), expected: dsp{FieldUI32: nil}},
		{name: "struct type - uint32 pointer - not uint32", q: "fieldUI32=a", v: dspP(), expected: dsp{FieldUI32: nil}},

		{name: "struct type - uint64 - min", q: "fieldUI64=0", v: dsP(), expected: ds{FieldUI64: 0}},
		{name: "struct type - uint64 - min out of range", q: "fieldUI64=-1", v: dsP(), expected: ds{FieldUI64: 0}},
		{name: "struct type - uint64 - max", q: "fieldUI64=18446744073709551615", v: dsP(), expected: ds{FieldUI64: 18446744073709551615}},
		{name: "struct type - uint64 - max out of range", q: "fieldUI64=18446744073709551616", v: dsP(), expected: ds{FieldUI64: 0}},
		{name: "struct type - uint64 - not uint64", q: "fieldUI64=a", v: dsP(), expected: ds{FieldUI64: 0}},
		{name: "struct type - uint64 pointer - min", q: "fieldUI64=0", v: dspP(), expected: dsp{FieldUI64: uint64P(0)}},
		{name: "struct type - uint64 pointer - min out of range", q: "fieldUI64=-1", v: dspP(), expected: dsp{FieldUI64: nil}},
		{name: "struct type - uint64 pointer - max", q: "fieldUI64=18446744073709551615", v: dspP(), expected: dsp{FieldUI64: uint64P(18446744073709551615)}},
		{name: "struct type - uint64 pointer - max out of range", q: "fieldUI64=18446744073709551616", v: dspP(), expected: dsp{FieldUI64: nil}},
		{name: "struct type - uint64 pointer - not uint64", q: "fieldUI64=a", v: dspP(), expected: dsp{FieldUI64: nil}},

		{name: "struct type - string", q: "json_str=1", v: dsP(), expected: ds{JSONStr: "1"}},
		{name: "struct type - string pointer", q: "json_str=1", v: dspP(), expected: dsp{JSONStr: stringP("1")}},
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
