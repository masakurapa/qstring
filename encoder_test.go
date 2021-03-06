package qstring_test

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/masakurapa/qstring"
)

type encodeCase struct {
	name     string
	q        interface{}
	err      error
	expected string
}

func TestEncode(t *testing.T) {
	t.Run("unsupported type", func(t *testing.T) {
		runEncodeTest(t, []encodeCase{
			{name: "nil", q: nil, err: fmt.Errorf("nil is not supported")},
			{name: "bool", q: true, err: fmt.Errorf("bool is not supported")},
			{name: "int", q: int(123), err: fmt.Errorf("int is not supported")},
			{name: "int64", q: int64(123), err: fmt.Errorf("int64 is not supported")},
			{name: "int32", q: int32(123), err: fmt.Errorf("int32 is not supported")},
			{name: "int16", q: int16(123), err: fmt.Errorf("int16 is not supported")},
			{name: "int8", q: int8(123), err: fmt.Errorf("int8 is not supported")},
			{name: "rune", q: rune(123), err: fmt.Errorf("int32 is not supported")},
			{name: "char", q: '1', err: fmt.Errorf("int32 is not supported")},
			{name: "uint", q: uint(123), err: fmt.Errorf("uint is not supported")},
			{name: "uint64", q: uint64(123), err: fmt.Errorf("uint64 is not supported")},
			{name: "uint32", q: uint32(123), err: fmt.Errorf("uint32 is not supported")},
			{name: "uint16", q: uint16(123), err: fmt.Errorf("uint16 is not supported")},
			{name: "uint8", q: uint8(123), err: fmt.Errorf("uint8 is not supported")},
			{name: "byte", q: byte(123), err: fmt.Errorf("uint8 is not supported")},
			{name: "float64", q: float64(123.456), err: fmt.Errorf("float64 is not supported")},
			{name: "float32", q: float32(123.456), err: fmt.Errorf("float32 is not supported")},
			{name: "complex128", q: complex128(123.456), err: fmt.Errorf("complex128 is not supported")},
			{name: "complex64", q: complex64(123.456), err: fmt.Errorf("complex64 is not supported")},
			{name: "array", q: [3]string{"1", "2", "3"}, err: fmt.Errorf("[3]string is not supported")},
			{name: "slice", q: []string{"1", "2", "3"}, err: fmt.Errorf("[]string is not supported")},
			{name: "uintptr", q: uintptr(1), err: fmt.Errorf("uintptr is not supported")},
			{name: "complex64", q: complex64(1), err: fmt.Errorf("complex64 is not supported")},
			{name: "complex128", q: complex128(1), err: fmt.Errorf("complex128 is not supported")},
			{name: "chan", q: make(chan int), err: fmt.Errorf("chan is not supported")},
			{name: "func", q: func() {}, err: fmt.Errorf("func is not supported")},
			{name: "ptr", q: uintP(123), err: fmt.Errorf("uint is not supported")},
			{name: "nil ptr", q: func() *string { return nil }(), err: fmt.Errorf("nil ptr is not supported")},
			{name: "unsafe pointer", q: unsafeP("1"), err: fmt.Errorf("unsafe.Pointer is not supported")},
		})
	})

	t.Run("string", func(t *testing.T) {
		runEncodeTest(t, []encodeCase{
			{name: "single value", q: "hoge=fuga", expected: "hoge=fuga"},
			{name: "single value and has quote", q: "?hoge=fuga", expected: "?hoge=fuga"},
			{name: "multiple value", q: "hoge=fuga&fuga=hoge", expected: "hoge=fuga&fuga=hoge"},
			{name: "array value", q: "hoge[]=fuga", expected: "hoge[]=fuga"},
			{name: "map value", q: "hoge[key]=fuga", expected: "hoge[key]=fuga"},
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("unsupported type value", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "uintptr", q: map[string]uintptr{"key": uintptr(1)}, err: fmt.Errorf("uintptr is not supported")},
				{name: "complex64", q: map[string]complex64{"key": complex64(1)}, err: fmt.Errorf("complex64 is not supported")},
				{name: "complex128", q: map[string]complex128{"key": complex128(1)}, err: fmt.Errorf("complex128 is not supported")},
				{name: "chan", q: qstring.Q{"key": make(chan int)}, err: fmt.Errorf("chan is not supported")},
				{name: "func", q: map[string]func(){"key": func() {}}, err: fmt.Errorf("func is not supported")},
				{name: "unsafe pointer", q: map[string]unsafe.Pointer{"key": unsafeP("1")}, err: fmt.Errorf("unsafe.Pointer is not supported")},
				{name: "key is not string", q: map[int]string{100: "val"}, err: fmt.Errorf("map[int]string is not supported")},
			})
		})

		t.Run("empty map", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "string type value", q: map[string]string{}, expected: ""},
				{name: "string type value pointer", q: &map[string]string{}, expected: ""},
				{name: "interface type value", q: qstring.Q{}, expected: ""},
				{name: "interface type value pointer", q: &qstring.Q{}, expected: ""},
			})
		})

		t.Run("supported type value", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "bool value", q: map[string]bool{"key": true}, expected: "key=true"},
				{name: "bool value pointer", q: map[string]*bool{"key": boolP(false)}, expected: "key=false"},

				{name: "int value", q: map[string]int{"key": int(123)}, expected: "key=123"},
				{name: "int value pointer", q: map[string]*int{"key": intP(123)}, expected: "key=123"},
				{name: "int64 value", q: map[string]int64{"key": int64(123)}, expected: "key=123"},
				{name: "int64 value pointer", q: map[string]*int64{"key": int64P(123)}, expected: "key=123"},
				{name: "int32 value", q: map[string]int32{"key": int32(123)}, expected: "key=123"},
				{name: "int32 value pointer", q: map[string]*int32{"key": int32P(123)}, expected: "key=123"},
				{name: "int16 value", q: map[string]int16{"key": int16(123)}, expected: "key=123"},
				{name: "int16 value pointer", q: map[string]*int16{"key": int16P(123)}, expected: "key=123"},
				{name: "int8 value", q: map[string]int8{"key": int8(123)}, expected: "key=123"},
				{name: "int8 value pointer", q: map[string]*int8{"key": int8P(123)}, expected: "key=123"},
				{name: "rune value", q: map[string]rune{"key": rune(123)}, expected: "key=123"},
				{name: "rune value pointer", q: map[string]*rune{"key": runeP(123)}, expected: "key=123"},
				{name: "char value", q: qstring.Q{"key": '1'}, expected: "key=49"},         // 1 is 49 in ascii
				{name: "char value pointer", q: qstring.Q{"key": '1'}, expected: "key=49"}, // 1 is 49 in ascii
				{name: "uint value", q: map[string]uint{"key": uint(123)}, expected: "key=123"},
				{name: "uint value pointer", q: map[string]*uint{"key": uintP(123)}, expected: "key=123"},
				{name: "uint64 value", q: map[string]uint64{"key": uint64(123)}, expected: "key=123"},
				{name: "uint64 value pointer", q: map[string]*uint64{"key": uint64P(123)}, expected: "key=123"},
				{name: "uint32 value", q: map[string]uint32{"key": uint32(123)}, expected: "key=123"},
				{name: "uint32 value pointer", q: map[string]*uint32{"key": uint32P(123)}, expected: "key=123"},
				{name: "uint16 value", q: map[string]uint16{"key": uint16(123)}, expected: "key=123"},
				{name: "uint16 value pointer", q: map[string]*uint16{"key": uint16P(123)}, expected: "key=123"},
				{name: "uint8 value", q: map[string]uint8{"key": uint8(123)}, expected: "key=123"},
				{name: "uint8 value pointer", q: map[string]*uint8{"key": uint8P(123)}, expected: "key=123"},
				{name: "byte value", q: map[string]byte{"key": byte(123)}, expected: "key=123"},
				{name: "byte value pointer", q: map[string]*byte{"key": byteP(123)}, expected: "key=123"},

				{name: "float64 value", q: map[string]float64{"key": float64(123.456)}, expected: "key=123.456"},
				{name: "float64 value pointer", q: map[string]*float64{"key": float64P(123.456)}, expected: "key=123.456"},
				{name: "float32 value", q: map[string]float32{"key": float32(123.456)}, expected: "key=123.456"},
				{name: "float32 value pointer", q: map[string]*float32{"key": float32P(123.456)}, expected: "key=123.456"},

				{name: "string value", q: map[string]string{"key": "hoge"}, expected: "key=hoge"},
				{name: "string value pointer", q: map[string]*string{"key": stringP("hoge")}, expected: "key=hoge"},
			})
		})

		t.Run("interface type value", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "bool value", q: qstring.Q{"key": true}, expected: "key=true"},
				{name: "bool value pointer", q: qstring.Q{"key": boolP(false)}, expected: "key=false"},

				{name: "int value", q: qstring.Q{"key": int(123)}, expected: "key=123"},
				{name: "int value pointer", q: qstring.Q{"key": intP(123)}, expected: "key=123"},
				{name: "int64 value", q: qstring.Q{"key": int64(123)}, expected: "key=123"},
				{name: "int64 value pointer", q: qstring.Q{"key": int64P(123)}, expected: "key=123"},
				{name: "int32 value", q: qstring.Q{"key": int32(123)}, expected: "key=123"},
				{name: "int32 value pointer", q: qstring.Q{"key": int32P(123)}, expected: "key=123"},
				{name: "int16 value", q: qstring.Q{"key": int16(123)}, expected: "key=123"},
				{name: "int16 value pointer", q: qstring.Q{"key": int16P(123)}, expected: "key=123"},
				{name: "int8 value", q: qstring.Q{"key": int8(123)}, expected: "key=123"},
				{name: "int8 value pointer", q: qstring.Q{"key": int8P(123)}, expected: "key=123"},
				{name: "rune value", q: qstring.Q{"key": rune(123)}, expected: "key=123"},
				{name: "rune value pointer", q: qstring.Q{"key": runeP(123)}, expected: "key=123"},
				{name: "char value", q: qstring.Q{"key": '1'}, expected: "key=49"},         // 1 is 49 in ascii
				{name: "char value pointer", q: qstring.Q{"key": '1'}, expected: "key=49"}, // 1 is 49 in ascii
				{name: "uint value", q: qstring.Q{"key": uint(123)}, expected: "key=123"},
				{name: "uint value pointer", q: qstring.Q{"key": uintP(123)}, expected: "key=123"},
				{name: "uint64 value", q: qstring.Q{"key": uint64(123)}, expected: "key=123"},
				{name: "uint64 value pointer", q: qstring.Q{"key": uint64P(123)}, expected: "key=123"},
				{name: "uint32 value", q: qstring.Q{"key": uint32(123)}, expected: "key=123"},
				{name: "uint32 value pointer", q: qstring.Q{"key": uint32P(123)}, expected: "key=123"},
				{name: "uint16 value", q: qstring.Q{"key": uint16(123)}, expected: "key=123"},
				{name: "uint16 value pointer", q: qstring.Q{"key": uint16P(123)}, expected: "key=123"},
				{name: "uint8 value", q: qstring.Q{"key": uint8(123)}, expected: "key=123"},
				{name: "uint8 value pointer", q: qstring.Q{"key": uint8P(123)}, expected: "key=123"},
				{name: "byte value", q: qstring.Q{"key": byte(123)}, expected: "key=123"},
				{name: "byte value pointer", q: qstring.Q{"key": byteP(123)}, expected: "key=123"},

				{name: "float64 value", q: qstring.Q{"key": float64(123.456)}, expected: "key=123.456"},
				{name: "float64 value pointer", q: qstring.Q{"key": float64P(123.456)}, expected: "key=123.456"},
				{name: "float32 value", q: qstring.Q{"key": float32(123.456)}, expected: "key=123.456"},
				{name: "float32 value pointer", q: qstring.Q{"key": float32P(123.456)}, expected: "key=123.456"},

				{name: "string value", q: qstring.Q{"key": "hoge"}, expected: "key=hoge"},
				{name: "string value pointer", q: qstring.Q{"key": stringP("hoge")}, expected: "key=hoge"},
			})
		})

		t.Run("array type value", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "string", q: qstring.Q{"key": [3]string{"1", "2", "3"}}, expected: "key[0]=1&key[1]=2&key[2]=3"},
				{name: "string pointer", q: qstring.Q{"key": &[3]string{"1", "2", "3"}}, expected: "key[0]=1&key[1]=2&key[2]=3"},
				{name: "interface", q: qstring.Q{"key": [3]interface{}{1, "2", true}}, expected: "key[0]=1&key[1]=2&key[2]=true"},
				{name: "interface pointer", q: qstring.Q{"key": &[3]interface{}{1, "2", true}}, expected: "key[0]=1&key[1]=2&key[2]=true"},
			})
		})

		t.Run("slice type value", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "string", q: qstring.Q{"key": []string{"1", "2", "3"}}, expected: "key[0]=1&key[1]=2&key[2]=3"},
				{name: "string pointer", q: qstring.Q{"key": &[]string{"1", "2", "3"}}, expected: "key[0]=1&key[1]=2&key[2]=3"},
				{name: "interface", q: qstring.Q{"key": qstring.S{1, "2", true}}, expected: "key[0]=1&key[1]=2&key[2]=true"},
				{name: "interface pointer", q: qstring.Q{"key": &qstring.S{1, "2", true}}, expected: "key[0]=1&key[1]=2&key[2]=true"},
			})
		})

		t.Run("map tyep value", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "string type value", q: qstring.Q{"key": map[string]string{"key1": "1", "key2": "2", "key3": "3"}}, expected: "key[key1]=1&key[key2]=2&key[key3]=3"},
				{name: "string type value pointer", q: qstring.Q{"key": &map[string]string{"key1": "1", "key2": "2", "key3": "3"}}, expected: "key[key1]=1&key[key2]=2&key[key3]=3"},
				{name: "interface type value", q: qstring.Q{"key": qstring.Q{"key1": 1, "key2": "2", "key3": true}}, expected: "key[key1]=1&key[key2]=2&key[key3]=true"},
				{name: "interface type value pointer", q: qstring.Q{"key": &qstring.Q{"key1": 1, "key2": "2", "key3": true}}, expected: "key[key1]=1&key[key2]=2&key[key3]=true"},
			})
		})

		t.Run("multiple tyep value", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{
					name: "string value map and interface value map into slice",
					q: qstring.Q{
						"key": qstring.S{
							map[string]string{"key1": "1", "key2": "2", "key3": "3"},
							qstring.Q{"key1": 1, "key2": "2", "key3": true},
						},
					},
					expected: "key[0][key1]=1&key[0][key2]=2&key[0][key3]=3&key[1][key1]=1&key[1][key2]=2&key[1][key3]=true",
				},
				{
					name: "nested map",
					q: qstring.Q{
						"key": qstring.Q{
							"key1": []string{"1", "2", "3"},
							"key2": qstring.S{1, "2", true},
						},
					},
					expected: "key[key1][0]=1&key[key1][1]=2&key[key1][2]=3&key[key2][0]=1&key[key2][1]=2&key[key2][2]=true",
				},
			})
		})

		t.Run("struct type vlaue", func(t *testing.T) {
			type s struct {
				Field string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: qstring.Q{"key": s{Field: "123"}}, expected: "key[field]=123"},
				{name: "empty value", q: qstring.Q{"key": s{}}, expected: "key[field]="},
			})
		})
	})

	t.Run("struct", func(t *testing.T) {
		t.Run("unsupported field type", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "uintptr", q: struct {
					Field uintptr `qstring:"field"`
				}{Field: uintptr(1)}, err: fmt.Errorf("uintptr is not supported")},
				{name: "complex64", q: struct {
					Field complex64 `qstring:"field"`
				}{Field: complex64(1)}, err: fmt.Errorf("complex64 is not supported")},
				{name: "complex128", q: struct {
					Field complex128 `qstring:"field"`
				}{Field: complex128(1)}, err: fmt.Errorf("complex128 is not supported")},
				{name: "func", q: struct {
					Field func() `qstring:"field"`
				}{Field: func() {}}, err: fmt.Errorf("func is not supported")},
				{name: "unsafe pointer", q: struct {
					Field unsafe.Pointer `qstring:"field"`
				}{Field: unsafeP("1")}, err: fmt.Errorf("unsafe.Pointer is not supported")},
			})
		})

		t.Run("excluded", func(t *testing.T) {
			runEncodeTest(t, []encodeCase{
				{name: "no tag", q: struct{ Field string }{}, expected: ""},
				{name: "private", q: struct{ field string }{}, expected: ""},
				{name: "private and has tag", q: struct {
					field string `qstring:"field"`
				}{}, expected: ""},
			})
		})

		t.Run("bool", func(t *testing.T) {
			type s struct {
				Field bool `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: true}, expected: "field=true"},
				{name: "empty value", q: s{}, expected: "field=false"},
			})
		})
		t.Run("bool omitempty", func(t *testing.T) {
			type s struct {
				Field bool `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: true}, expected: "field=true"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("bool pointer", func(t *testing.T) {
			type s struct {
				Field *bool `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: boolP(true)}, expected: "field=true"},
				{name: "empty value", q: s{Field: boolP(false)}, expected: "field=false"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("bool pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *bool `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: boolP(true)}, expected: "field=true"},
				{name: "empty value", q: s{Field: boolP(false)}, expected: "field=false"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("string", func(t *testing.T) {
			type s struct {
				Field string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: "123"}, expected: "field=123"},
				{name: "empty value", q: s{}, expected: "field="},
			})
		})
		t.Run("string omitempty", func(t *testing.T) {
			type s struct {
				Field string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: "123"}, expected: "field=123"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("string pointer", func(t *testing.T) {
			type s struct {
				Field *string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: stringP("123")}, expected: "field=123"},
				{name: "empty value", q: s{stringP("")}, expected: "field="},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("string pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: stringP("123")}, expected: "field=123"},
				{name: "empty value", q: s{stringP("")}, expected: "field="},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("int", func(t *testing.T) {
			type s struct {
				Field int `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("int omitempty", func(t *testing.T) {
			type s struct {
				Field int `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("int pointer", func(t *testing.T) {
			type s struct {
				Field *int `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: intP(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: intP(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("int pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *int `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: intP(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: intP(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("int64", func(t *testing.T) {
			type s struct {
				Field int64 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("int64 omitempty", func(t *testing.T) {
			type s struct {
				Field int64 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("int64 pointer", func(t *testing.T) {
			type s struct {
				Field *int64 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: int64P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: int64P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("int64 pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *int64 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: int64P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: int64P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("int32", func(t *testing.T) {
			type s struct {
				Field int32 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("int32 omitempty", func(t *testing.T) {
			type s struct {
				Field int32 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("int32 pointer", func(t *testing.T) {
			type s struct {
				Field *int32 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: int32P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: int32P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("int32 pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *int32 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: int32P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: int32P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("int16", func(t *testing.T) {
			type s struct {
				Field int16 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("int16 omitempty", func(t *testing.T) {
			type s struct {
				Field int16 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("int16 pointer", func(t *testing.T) {
			type s struct {
				Field *int16 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: int16P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: int16P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("int16 pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *int16 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: int16P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: int16P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("int8", func(t *testing.T) {
			type s struct {
				Field int8 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("int8 omitempty", func(t *testing.T) {
			type s struct {
				Field int8 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("int8 pointer", func(t *testing.T) {
			type s struct {
				Field *int8 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: int8P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: int8P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("int8 pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *int8 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: int8P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: int8P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("uint", func(t *testing.T) {
			type s struct {
				Field uint `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("uint omitempty", func(t *testing.T) {
			type s struct {
				Field uint `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("uint pointer", func(t *testing.T) {
			type s struct {
				Field *uint `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uintP(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uintP(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("uint pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *uint `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uintP(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uintP(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("uint64", func(t *testing.T) {
			type s struct {
				Field uint64 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("uint64 omitempty", func(t *testing.T) {
			type s struct {
				Field uint64 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("uint64 pointer", func(t *testing.T) {
			type s struct {
				Field *uint64 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uint64P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uint64P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("uint64 pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *uint64 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uint64P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uint64P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("uint32", func(t *testing.T) {
			type s struct {
				Field uint32 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("uint32 omitempty", func(t *testing.T) {
			type s struct {
				Field uint32 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("uint32 pointer", func(t *testing.T) {
			type s struct {
				Field *uint32 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uint32P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uint32P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("uint32 pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *uint32 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uint32P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uint32P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("uint16", func(t *testing.T) {
			type s struct {
				Field uint16 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("uint16 omitempty", func(t *testing.T) {
			type s struct {
				Field uint16 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("uint16 pointer", func(t *testing.T) {
			type s struct {
				Field *uint16 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uint16P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uint16P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("uint16 pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *uint16 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uint16P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uint16P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("uint8", func(t *testing.T) {
			type s struct {
				Field uint8 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: "field=0"},
			})
		})
		t.Run("uint8 omitempty", func(t *testing.T) {
			type s struct {
				Field uint8 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: 100}, expected: "field=100"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("uint8 pointer", func(t *testing.T) {
			type s struct {
				Field *uint8 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uint8P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uint8P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("uint8 pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *uint8 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: uint8P(100)}, expected: "field=100"},
				{name: "empty value", q: s{Field: uint8P(0)}, expected: "field=0"},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("array", func(t *testing.T) {
			type s struct {
				Field [3]string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: [3]string{"1", "2", "3"}}, expected: "field[0]=1&field[1]=2&field[2]=3"},
				{name: "empty value", q: s{}, expected: "field[0]=&field[1]=&field[2]="},
			})
		})
		t.Run("array omitempty", func(t *testing.T) {
			type s struct {
				Field [3]string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: [3]string{"1", "2", "3"}}, expected: "field[0]=1&field[1]=2&field[2]=3"},
				{name: "empty value", q: s{}, expected: "field[0]=&field[1]=&field[2]="},
			})
		})
		t.Run("array pointer", func(t *testing.T) {
			type s struct {
				Field *[3]string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &[3]string{"1", "2", "3"}}, expected: "field[0]=1&field[1]=2&field[2]=3"},
				{name: "empty value", q: s{Field: &[3]string{}}, expected: "field[0]=&field[1]=&field[2]="},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("array pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *[3]string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &[3]string{"1", "2", "3"}}, expected: "field[0]=1&field[1]=2&field[2]=3"},
				{name: "empty value", q: s{Field: &[3]string{}}, expected: "field[0]=&field[1]=&field[2]="},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("empty array", func(t *testing.T) {
			type s struct {
				Field [0]string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: [0]string{}}, expected: "field="},
				{name: "empty value", q: s{}, expected: "field="},
			})
		})
		t.Run("empty array omitempty", func(t *testing.T) {
			type s struct {
				Field [0]string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: [0]string{}}, expected: ""},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("empty array pointer", func(t *testing.T) {
			type s struct {
				Field *[0]string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &[0]string{}}, expected: "field="},
				{name: "empty value", q: s{Field: &[0]string{}}, expected: "field="},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("empty array pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *[0]string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &[0]string{}}, expected: "field="},
				{name: "empty value", q: s{Field: &[0]string{}}, expected: "field="},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("slice", func(t *testing.T) {
			type s struct {
				Field []string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: []string{"1", "2", "3"}}, expected: "field[0]=1&field[1]=2&field[2]=3"},
				{name: "empty value", q: s{}, expected: "field="},
				{name: "nil value", q: s{Field: nil}, expected: "field="},
			})
		})
		t.Run("slice omitempty", func(t *testing.T) {
			type s struct {
				Field []string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: []string{"1", "2", "3"}}, expected: "field[0]=1&field[1]=2&field[2]=3"},
				{name: "empty value", q: s{}, expected: ""},
				{name: "nil value", q: s{Field: nil}, expected: ""},
			})
		})
		t.Run("slice pointer", func(t *testing.T) {
			type s struct {
				Field *[]string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &[]string{"1", "2", "3"}}, expected: "field[0]=1&field[1]=2&field[2]=3"},
				{name: "empty value", q: s{Field: &[]string{}}, expected: "field="},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("slice pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *[]string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &[]string{"1", "2", "3"}}, expected: "field[0]=1&field[1]=2&field[2]=3"},
				{name: "empty value", q: s{Field: &[]string{}}, expected: "field="},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("interface", func(t *testing.T) {
			type s struct {
				Field interface{} `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: "123"}, expected: "field=123"},
				{name: "empty value", q: s{Field: ""}, expected: "field="},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("interface omitempty", func(t *testing.T) {
			type s struct {
				Field interface{} `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: "123"}, expected: "field=123"},
				{name: "empty value", q: s{Field: ""}, expected: "field="},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("string value map", func(t *testing.T) {
			type s struct {
				Field map[string]string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: map[string]string{"a": "1", "1": "2", "true": "3"}}, expected: "field[1]=2&field[a]=1&field[true]=3"},
				{name: "empty value", q: s{}, expected: "field="},
				{name: "nil value", q: s{Field: nil}, expected: "field="},
			})
		})
		t.Run("string value map omitempty", func(t *testing.T) {
			type s struct {
				Field map[string]string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: map[string]string{"a": "1", "1": "2", "true": "3"}}, expected: "field[1]=2&field[a]=1&field[true]=3"},
				{name: "empty value", q: s{}, expected: ""},
				{name: "nil value", q: s{Field: nil}, expected: ""},
			})
		})
		t.Run("string value map pointer", func(t *testing.T) {
			type s struct {
				Field *map[string]string `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &map[string]string{"a": "1", "1": "2", "true": "3"}}, expected: "field[1]=2&field[a]=1&field[true]=3"},
				{name: "empty value", q: s{Field: &map[string]string{}}, expected: "field="},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("string value map pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *map[string]string `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &map[string]string{"a": "1", "1": "2", "true": "3"}}, expected: "field[1]=2&field[a]=1&field[true]=3"},
				{name: "empty value", q: s{Field: &map[string]string{}}, expected: "field="},
				{name: "nil value", q: s{}, expected: ""},
			})
		})

		t.Run("interface value map", func(t *testing.T) {
			type s struct {
				Field map[string]interface{} `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: map[string]interface{}{"a": "1", "1": "2", "true": "3"}}, expected: "field[1]=2&field[a]=1&field[true]=3"},
				{name: "empty value", q: s{}, expected: "field="},
				{name: "nil value", q: s{Field: nil}, expected: "field="},
			})
		})
		t.Run("interface value map omitempty", func(t *testing.T) {
			type s struct {
				Field map[string]interface{} `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: map[string]interface{}{"a": "1", "1": "2", "true": "3"}}, expected: "field[1]=2&field[a]=1&field[true]=3"},
				{name: "empty value", q: s{}, expected: ""},
				{name: "nil value", q: s{Field: nil}, expected: ""},
			})
		})
		t.Run("interface value map pointer", func(t *testing.T) {
			type s struct {
				Field *map[string]interface{} `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &map[string]interface{}{"a": "1", "1": "2", "true": "3"}}, expected: "field[1]=2&field[a]=1&field[true]=3"},
				{name: "empty value", q: s{Field: &map[string]interface{}{}}, expected: "field="},
				{name: "nil value", q: s{}, expected: "field="},
			})
		})
		t.Run("interface value map pointer omitempty", func(t *testing.T) {
			type s struct {
				Field *map[string]interface{} `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: &map[string]interface{}{"a": "1", "1": "2", "true": "3"}}, expected: "field[1]=2&field[a]=1&field[true]=3"},
				{name: "empty value", q: s{Field: &map[string]interface{}{}}, expected: "field="},
				{name: "nil value", q: s{Field: nil}, expected: ""},
			})
		})

		t.Run("struct", func(t *testing.T) {
			type s2 struct {
				Field2 string `qstring:"field-2"`
			}
			type s struct {
				Field s2 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: s2{Field2: "123"}}, expected: "field[field-2]=123"},
				{name: "empty value", q: s{}, expected: "field[field-2]="},
			})
		})
		t.Run("struct omitempty", func(t *testing.T) {
			type s2 struct {
				Field2 string `qstring:"field-2,omitempty"`
			}
			type s struct {
				Field s2 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: s{Field: s2{Field2: "123"}}, expected: "field[field-2]=123"},
				{name: "empty value", q: s{}, expected: ""},
			})
		})
		t.Run("struct pointer", func(t *testing.T) {
			type s2 struct {
				Field2 *string `qstring:"field-2"`
			}
			type s struct {
				Field *s2 `qstring:"field"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: &s{Field: &s2{Field2: stringP("123")}}, expected: "field[field-2]=123"},
				{name: "empty child value", q: s{Field: &s2{Field2: stringP("")}}, expected: "field[field-2]="},
				{name: "nil child value", q: s{Field: &s2{Field2: nil}}, expected: "field[field-2]="},
				{name: "nil child struct", q: s{Field: nil}, expected: "field="},
			})
		})
		t.Run("struct pointer omitempty", func(t *testing.T) {
			type s2 struct {
				Field2 *string `qstring:"field-2,omitempty"`
			}
			type s struct {
				Field *s2 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: &s{Field: &s2{Field2: stringP("123")}}, expected: "field[field-2]=123"},
				{name: "empty child value", q: s{Field: &s2{Field2: stringP("")}}, expected: "field[field-2]="},
				{name: "nil child value", q: s{Field: &s2{Field2: nil}}, expected: ""},
				{name: "nil child struct", q: s{Field: nil}, expected: ""},
			})
		})
		t.Run("struct pointer omitempty (multiple fields)", func(t *testing.T) {
			type s2 struct {
				Field2 *string `qstring:"field-2,omitempty"`
				Field3 *string `qstring:"field-3,omitempty"`
			}
			type s struct {
				Field *s2 `qstring:"field,omitempty"`
			}
			runEncodeTest(t, []encodeCase{
				{name: "has value", q: &s{Field: &s2{Field2: stringP("123"), Field3: stringP("true")}}, expected: "field[field-2]=123&field[field-3]=true"},
				{name: "empty child value", q: s{Field: &s2{Field2: stringP(""), Field3: stringP("")}}, expected: "field[field-2]=&field[field-3]="},
				{name: "one is nil child value", q: s{Field: &s2{Field2: nil, Field3: stringP("true")}}, expected: "field[field-3]=true"},
				{name: "both nil child value", q: s{Field: &s2{Field2: nil, Field3: nil}}, expected: ""},
				{name: "nil child struct", q: s{Field: nil}, expected: ""},
			})
		})
	})
}

func runEncodeTest(t *testing.T, testCases []encodeCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := qstring.Encode(tc.q)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("Encode() should not returns error, got %q", err)
				}
				if err.Error() != tc.err.Error() {
					t.Fatalf("Encode() error returns %q, want %q", err, tc.err)
				}
			}

			if err == nil && tc.err != nil {
				t.Errorf("Encode() should returns error, want %q", tc.err)
			}

			a := strings.ReplaceAll(actual, "%5B", "[")
			a = strings.ReplaceAll(a, "%5D", "]")
			if a != tc.expected {
				t.Errorf("Encode() returns \n%q\nwant \n%q", a, tc.expected)
			}
		})
	}
}
