package qstring_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/masakurapa/qstring"
)

type decodeCase struct {
	name     string
	q        string
	v        interface{}
	err      error
	expected interface{}
}

func TestDecode(t *testing.T) {
	t.Run("unsupported type", func(t *testing.T) {
		q := "test=1"
		runDecodeTest(t, []decodeCase{
			{name: "nil", q: q, v: nil, err: fmt.Errorf("nil is not supported")},
			{name: "not pointer", q: q, v: func() string { return "" }(), err: fmt.Errorf("non-pointer is not supported")},

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
			{name: "uintptr type", q: q, v: func() *uintptr { var a uintptr; return &a }(), err: fmt.Errorf("uintptr is not supported")},
			{name: "complex64 type", q: q, v: func() *complex64 { var a complex64; return &a }(), err: fmt.Errorf("complex64 is not supported")},
			{name: "complex128 type", q: q, v: func() *complex128 { var a complex128; return &a }(), err: fmt.Errorf("complex128 is not supported")},
			// {name: "chan type", q: q, v: func() *chan { var a chan; return &a }(), err: fmt.Errorf("chan is not supported")},
			{name: "func type", q: q, v: func() *string { return nil }, err: fmt.Errorf("non-pointer is not supported")},
			{name: "nil ptr type", q: q, v: func() *bool { return nil }(), err: fmt.Errorf("nil ptr is not supported")},
			{name: "unsafe pointer type", q: q, v: func() *unsafe.Pointer { var a unsafe.Pointer; return &a }(), err: fmt.Errorf("unsafe.Pointer is not supported")},
		})
	})

	t.Run("string", func(t *testing.T) {
		runDecodeTest(t, []decodeCase{
			{name: "with quote", q: "?hoge[key]=fuga", v: stringP(""), expected: "?hoge[key]=fuga"},
			{name: "without quote", q: "hoge[key]=fuga", v: stringP(""), expected: "hoge[key]=fuga"},
		})
	})

	t.Run("array", func(t *testing.T) {
		runDecodeTest(t, []decodeCase{
			{name: "string type", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *[3]string { var a [3]string; return &a }(), expected: [3]string{"a", "2", "3"}},
			{name: "string type - capacity exceeded", q: "hoge[]=a&hoge[]=2&hoge[]=3&hoge[]=4", v: func() *[3]string { var a [3]string; return &a }(), err: fmt.Errorf("index out of range [4] with [3]string")},
			{name: "string type - multiple key name", q: "hoge[]=a&fuga[]=2", v: func() *[3]string { var a [3]string; return &a }(), expected: [3]string{}}, //TODO error?
			{name: "int type", q: "hoge[]=1", v: func() *[3]int { var a [3]int; return &a }(), err: fmt.Errorf("[3]int is not supported")},
		})
	})

	t.Run("slice", func(t *testing.T) {
		runDecodeTest(t, []decodeCase{
			{name: "string type", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *[]string { var a []string; return &a }(), expected: []string{"a", "2", "3"}},
			{name: "int type", q: "hoge[]=1", v: func() *[]int { var a []int; return &a }(), err: fmt.Errorf("[]int is not supported")},
			{name: "string type - multiple key name", q: "hoge[]=a&fuga[]=2", v: func() *[]string { var a []string; return &a }(), expected: func() []string { var a []string; return a }()}, //TODO error?
		})
	})

	t.Run("map", func(t *testing.T) {
		runDecodeTest(t, []decodeCase{
			{name: "single key name", q: "hoge=a", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": "a"}},
			{name: "multiple key name", q: "hoge=a&fuga=1", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": "a", "fuga": "1"}},
			{name: "duplicate key name", q: "hoge=a&hoge=1", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": []string{"a", "1"}}},
			{name: "array key", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": []string{"a", "2", "3"}}},
			{name: "index array key", q: "hoge[0]=a&hoge[1]=2&hoge[2]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": []string{"a", "2", "3"}}},
			{name: "nested array and array", q: "hoge[0][0]=a&hoge[0][1]=2&hoge[0][2]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": qstring.ArrayQ{[]string{"a", "2", "3"}}}},
			{name: "nested map and array", q: "hoge[a][0]=a&hoge[a][1]=2&hoge[a][2]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": qstring.Q{"a": []string{"a", "2", "3"}}}},
			{name: "nested map and map", q: "hoge[a][b]=a&hoge[a][1]=2&hoge[0][a]=3", v: func() *qstring.Q { var a qstring.Q; return &a }(), expected: qstring.Q{"hoge": qstring.Q{"0": qstring.Q{"a": "3"}, "a": qstring.Q{"1": "2", "b": "a"}}}},
		})
	})

	t.Run("struct", func(t *testing.T) {
		t.Run("bool", func(t *testing.T) {
			type s struct {
				Field bool `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "false", q: "field=false", v: &s{}, expected: s{Field: false}},
				{name: "0 to false", q: "field=0", v: &s{}, expected: s{Field: false}},
				{name: "true", q: "field=true", v: &s{}, expected: s{Field: true}},
				{name: "1 to true", q: "field=1", v: &s{}, expected: s{Field: true}},
				{name: "not bool", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to bool`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: false}},
			})
		})
		t.Run("bool pointer", func(t *testing.T) {
			type s struct {
				Field *bool `qstring:"field"`
			}

			runDecodeTest(t, []decodeCase{
				{name: "false", q: "field=false", v: &s{}, expected: s{Field: boolP(false)}},
				{name: "0 to false", q: "field=0", v: &s{}, expected: s{Field: boolP(false)}},
				{name: "true", q: "field=true", v: &s{}, expected: s{Field: boolP(true)}},
				{name: "1 to true", q: "field=1", v: &s{}, expected: s{Field: boolP(true)}},
				{name: "not bool", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *bool`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("int", func(t *testing.T) {
			type s struct {
				Field int `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-9223372036854775808", v: &s{}, expected: s{Field: -9223372036854775808}},
				{name: "min out of range", q: "field=-9223372036854775809", v: &s{}, err: fmt.Errorf(`"-9223372036854775809" can not be assign to int`)},
				{name: "max", q: "field=9223372036854775807", v: &s{}, expected: s{Field: 9223372036854775807}},
				{name: "max out of range", q: "field=9223372036854775808", v: &s{}, err: fmt.Errorf(`"9223372036854775808" can not be assign to int`)},
				{name: "not int", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to int`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("int pointer", func(t *testing.T) {
			type s struct {
				Field *int `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-9223372036854775808", v: &s{}, expected: s{Field: intP(-9223372036854775808)}},
				{name: "min out of range", q: "field=-9223372036854775809", v: &s{}, err: fmt.Errorf(`"-9223372036854775809" can not be assign to *int`)},
				{name: "max", q: "field=9223372036854775807", v: &s{}, expected: s{Field: intP(9223372036854775807)}},
				{name: "max out of range", q: "field=9223372036854775808", v: &s{}, err: fmt.Errorf(`"9223372036854775808" can not be assign to *int`)},
				{name: "not int", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *int`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("int8", func(t *testing.T) {
			type s struct {
				Field int8 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-128", v: &s{}, expected: s{Field: -128}},
				{name: "min out of range", q: "field=-129", v: &s{}, err: fmt.Errorf(`"-129" can not be assign to int8`)},
				{name: "max", q: "field=127", v: &s{}, expected: s{Field: 127}},
				{name: "max out of range", q: "field=128", v: &s{}, err: fmt.Errorf(`"128" can not be assign to int8`)},
				{name: "not int8", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to int8`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("int8 pointer", func(t *testing.T) {
			type s struct {
				Field *int8 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-128", v: &s{}, expected: s{Field: int8P(-128)}},
				{name: "min out of range", q: "field=-129", v: &s{}, err: fmt.Errorf(`"-129" can not be assign to *int8`)},
				{name: "max", q: "field=127", v: &s{}, expected: s{Field: int8P(127)}},
				{name: "max out of range", q: "field=128", v: &s{}, err: fmt.Errorf(`"128" can not be assign to *int8`)},
				{name: "not int8", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *int8`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("int16", func(t *testing.T) {
			type s struct {
				Field int16 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-32768", v: &s{}, expected: s{Field: -32768}},
				{name: "min out of range", q: "field=-32769", v: &s{}, err: fmt.Errorf(`"-32769" can not be assign to int16`)},
				{name: "max", q: "field=32767", v: &s{}, expected: s{Field: 32767}},
				{name: "max out of range", q: "field=32768", v: &s{}, err: fmt.Errorf(`"32768" can not be assign to int16`)},
				{name: "not int16", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to int16`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("int16 pointer", func(t *testing.T) {
			type s struct {
				Field *int16 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-32768", v: &s{}, expected: s{Field: int16P(-32768)}},
				{name: "min out of range", q: "field=-32769", v: &s{}, err: fmt.Errorf(`"-32769" can not be assign to *int16`)},
				{name: "max", q: "field=32767", v: &s{}, expected: s{Field: int16P(32767)}},
				{name: "max out of range", q: "field=32768", v: &s{}, err: fmt.Errorf(`"32768" can not be assign to *int16`)},
				{name: "not int8", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *int16`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("int32", func(t *testing.T) {
			type s struct {
				Field int32 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-2147483648", v: &s{}, expected: s{Field: -2147483648}},
				{name: "min out of range", q: "field=-2147483649", v: &s{}, err: fmt.Errorf(`"-2147483649" can not be assign to int32`)},
				{name: "max", q: "field=2147483647", v: &s{}, expected: s{Field: 2147483647}},
				{name: "max out of range", q: "field=2147483648", v: &s{}, err: fmt.Errorf(`"2147483648" can not be assign to int32`)},
				{name: "not int32", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to int32`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("int32 pointer", func(t *testing.T) {
			type s struct {
				Field *int32 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-2147483648", v: &s{}, expected: s{Field: int32P(-2147483648)}},
				{name: "min out of range", q: "field=-2147483649", v: &s{}, err: fmt.Errorf(`"-2147483649" can not be assign to *int32`)},
				{name: "max", q: "field=2147483647", v: &s{}, expected: s{Field: int32P(2147483647)}},
				{name: "max out of range", q: "field=2147483648", v: &s{}, err: fmt.Errorf(`"2147483648" can not be assign to *int32`)},
				{name: "not int32", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *int32`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("int64", func(t *testing.T) {
			type s struct {
				Field int64 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-9223372036854775808", v: &s{}, expected: s{Field: -9223372036854775808}},
				{name: "min out of range", q: "field=-9223372036854775809", v: &s{}, err: fmt.Errorf(`"-9223372036854775809" can not be assign to int64`)},
				{name: "max", q: "field=9223372036854775807", v: &s{}, expected: s{Field: 9223372036854775807}},
				{name: "max out of range", q: "field=9223372036854775808", v: &s{}, err: fmt.Errorf(`"9223372036854775808" can not be assign to int64`)},
				{name: "not int64", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to int64`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("int64 pointer", func(t *testing.T) {
			type s struct {
				Field *int64 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=-9223372036854775808", v: &s{}, expected: s{Field: int64P(-9223372036854775808)}},
				{name: "min out of range", q: "field=-9223372036854775809", v: &s{}, err: fmt.Errorf(`"-9223372036854775809" can not be assign to *int64`)},
				{name: "max", q: "field=9223372036854775807", v: &s{}, expected: s{Field: int64P(9223372036854775807)}},
				{name: "max out of range", q: "field=9223372036854775808", v: &s{}, err: fmt.Errorf(`"9223372036854775808" can not be assign to *int64`)},
				{name: "not int64", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *int64`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("uint", func(t *testing.T) {
			type s struct {
				Field uint `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: 0}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to uint`)},
				{name: "max", q: "field=18446744073709551615", v: &s{}, expected: s{Field: 18446744073709551615}},
				{name: "max out of range", q: "field=18446744073709551616", v: &s{}, err: fmt.Errorf(`"18446744073709551616" can not be assign to uint`)},
				{name: "not uint", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to uint`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("uint pointer", func(t *testing.T) {
			type s struct {
				Field *uint `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: uintP(0)}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to *uint`)},
				{name: "max", q: "field=18446744073709551615", v: &s{}, expected: s{Field: uintP(18446744073709551615)}},
				{name: "max out of range", q: "field=18446744073709551616", v: &s{}, err: fmt.Errorf(`"18446744073709551616" can not be assign to *uint`)},
				{name: "not uint", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *uint`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("uint8", func(t *testing.T) {
			type s struct {
				Field uint8 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: 0}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to uint8`)},
				{name: "max", q: "field=128", v: &s{}, expected: s{Field: 128}},
				{name: "max out of range", q: "field=256", v: &s{}, err: fmt.Errorf(`"256" can not be assign to uint8`)},
				{name: "not uint8", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to uint8`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("uint8 pointer", func(t *testing.T) {
			type s struct {
				Field *uint8 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: uint8P(0)}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to *uint8`)},
				{name: "max", q: "field=128", v: &s{}, expected: s{Field: uint8P(128)}},
				{name: "max out of range", q: "field=256", v: &s{}, err: fmt.Errorf(`"256" can not be assign to *uint8`)},
				{name: "not uint8", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *uint8`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("uint16", func(t *testing.T) {
			type s struct {
				Field uint16 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: 0}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to uint16`)},
				{name: "max", q: "field=65535", v: &s{}, expected: s{Field: 65535}},
				{name: "max out of range", q: "field=65536", v: &s{}, err: fmt.Errorf(`"65536" can not be assign to uint16`)},
				{name: "not uint16", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to uint16`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("uint16 pointer", func(t *testing.T) {
			type s struct {
				Field *uint16 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: uint16P(0)}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to *uint16`)},
				{name: "max", q: "field=65535", v: &s{}, expected: s{Field: uint16P(65535)}},
				{name: "max out of range", q: "field=65536", v: &s{}, err: fmt.Errorf(`"65536" can not be assign to *uint16`)},
				{name: "not uint16", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *uint16`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("uint32", func(t *testing.T) {
			type s struct {
				Field uint32 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: 0}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to uint32`)},
				{name: "max", q: "field=4294967295", v: &s{}, expected: s{Field: 4294967295}},
				{name: "max out of range", q: "field=4294967296", v: &s{}, err: fmt.Errorf(`"4294967296" can not be assign to uint32`)},
				{name: "not uint32", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to uint32`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("uint32 pointer", func(t *testing.T) {
			type s struct {
				Field *uint32 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: uint32P(0)}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to *uint32`)},
				{name: "max", q: "field=4294967295", v: &s{}, expected: s{Field: uint32P(4294967295)}},
				{name: "max out of range", q: "field=4294967296", v: &s{}, err: fmt.Errorf(`"4294967296" can not be assign to *uint32`)},
				{name: "not uint32", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *uint32`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("uint64", func(t *testing.T) {
			type s struct {
				Field uint64 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: 0}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to uint64`)},
				{name: "max", q: "field=18446744073709551615", v: &s{}, expected: s{Field: 18446744073709551615}},
				{name: "max out of range", q: "field=18446744073709551616", v: &s{}, err: fmt.Errorf(`"18446744073709551616" can not be assign to uint64`)},
				{name: "not uint64", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to uint64`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: 0}},
			})
		})
		t.Run("uint64 pointer", func(t *testing.T) {
			type s struct {
				Field *uint64 `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "min", q: "field=0", v: &s{}, expected: s{Field: uint64P(0)}},
				{name: "min out of range", q: "field=-1", v: &s{}, err: fmt.Errorf(`"-1" can not be assign to *uint64`)},
				{name: "max", q: "field=18446744073709551615", v: &s{}, expected: s{Field: uint64P(18446744073709551615)}},
				{name: "max out of range", q: "field=18446744073709551616", v: &s{}, err: fmt.Errorf(`"18446744073709551616" can not be assign to *uint64`)},
				{name: "not uint64", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *uint64`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("string", func(t *testing.T) {
			type s struct {
				Field string `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "has field", q: "field=1", v: &s{}, expected: s{Field: "1"}},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: ""}},
			})
		})
		t.Run("string pointer", func(t *testing.T) {
			type s struct {
				Field *string `qstring:"field"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "has field", q: "field=1", v: &s{}, expected: s{Field: stringP("1")}},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("array", func(t *testing.T) {
			type s struct {
				Field  [3]string `qstring:"field"`
				FieldI [3]int    `qstring:"field_i"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "no index", q: "field[]=1&field[]=a&field[]=true", v: &s{}, expected: s{Field: [3]string{"1", "a", "true"}}},
				{name: "has index", q: "field[0]=1&field[1]=a&field[2]=true", v: &s{}, expected: s{Field: [3]string{"1", "a", "true"}}},
				{name: "out of range", q: "field[0]=1&field[1]=a&field[2]=true&field[3]=b", v: &s{}, err: fmt.Errorf("index out of range [3] with [3]string")},
				{name: "int", q: "field_i[0]=1&field_i[1]=a&field_i[2]=true", v: &s{}, err: fmt.Errorf("[3]int is not supported")},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: [3]string{"", "", ""}}},
			})
		})
		t.Run("array pointer", func(t *testing.T) {
			type s struct {
				Field  *[3]string `qstring:"field"`
				FieldI *[3]int    `qstring:"field_i"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "no index", q: "field[]=1&field[]=a&field[]=true", v: &s{}, expected: s{Field: &[3]string{"1", "a", "true"}}},
				{name: "has index", q: "field[0]=1&field[1]=a&field[2]=true", v: &s{}, expected: s{Field: &[3]string{"1", "a", "true"}}},
				{name: "out of range", q: "field[0]=1&field[1]=a&field[2]=true&field[3]=b", v: &s{}, err: fmt.Errorf("index out of range [3] with [3]string")},
				{name: "int", q: "field_i[0]=1&field_i[1]=a&field_i[2]=true", v: &s{}, err: fmt.Errorf("[3]int is not supported")},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("nested array", func(t *testing.T) {
			type s struct {
				Field  [3][2]string `qstring:"field"`
				FieldI [3][2]int    `qstring:"field_i"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "no index", q: "field[][]=1&field[][]=a", v: &s{}, expected: s{Field: [3][2]string{{"1", "a"}}}},
				{name: "string", q: "field[0][0]=1&field[0][1]=a&field[1][0]=true", v: &s{}, expected: s{Field: [3][2]string{{"1", "a"}, {"true", ""}, {"", ""}}}},
				{name: "no index and out of range", q: "field[][]=1&field[][]=a&field[][]=true", v: &s{}, err: fmt.Errorf("index out of range [2] with [2]string")},
				{name: "has index and out of range", q: "field[0][0]=1&field[1][0]=a&field[2][0]=true&field[3][0]=b", v: &s{}, err: fmt.Errorf("index out of range [3] with [3][2]string")},
				{name: "child out of range", q: "field[0][0]=1&field[0][1]=a&field[0][2]=true", v: &s{}, err: fmt.Errorf("index out of range [2] with [2]string")},
				{name: "int", q: "field_i[0]=1&field_i[1]=a&field_i[2]=true", v: &s{}, err: fmt.Errorf("[3][2]int is not supported")},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: [3][2]string{{"", ""}, {"", ""}, {"", ""}}}},
			})
		})
		t.Run("nested array pointer", func(t *testing.T) {
			type s struct {
				Field  *[3][2]string `qstring:"field"`
				FieldI *[3][2]int    `qstring:"field_i"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "no index", q: "field[][]=1&field[][]=a", v: &s{}, expected: s{Field: &[3][2]string{{"1", "a"}}}},
				{name: "string", q: "field[0][0]=1&field[0][1]=a&field[1][0]=true", v: &s{}, expected: s{Field: &[3][2]string{{"1", "a"}, {"true", ""}, {"", ""}}}},
				{name: "no index and out of range", q: "field[][]=1&field[][]=a&field[][]=true", v: &s{}, err: fmt.Errorf("index out of range [2] with [2]string")},
				{name: "has index and out of range", q: "field[0][0]=1&field[1][0]=a&field[2][0]=true&field[3][0]=b", v: &s{}, err: fmt.Errorf("index out of range [3] with [3][2]string")},
				{name: "child out of range", q: "field[0][0]=1&field[0][1]=a&field[0][2]=true", v: &s{}, err: fmt.Errorf("index out of range [2] with [2]string")},
				{name: "int", q: "field_i[0]=1&field_i[1]=a&field_i[2]=true", v: &s{}, err: fmt.Errorf("[3][2]int is not supported")},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("slice", func(t *testing.T) {
			type s struct {
				Field  []string `qstring:"field"`
				FieldI []int    `qstring:"field_i"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "no index", q: "field[]=1&field[]=a&field[]=true", v: &s{}, expected: s{Field: []string{"1", "a", "true"}}},
				{name: "has index", q: "field[0]=1&field[1]=a&field[2]=true", v: &s{}, expected: s{Field: []string{"1", "a", "true"}}},
				{name: "int", q: "field_i[0]=1&field_i[1]=a&field_i[2]=true", v: &s{}, err: fmt.Errorf("[]int is not supported")},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})
		t.Run("slice pointer", func(t *testing.T) {
			type s struct {
				Field  *[]string `qstring:"field"`
				FieldI []int     `qstring:"field_i"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "no index", q: "field[]=1&field[]=a&field[]=true", v: &s{}, expected: s{Field: &[]string{"1", "a", "true"}}},
				{name: "has index", q: "field[0]=1&field[1]=a&field[2]=true", v: &s{}, expected: s{Field: &[]string{"1", "a", "true"}}},
				{name: "int", q: "field_i[0]=1&field_i[1]=a&field_i[2]=true", v: &s{}, err: fmt.Errorf("[]int is not supported")},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("nested slice", func(t *testing.T) {
			type s struct {
				Field  [][]string `qstring:"field"`
				FieldI [][]int    `qstring:"field_i"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "no index", q: "field[][]=1&field[][]=a&field[][]=true", v: &s{}, expected: s{Field: [][]string{{"1", "a", "true"}}}},
				{name: "has index", q: "field[0][0]=1&field[0][1]=a&field[1][0]=true", v: &s{}, expected: s{Field: [][]string{{"1", "a"}, {"true"}}}},
				{name: "int", q: "field_i[0]=1&field_i[1]=a&field_i[2]=true", v: &s{}, err: fmt.Errorf("[][]int is not supported")},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})
		t.Run("nested slice pointer", func(t *testing.T) {
			type s struct {
				Field  *[][]string `qstring:"field"`
				FieldI *[][]int    `qstring:"field_i"`
			}
			runDecodeTest(t, []decodeCase{
				{name: "no index", q: "field[][]=1&field[][]=a&field[][]=true", v: &s{}, expected: s{Field: &[][]string{{"1", "a", "true"}}}},
				{name: "has index", q: "field[0][0]=1&field[0][1]=a&field[1][0]=true", v: &s{}, expected: s{Field: &[][]string{{"1", "a"}, {"true"}}}},
				{name: "int", q: "field_i[0]=1&field_i[1]=a&field_i[2]=true", v: &s{}, err: fmt.Errorf("[][]int is not supported")},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("struct", func(t *testing.T) {
			type c struct {
				Field  string `qstring:"child_field"`
				FieldI int    `qstring:"child_field_i"`
			}
			type s struct {
				Field c `qstring:"field"`
			}

			runDecodeTest(t, []decodeCase{
				{name: "has child field", q: "field[child_field]=a&field[child_field_i]=1", v: &s{}, expected: s{Field: c{Field: "a", FieldI: 1}}},
				{name: "no child field", q: "field[no]=1", v: &s{}, expected: s{Field: c{Field: "", FieldI: 0}}},
				{name: "not assign value", q: "field[child_field]=a&field[child_field_i]=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to int`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: c{Field: "", FieldI: 0}}},
			})
		})
		t.Run("struct pointer", func(t *testing.T) {
			type c struct {
				Field  *string `qstring:"child_field"`
				FieldI *int    `qstring:"child_field_i"`
			}
			type s struct {
				Field *c `qstring:"field"`
			}

			runDecodeTest(t, []decodeCase{
				{name: "has child field and nil child", q: "field[child_field]=a&field[child_field_i]=1", v: &s{}, expected: s{Field: &c{Field: stringP("a"), FieldI: intP(1)}}},
				{name: "has child field and non-nil child", q: "field[child_field]=a&field[child_field_i]=1", v: &s{Field: &c{}}, expected: s{Field: &c{Field: stringP("a"), FieldI: intP(1)}}},
				{name: "no child field", q: "field[no]=1", v: &s{}, expected: s{Field: &c{Field: nil, FieldI: nil}}},
				{name: "not assign value", q: "field[child_field]=a&field[child_field_i]=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to *int`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})

		t.Run("map", func(t *testing.T) {
			type s struct {
				Field qstring.Q `qstring:"field"`
			}

			runDecodeTest(t, []decodeCase{
				{name: "nil", q: "field[a]=1&field[1]=b&field[c]=true", v: &s{Field: nil}, expected: s{Field: qstring.Q{"a": "1", "1": "b", "c": "true"}}},
				{name: "non-nil", q: "field[a]=1&field[1]=b&field[c]=true", v: &s{}, expected: s{Field: qstring.Q{"a": "1", "1": "b", "c": "true"}}},
				{name: "not assign value", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to qstring.Q`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})
		t.Run("map pointer", func(t *testing.T) {
			type s struct {
				Field *qstring.Q `qstring:"field"`
			}

			runDecodeTest(t, []decodeCase{
				{name: "nil", q: "field[a]=1&field[1]=b&field[c]=true", v: &s{}, expected: s{Field: &qstring.Q{"a": "1", "1": "b", "c": "true"}}},
				{name: "non-nil", q: "field[a]=1&field[1]=b&field[c]=true", v: &s{Field: &qstring.Q{}}, expected: s{Field: &qstring.Q{"a": "1", "1": "b", "c": "true"}}},
				{name: "not assign value", q: "field=a", v: &s{}, err: fmt.Errorf(`"a" can not be assign to qstring.Q`)},
				{name: "no field", q: "no=1", v: &s{}, expected: s{Field: nil}},
			})
		})
	})
}

// helpers
func runDecodeTest(t *testing.T, testCases []decodeCase) {
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

func boolP(v bool) *bool              { return &v }
func intP(v int) *int                 { return &v }
func int8P(v int8) *int8              { return &v }
func int16P(v int16) *int16           { return &v }
func int32P(v int32) *int32           { return &v }
func int64P(v int64) *int64           { return &v }
func runeP(v rune) *rune              { return &v }
func uint8P(v uint8) *uint8           { return &v }
func uintP(v uint) *uint              { return &v }
func uint16P(v uint16) *uint16        { return &v }
func uint32P(v uint32) *uint32        { return &v }
func uint64P(v uint64) *uint64        { return &v }
func byteP(v byte) *byte              { return &v }
func float64P(v float64) *float64     { return &v }
func float32P(v float32) *float32     { return &v }
func stringP(v string) *string        { return &v }
func unsafeP(v string) unsafe.Pointer { return unsafe.Pointer(stringP(v)) }
