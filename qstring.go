package qstring

// Q is the type of the query string parameters.
type Q map[string]interface{}

// S is a type of query string in slice format.
type S []interface{}

// Encode returns the URL-encoded query string.
//
// The argument supports
// string type, struct, map type where the key is a string.
//
// The struct needs to specify the "qstring" tag in the public field.
// If you don't want to output zero-value, Please specify option "omitempty" in the tag.
func Encode(v interface{}) (string, error) {
	e := encoder{}
	return e.encode(v)
}

// Decode is URL-decodes query string.
//
// The second argument supports
// string type, array, slice, struct, map type where the key is a string.
//
// The struct needs to specify the "qstring" tag in the public field.
func Decode(s string, v interface{}) error {
	d := decoder{query: s}
	return d.decode(v)
}

// DecodeToString returns the URL-decoded query string.
func DecodeToString(s string) (string, error) {
	var v string
	err := Decode(s, &v)
	if err != nil {
		return "", err
	}
	return v, nil
}

// DecodeToMap returns the URL-decoded query string as map type
func DecodeToMap(s string) (Q, error) {
	var v Q
	err := Decode(s, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeToMap returns the URL-decoded query string as slice type
func DecodeToSlice(s string) ([]string, error) {
	var v []string
	err := Decode(s, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
