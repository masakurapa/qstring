package qstring_test

import (
	"fmt"

	"github.com/masakurapa/qstring"
)

func ExampleEncode_fromString() {
	s, _ := qstring.Encode("key[a]=1&key[b]=2")
	fmt.Println(s)

	// Output: key%5Ba%5D=1&key%5Bb%5D=2
}

func ExampleEncode_fromMap() {
	s, _ := qstring.Encode(qstring.Q{
		"key": qstring.Q{
			"a": "1",
			"b": "2",
		},
	})
	fmt.Println(s)

	// Output: key%5Ba%5D=1&key%5Bb%5D=2
}

func ExampleEncode_fromStruct() {
	type b struct {
		A string `qstring:"a"`
		B string `qstring:"b"`
		C string // untagged field will be ignored
		D string `qstring:"d,omitempty"` // "omitempty" ignores the field if it has a zero-value
	}
	type a struct {
		Key b `qstring:"key"`
	}

	s, _ := qstring.Encode(a{
		Key: b{
			A: "1",
			B: "2",
			C: "3",
			D: "",
		},
	})
	fmt.Println(s)

	// Output: key%5Ba%5D=1&key%5Bb%5D=2
}

func ExampleDecode_toString() {
	v := ""
	_ = qstring.Decode("key%5Ba%5D=1&key%5Bb%5D=2", &v)
	fmt.Println(v)

	// Output: key[a]=1&key[b]=2
}

func ExampleDecode_toArray() {
	v := [3]string{}
	_ = qstring.Decode("key%5B0%5D=a&key%5B1%5D=b&key%5B2%5D=c", &v)
	fmt.Println(v)

	// Output: [a b c]
}

func ExampleDecode_toSlice() {
	v := []string{}
	_ = qstring.Decode("key%5B0%5D=a&key%5B1%5D=b&key%5B2%5D=c", &v)
	fmt.Println(v)

	// Output: [a b c]
}

func ExampleDecode_toMap() {
	v := qstring.Q{}
	_ = qstring.Decode("key%5Ba%5D=1&key%5Bb%5D=2", &v)
	fmt.Println(v)

	// Output: map[key:map[a:1 b:2]]
}

func ExampleDecode_toStruct() {
	type b struct {
		A string `qstring:"a"`
		B string `qstring:"b"`
		C string // Untagged field will be ignored
	}
	type a struct {
		Key b `qstring:"key"`
	}

	v := a{}
	_ = qstring.Decode("key%5Ba%5D=1&key%5Bb%5D=2&key%5Bc%5D=3", &v)
	fmt.Printf("%+v", v)

	// Output: {Key:{A:1 B:2 C:}}
}

func ExampleDecodeToString() {
	q, _ := qstring.DecodeToString("key%5Ba%5D=1&key%5Bb%5D=2")
	fmt.Println(q)

	// Output: key[a]=1&key[b]=2
}

func ExampleDecodeToMap() {
	q, _ := qstring.DecodeToMap("key%5Ba%5D=1&key%5Bb%5D=2")
	fmt.Println(q)

	// Output: map[key:map[a:1 b:2]]
}

func ExampleDecodeToSlice() {
	q, _ := qstring.DecodeToSlice("key%5B0%5D=a&key%5B1%5D=b&key%5B2%5D=c")
	fmt.Println(q)

	// Output: [a b c]
}
