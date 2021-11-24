package qstringer_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/masakurapa/qstringer"
)

func TestQ_Encode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testCases := []struct {
			name     string
			q        qstringer.Q
			expected string
		}{
			{name: "empty map", q: qstringer.Q{}, expected: ""},
			// bool type
			{name: "bool type (true)", q: qstringer.Q{"key": true}, expected: "?key=true"},
			{name: "bool type (false)", q: qstringer.Q{"key": false}, expected: "?key=false"},
			// int type
			{name: "int type", q: qstringer.Q{"key": int(123)}, expected: "?key=123"},
			{name: "int64 type", q: qstringer.Q{"key": int64(123)}, expected: "?key=123"},
			{name: "int32 type", q: qstringer.Q{"key": int32(123)}, expected: "?key=123"},
			{name: "int16 type", q: qstringer.Q{"key": int16(123)}, expected: "?key=123"},
			{name: "int8 type", q: qstringer.Q{"key": int8(123)}, expected: "?key=123"},
			// uint type
			{name: "uint type", q: qstringer.Q{"key": uint(123)}, expected: "?key=123"},
			{name: "uint64 type", q: qstringer.Q{"key": uint64(123)}, expected: "?key=123"},
			{name: "uint32 type", q: qstringer.Q{"key": uint32(123)}, expected: "?key=123"},
			{name: "uint16 type", q: qstringer.Q{"key": uint16(123)}, expected: "?key=123"},
			{name: "uint8 type", q: qstringer.Q{"key": uint8(123)}, expected: "?key=123"},
			// float type
			{name: "float64 type", q: qstringer.Q{"key": float64(123.456)}, expected: "?key=123.456"},
			{name: "float32 type", q: qstringer.Q{"key": float32(123.456)}, expected: "?key=123.456"},
			// string type
			{name: "string type", q: qstringer.Q{"key": "hoge"}, expected: "?key=hoge"},
			// array type
			{name: "array type (string)", q: qstringer.Q{"key": [3]string{"1", "2", "3"}}, expected: "?key%5B0%5D=1&key%5B1%5D=2&key%5B2%5D=3"},
			{name: "array type (interface)", q: qstringer.Q{"key": [3]interface{}{1, "2", true}}, expected: "?key%5B0%5D=1&key%5B1%5D=2&key%5B2%5D=true"},
			// slice type
			{name: "slice type (string)", q: qstringer.Q{"key": []string{"1", "2", "3"}}, expected: "?key%5B0%5D=1&key%5B1%5D=2&key%5B2%5D=3"},
			{name: "slice type (interface)", q: qstringer.Q{"key": qstringer.ArrayQ{1, "2", true}}, expected: "?key%5B0%5D=1&key%5B1%5D=2&key%5B2%5D=true"},
			// map type
			{name: "map type (string)", q: qstringer.Q{"key": map[string]string{"key1": "1", "key2": "2", "key3": "3"}}, expected: "?key%5Bkey1%5D=1&key%5Bkey2%5D=2&key%5Bkey3%5D=3"},
			{name: "map type (interface)", q: qstringer.Q{"key": qstringer.MapQ{"key1": 1, "key2": "2", "key3": true}}, expected: "?key%5Bkey1%5D=1&key%5Bkey2%5D=2&key%5Bkey3%5D=true"},

			// map type inside slice type
			{name: "slice type (interface)", q: qstringer.Q{"key": qstringer.ArrayQ{
				map[string]string{"key1": "1", "key2": "2", "key3": "3"},
				qstringer.MapQ{"key1": 1, "key2": "2", "key3": true},
			}}, expected: "?key%5B0%5D%5Bkey1%5D=1&key%5B0%5D%5Bkey2%5D=2&key%5B0%5D%5Bkey3%5D=3&key%5B1%5D%5Bkey1%5D=1&key%5B1%5D%5Bkey2%5D=2&key%5B1%5D%5Bkey3%5D=true"},

			// slice type inside map type
			{name: "slice type (interface)", q: qstringer.Q{"key": qstringer.MapQ{
				"key1": []string{"1", "2", "3"},
				"key2": qstringer.ArrayQ{1, "2", true},
			}}, expected: "?key%5Bkey1%5D%5B0%5D=1&key%5Bkey1%5D%5B1%5D=2&key%5Bkey1%5D%5B2%5D=3&key%5Bkey2%5D%5B0%5D=1&key%5Bkey2%5D%5B1%5D=2&key%5Bkey2%5D%5B2%5D=true"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				a, err := qstringer.Encode(tc.q)
				if err != nil {
					t.Fatalf("Encode() should not returns error: %s", err)
				}
				if a != tc.expected {
					t.Errorf("Encode() returns %q, want %q", a, tc.expected)
				}
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		v := "1"

		testCases := []struct {
			name string
			q    interface{}
			err  error
		}{
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
			{name: "string type", q: "hoge", err: fmt.Errorf("type string is not available")},
			{name: "array type", q: [3]string{"1", "2", "3"}, err: fmt.Errorf("type array is not available")},
			{name: "slice type", q: []string{"1", "2", "3"}, err: fmt.Errorf("type slice is not available")},
			{name: "uintptr type", q: uintptr(1), err: fmt.Errorf("type uintptr is not available")},
			{name: "complex64 type", q: complex64(1), err: fmt.Errorf("type complex64 is not available")},
			{name: "complex128 type", q: complex128(1), err: fmt.Errorf("type complex128 is not available")},
			{name: "chan type", q: make(chan int), err: fmt.Errorf("type chan is not available")},
			{name: "func type", q: func() {}, err: fmt.Errorf("type func is not available")},
			{name: "ptr type", q: func() *string { return nil }(), err: fmt.Errorf("type ptr is not available")},
			{name: "struct type", q: struct{}{}, err: fmt.Errorf("type struct is not available")},
			{name: "unsafe pointer type", q: unsafe.Pointer(&v), err: fmt.Errorf("type unsafe.Pointer is not available")},

			// unavailable map value types
			{name: "map vlalue of uintptr type", q: qstringer.Q{"key": uintptr(1)}, err: fmt.Errorf("type uintptr is not available (key: key)")},
			{name: "map vlalue of complex64 type", q: qstringer.Q{"key": complex64(1)}, err: fmt.Errorf("type complex64 is not available (key: key)")},
			{name: "map vlalue of complex128 type", q: qstringer.Q{"key": complex128(1)}, err: fmt.Errorf("type complex128 is not available (key: key)")},
			{name: "map vlalue of chan type", q: qstringer.Q{"key": make(chan int)}, err: fmt.Errorf("type chan is not available (key: key)")},
			{name: "map vlalue of func type", q: qstringer.Q{"key": func() {}}, err: fmt.Errorf("type func is not available (key: key)")},
			{name: "map vlalue of ptr type", q: qstringer.Q{"key": func() *string { return nil }()}, err: fmt.Errorf("type ptr is not available (key: key)")},
			{name: "map vlalue of struct type", q: qstringer.Q{"key": struct{}{}}, err: fmt.Errorf("type struct is not available (key: key)")},
			{name: "map vlalue of unsafe pointer type", q: qstringer.Q{"key": unsafe.Pointer(&v)}, err: fmt.Errorf("type unsafe.Pointer is not available (key: key)")},

			// map key is not string type
			{name: "uintptr type", q: map[int]string{100: "val"}, err: fmt.Errorf("the key of the map type must be a string")},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				a, err := qstringer.Encode(tc.q)
				if err == nil {
					t.Fatalf("Encode() should returns error")
				}
				if err.Error() != tc.err.Error() {
					t.Errorf("Encode() error returns %s, want %s", err, tc.err)
				}
				if a != "" {
					t.Errorf("Encode() returns %q, want %q", a, "")
				}
			})
		}
	})
}
