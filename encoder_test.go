package qstringer_test

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/masakurapa/qstringer"
)

type s struct {
	FieldB      bool              `qstringer:"field_b"`
	FieldI      int               `qstringer:"fieldI"`
	FieldUIP    *uint             `qstringer:"field-uip"`
	JSONStr     string            `qstringer:"json_str"`
	JSONStrP    *string           `qstringer:"json_str-p"`
	Struct2     s2                `qstringer:"struct2"`
	Slice4Str   []string          `qstringer:"slice-for-str"`
	Map_Str_Str map[string]string `qstringer:"map_strStr"`
	Interface   interface{}       `qstringer:"interface"`
	NoTag       string
	privateS    string `qstringer:"private-S"`
}

type s2 struct {
	Field    string `qstringer:"field"`
	NoTag    string
	privateS string `qstringer:"private-S"`
}

func TestEncode(t *testing.T) {
	v := "1"
	var uv uint = 200

	sv := s{
		FieldB:      true,
		FieldI:      100,
		FieldUIP:    &uv,
		JSONStr:     "hoge",
		JSONStrP:    &v,
		Struct2:     s2{Field: "fuga", NoTag: "ninja", privateS: "fuji"},
		Slice4Str:   []string{"1", "3", "5"},
		Map_Str_Str: map[string]string{"k1": "2", "k2": "4", "k3": "6"},
		Interface:   "gumi",
		NoTag:       "ika",
		privateS:    "kani",
	}

	testCases := []struct {
		name     string
		q        interface{}
		err      error
		expected string
	}{
		{name: "nil", q: nil, err: fmt.Errorf("nil is not available")},

		// unavailable types
		{name: "bool type", q: true, err: fmt.Errorf("type bool is not available")},
		{name: "int type", q: int(123), err: fmt.Errorf("type int is not available")},
		{name: "int64 type", q: int64(123), err: fmt.Errorf("type int64 is not available")},
		{name: "int32 type", q: int32(123), err: fmt.Errorf("type int32 is not available")},
		{name: "int16 type", q: int16(123), err: fmt.Errorf("type int16 is not available")},
		{name: "int8 type", q: int8(123), err: fmt.Errorf("type int8 is not available")},
		{name: "uint type", q: uint(123), err: fmt.Errorf("type uint is not available")},
		{name: "uint64 type", q: uint64(123), err: fmt.Errorf("type uint64 is not available")},
		{name: "uint32 type", q: uint32(123), err: fmt.Errorf("type uint32 is not available")},
		{name: "uint16 type", q: uint16(123), err: fmt.Errorf("type uint16 is not available")},
		{name: "uint8 type", q: uint8(123), err: fmt.Errorf("type uint8 is not available")},
		{name: "float64 type", q: float64(123.456), err: fmt.Errorf("type float64 is not available")},
		{name: "float32 type", q: float32(123.456), err: fmt.Errorf("type float32 is not available")},
		{name: "array type", q: [3]string{"1", "2", "3"}, err: fmt.Errorf("type array is not available")},
		{name: "slice type", q: []string{"1", "2", "3"}, err: fmt.Errorf("type slice is not available")},
		{name: "uintptr type", q: uintptr(1), err: fmt.Errorf("type uintptr is not available")},
		{name: "complex64 type", q: complex64(1), err: fmt.Errorf("type complex64 is not available")},
		{name: "complex128 type", q: complex128(1), err: fmt.Errorf("type complex128 is not available")},
		{name: "chan type", q: make(chan int), err: fmt.Errorf("type chan is not available")},
		{name: "func type", q: func() {}, err: fmt.Errorf("type func is not available")},
		{name: "ptr type", q: func() *uint { return &uv }(), err: fmt.Errorf("type uint is not available")},
		{name: "nil ptr type", q: func() *string { return nil }(), err: fmt.Errorf("nil is not available")},
		{name: "unsafe pointer type", q: unsafe.Pointer(&v), err: fmt.Errorf("type unsafe.Pointer is not available")},

		//
		// string type
		//
		{name: "string type with quote", q: "?hoge[key]=fuga", expected: "?hoge[key]=fuga"},
		{name: "string type without quote", q: "hoge[key]=fuga", expected: "hoge[key]=fuga"},

		//
		// map type
		//
		{name: "empty map", q: qstringer.Q{}, expected: ""},
		{name: "empty map pointer", q: &qstringer.Q{}, expected: ""},
		// map value type is bool
		{name: "map value type is bool (true)", q: qstringer.Q{"key": true}, expected: "key=true"},
		{name: "map value type is bool (false)", q: qstringer.Q{"key": false}, expected: "key=false"},
		// map value type is int
		{name: "map value type is int", q: qstringer.Q{"key": int(123)}, expected: "key=123"},
		{name: "map value type is int64", q: qstringer.Q{"key": int64(123)}, expected: "key=123"},
		{name: "map value type is int32", q: qstringer.Q{"key": int32(123)}, expected: "key=123"},
		{name: "map value type is int16", q: qstringer.Q{"key": int16(123)}, expected: "key=123"},
		{name: "map value type is int8", q: qstringer.Q{"key": int8(123)}, expected: "key=123"},
		// map value type is uint
		{name: "map value type is uint", q: qstringer.Q{"key": uint(123)}, expected: "key=123"},
		{name: "map value type is uint64", q: qstringer.Q{"key": uint64(123)}, expected: "key=123"},
		{name: "map value type is uint32", q: qstringer.Q{"key": uint32(123)}, expected: "key=123"},
		{name: "map value type is uint16", q: qstringer.Q{"key": uint16(123)}, expected: "key=123"},
		{name: "map value type is uint8", q: qstringer.Q{"key": uint8(123)}, expected: "key=123"},
		// map value type is float
		{name: "map value type is float64", q: qstringer.Q{"key": float64(123.456)}, expected: "key=123.456"},
		{name: "map value type is float32", q: qstringer.Q{"key": float32(123.456)}, expected: "key=123.456"},
		// map value type is string
		{name: "map value type is string", q: qstringer.Q{"key": "hoge"}, expected: "key=hoge"},
		// map value type is ptr
		{name: "map value type is ptr", q: qstringer.Q{"key": func() *string {
			s := "pointer"
			return &s
		}()}, expected: "key=pointer"},
		// map value type is array
		{name: "map value type is array (string)", q: qstringer.Q{"key": [3]string{"1", "2", "3"}}, expected: "key[0]=1&key[1]=2&key[2]=3"},
		{name: "map value type is array (interface)", q: qstringer.Q{"key": [3]interface{}{1, "2", true}}, expected: "key[0]=1&key[1]=2&key[2]=true"},
		// map value type is slice
		{name: "map value type is slice (string)", q: qstringer.Q{"key": []string{"1", "2", "3"}}, expected: "key[0]=1&key[1]=2&key[2]=3"},
		{name: "map value type is slice (interface)", q: qstringer.Q{"key": qstringer.ArrayQ{1, "2", true}}, expected: "key[0]=1&key[1]=2&key[2]=true"},
		// map value type is map
		{name: "map value type is map (string)", q: qstringer.Q{"key": map[string]string{"key1": "1", "key2": "2", "key3": "3"}}, expected: "key[key1]=1&key[key2]=2&key[key3]=3"},
		{name: "map value type is map (interface)", q: qstringer.Q{"key": qstringer.Q{"key1": 1, "key2": "2", "key3": true}}, expected: "key[key1]=1&key[key2]=2&key[key3]=true"},
		// map value type is slice
		{
			name: "map value type is slice (interface)",
			q: qstringer.Q{"key": qstringer.ArrayQ{
				map[string]string{"key1": "1", "key2": "2", "key3": "3"},
				qstringer.Q{"key1": 1, "key2": "2", "key3": true},
			}},
			expected: "key[0][key1]=1&key[0][key2]=2&key[0][key3]=3&key[1][key1]=1&key[1][key2]=2&key[1][key3]=true",
		},
		// map value type is map
		{
			name: "map value type is map",
			q: qstringer.Q{"key": qstringer.Q{
				"key1": []string{"1", "2", "3"},
				"key2": qstringer.ArrayQ{1, "2", true},
			}},
			expected: "key[key1][0]=1&key[key1][1]=2&key[key1][2]=3&key[key2][0]=1&key[key2][1]=2&key[key2][2]=true",
		},
		// unavailable map value types
		{name: "map value type is uintptr", q: qstringer.Q{"key": uintptr(1)}, err: fmt.Errorf("type uintptr is not available (key: key)")},
		{name: "map value type is complex64", q: qstringer.Q{"key": complex64(1)}, err: fmt.Errorf("type complex64 is not available (key: key)")},
		{name: "map value type is complex128", q: qstringer.Q{"key": complex128(1)}, err: fmt.Errorf("type complex128 is not available (key: key)")},
		{name: "map value type is chan", q: qstringer.Q{"key": make(chan int)}, err: fmt.Errorf("type chan is not available (key: key)")},
		{name: "map value type is func", q: qstringer.Q{"key": func() {}}, err: fmt.Errorf("type func is not available (key: key)")},
		{
			name: "map value type is struct",
			q:    qstringer.Q{"key": sv},
			expected: strings.Join([]string{
				"key[field-uip]=200",
				"key[fieldI]=100",
				"key[field_b]=true",
				"key[interface]=gumi",
				"key[json_str-p]=1",
				"key[json_str]=hoge",
				"key[map_strStr][k1]=2&key[map_strStr][k2]=4&key[map_strStr][k3]=6",
				"key[slice-for-str][0]=1&key[slice-for-str][1]=3&key[slice-for-str][2]=5",
				"key[struct2][field]=fuga",
			}, "&"),
		},

		{name: "map value type is unsafe pointer", q: qstringer.Q{"key": unsafe.Pointer(&v)}, err: fmt.Errorf("type unsafe.Pointer is not available (key: key)")},
		// map key is not string type
		{name: "uintptr type", q: map[int]string{100: "val"}, err: fmt.Errorf("the key of the map type must be a string")},

		//
		// struct type
		//
		{name: "empty struct", q: struct{}{}, expected: ""},
		{
			name: "struct fields are private",
			q: struct {
				field1 bool
				field2 int
				field3 string
			}{
				field1: true,
				field2: 100,
				field3: "value",
			},
			expected: "",
		},
		{
			name: "struct value",
			q:    sv,
			expected: strings.Join([]string{
				"field-uip=200",
				"fieldI=100",
				"field_b=true",
				"interface=gumi",
				"json_str=hoge",
				"json_str-p=1",
				"map_strStr[k1]=2&map_strStr[k2]=4&map_strStr[k3]=6",
				"slice-for-str[0]=1&slice-for-str[1]=3&slice-for-str[2]=5",
				"struct2[field]=fuga",
			}, "&"),
		},
		{
			name: "struct value is empty and not specify omitempty",
			q:    s{},
			expected: strings.Join([]string{
				"field-uip=",
				"fieldI=0",
				"field_b=false",
				"interface=",
				"json_str=",
				"json_str-p=",
				"map_strStr=",
				"slice-for-str=",
				"struct2[field]=",
			}, "&"),
		},
		{
			name: "struct value is empty and specify omitempty",
			q: struct {
				FieldB   bool    `qstringer:"field_b,omitempty"`
				FieldI   int     `qstringer:"fieldI,omitempty"`
				FieldUIP *uint   `qstringer:"field-uip,omitempty"`
				JSONStr  string  `qstringer:"json_str,omitempty"`
				JSONStrP *string `qstringer:"json_str-p,omitempty"`
				Struct2  struct {
					Field    string `qstringer:"field,omitempty"`
					NoTag    string
					privateS string `qstringer:"private-S,omitempty"`
				} `qstringer:"struct2,omitempty"`
				Slice4Str   []string          `qstringer:"slice-for-str,omitempty"`
				Map_Str_Str map[string]string `qstringer:"map_strStr,omitempty"`
				Interface   interface{}       `qstringer:"interface,omitempty"`
				NoTag       string
				privateS    string `qstringer:"private-S,omitempty"`
			}{},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := qstringer.Encode(tc.q)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("Encode() should not returns error, got %q", err)
				}
				if err.Error() != tc.err.Error() {
					t.Fatalf("Encode() error returns %q, got %q", err, tc.err)
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

func BenchmarkEncode(b *testing.B) {
	sv := s{
		FieldB:      true,
		FieldI:      100,
		JSONStr:     "hoge",
		Struct2:     s2{Field: "fuga"},
		Slice4Str:   []string{"1", "3", "5"},
		Map_Str_Str: map[string]string{"k1": "2", "k2": "4", "k3": "6"},
		Interface:   "gumi",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := qstringer.Encode(sv)
		if err != nil {
			b.Fatal(err)
		}
	}
}
