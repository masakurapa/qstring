package qstringer_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/masakurapa/qstringer"
)

func TestQstringer(t *testing.T) {
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
			{name: "array type (string)", q: qstringer.Q{"key": [3]string{"1", "2", "3"}}, expected: "?key%5B%5D=1&key%5B%5D=2&key%5B%5D=3"},
			{name: "array type (interface)", q: qstringer.Q{"key": [3]interface{}{1, "2", true}}, expected: "?key%5B%5D=1&key%5B%5D=2&key%5B%5D=true"},
			// slice type
			{name: "slice type (string)", q: qstringer.Q{"key": []string{"1", "2", "3"}}, expected: "?key%5B%5D=1&key%5B%5D=2&key%5B%5D=3"},
			{name: "slice type (interface)", q: qstringer.Q{"key": qstringer.ArrayQ{1, "2", true}}, expected: "?key%5B%5D=1&key%5B%5D=2&key%5B%5D=true"},
			// map type
			{name: "map type (string)", q: qstringer.Q{"key": map[string]string{"key1": "1", "key2": "2", "key3": "3"}}, expected: "?key%5Bkey1%5D=1&key%5Bkey2%5D=2&key%5Bkey3%5D=3"},
			{name: "map type (interface)", q: qstringer.Q{"key": qstringer.MapQ{"key1": 1, "key2": "2", "key3": true}}, expected: "?key%5Bkey1%5D=1&key%5Bkey2%5D=2&key%5Bkey3%5D=true"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				a, err := tc.q.Encode()
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
			q    qstringer.Q
			err  error
		}{
			// map key is not string type

			// unavailable types
			{name: "uintptr type", q: qstringer.Q{"key": uintptr(1)}, err: fmt.Errorf("type uintptr is not available (key: key)")},
			{name: "complex64 type", q: qstringer.Q{"key": complex64(1)}, err: fmt.Errorf("type complex64 is not available (key: key)")},
			{name: "complex128 type", q: qstringer.Q{"key": complex128(1)}, err: fmt.Errorf("type complex128 is not available (key: key)")},
			{name: "chan type", q: qstringer.Q{"key": make(chan int)}, err: fmt.Errorf("type chan is not available (key: key)")},
			{name: "func type", q: qstringer.Q{"key": func() {}}, err: fmt.Errorf("type func is not available (key: key)")},
			{name: "ptr type", q: qstringer.Q{"key": func() *string { return nil }()}, err: fmt.Errorf("type ptr is not available (key: key)")},
			{name: "struct type", q: qstringer.Q{"key": struct{}{}}, err: fmt.Errorf("type complex128 is not available (key: key)")},
			{name: "unsafe pointer type", q: qstringer.Q{"key": unsafe.Pointer(&v)}, err: fmt.Errorf("type unsafe.Pointer is not available (key: key)")},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				a, err := tc.q.Encode()
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
