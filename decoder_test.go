package qstringer_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/masakurapa/qstringer"
)

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
		{name: "bool type", q: q, v: func() *bool { var a bool; return &a }(), err: fmt.Errorf("type bool is not available")},
		{name: "int type", q: q, v: func() *int { var a int; return &a }(), err: fmt.Errorf("type int is not available")},
		{name: "int64 type", q: q, v: func() *int64 { var a int64; return &a }(), err: fmt.Errorf("type int64 is not available")},
		{name: "int32 type", q: q, v: func() *int32 { var a int32; return &a }(), err: fmt.Errorf("type int32 is not available")},
		{name: "int16 type", q: q, v: func() *int16 { var a int16; return &a }(), err: fmt.Errorf("type int16 is not available")},
		{name: "int8 type", q: q, v: func() *int8 { var a int8; return &a }(), err: fmt.Errorf("type int8 is not available")},
		{name: "uint type", q: q, v: func() *uint { var a uint; return &a }(), err: fmt.Errorf("type uint is not available")},
		{name: "uint64 type", q: q, v: func() *uint64 { var a uint64; return &a }(), err: fmt.Errorf("type uint64 is not available")},
		{name: "uint32 type", q: q, v: func() *uint32 { var a uint32; return &a }(), err: fmt.Errorf("type uint32 is not available")},
		{name: "uint16 type", q: q, v: func() *uint16 { var a uint16; return &a }(), err: fmt.Errorf("type uint16 is not available")},
		{name: "uint8 type", q: q, v: func() *uint8 { var a uint8; return &a }(), err: fmt.Errorf("type uint8 is not available")},
		{name: "float64 type", q: q, v: func() *float64 { var a float64; return &a }(), err: fmt.Errorf("type float64 is not available")},
		{name: "float32 type", q: q, v: func() *float32 { var a float32; return &a }(), err: fmt.Errorf("type float32 is not available")},
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
		{name: "string type with quote", q: "?hoge[key]=fuga", v: func() *string { var a string; return &a }(), expected: "?hoge[key]=fuga"},
		{name: "string type without quote", q: "hoge[key]=fuga", v: func() *string { var a string; return &a }(), expected: "hoge[key]=fuga"},

		// array type
		{name: "array of string type1", q: "hoge[]=a", v: func() *[3]string { var a [3]string; return &a }(), expected: [3]string{"a", "", ""}},
		{name: "array of string type2", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *[3]string { var a [3]string; return &a }(), expected: [3]string{"a", "2", "3"}},
		{name: "array of string type - capacity exceeded", q: "hoge[]=a&hoge[]=2&hoge[]=3&hoge[]=4", v: func() *[3]string { var a [3]string; return &a }(), err: fmt.Errorf("array capacity exceeded")},
		{name: "array of int type", q: "hoge[]=1", v: func() *[3]int { var a [3]int; return &a }(), err: fmt.Errorf("allocation type must be [n]stirng")},

		// slice type
		{name: "array of string type", q: "hoge[]=a&hoge[]=2&hoge[]=3", v: func() *[]string { var a []string; return &a }(), expected: []string{"a", "2", "3"}},
		{name: "array of int type", q: "hoge[]=1", v: func() *[]int { var a []int; return &a }(), err: fmt.Errorf("allocation type must be []stirng")},

		// map type

		// struct type
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := strings.ReplaceAll(tc.q, "%5B", "[")
			q = strings.ReplaceAll(q, "%5D", "]")
			err := qstringer.Decode(q, tc.v)
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
					t.Errorf("Decode() returns \n%v\nwant \n%v", tc.v, tc.expected)
				}
			}
		})
	}
}
