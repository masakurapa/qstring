package qstring

import (
	"reflect"
	"strconv"
)

// unsupportedTypeError is an error for unsupported types
type unsupportedTypeError struct {
	rt reflect.Type
}

func (e *unsupportedTypeError) Error() string {
	t := e.rt.String()
	switch e.rt.Kind() {
	case reflect.Func:
		t = "func"
	case reflect.Chan:
		t = "chan"
	}
	return t + " is not supported"
}

// invalidEncodeError is an error for unencodable arguments
type invalidEncodeError struct {
	rt reflect.Type
}

func (e *invalidEncodeError) Error() string {
	if e.rt == nil {
		return "nil is not supported"
	}
	// this error should only nil or nil-pointer errors
	return "nil " + e.rt.Kind().String() + " is not supported"
}

// InvalidEncodeError is an error for undecodable arguments
type invalidDecodeError struct {
	rt reflect.Type
}

func (e *invalidDecodeError) Error() string {
	if e.rt == nil {
		return "nil is not supported"
	}
	if e.rt.Kind() != reflect.Ptr {
		return "non-pointer is not supported"
	}
	return "nil " + e.rt.Kind().String() + " is not supported"
}

// arrayIndexOutOfRangeDecodeError is an error
// when the capacity of the array is exceeded during decoding
type arrayIndexOutOfRangeDecodeError struct {
	rt  reflect.Type
	len int
}

func (e *arrayIndexOutOfRangeDecodeError) Error() string {
	return "index out of range [" + strconv.Itoa(e.len) + "] with " + e.rt.String()
}

// noAssignableValueError is an error
// if the value cannot be assigned to a structure
type noAssignableValueError struct {
	rt    reflect.Type
	value string
}

func (e *noAssignableValueError) Error() string {
	return `"` + e.value + `" can not be assign to ` + e.rt.String()
}

// multipleKeysError is an error
// when attempting to convert a query string with multiple keys
// into an array or slice
type multipleKeysError struct{}

func (e *multipleKeysError) Error() string {
	return "cannot decode due to multiple keys"
}
