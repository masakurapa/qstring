package qstring

// Q is the type of the query string parameters.
type Q map[string]interface{}

// ArrayQ is a type of query string in array format.
type ArrayQ []interface{}

// Encode returns the URL-encoded query string.
//
// Support struct, map type where the key is a string.
//
// By default, a nil value will be output.
//
// If you don't want to output nil values,
// Please specify option "omitempty" in the tag.
func Encode(v interface{}) (string, error) {
	e := encoder{}
	return e.encode(v)
}

// Decode returns the URL-encoded query string
//
// add "?" to the beginning and return
func Decode(s string, v interface{}) error {
	d := decoder{query: s}
	return d.decode(v)
}

func DecodeToMap(s string) (Q, error) {
	var q Q
	err := Decode(s, &q)
	if err != nil {
		return nil, err
	}
	return q, nil
}
