package qstring

import (
	"reflect"
	"strconv"
)

// UnsupportedTypeError is an error for unsupported types
type UnsupportedTypeError struct {
	rt reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	t := e.rt.String()
	switch e.rt.Kind() {
	case reflect.Func:
		t = "func"
	case reflect.Chan:
		t = "chan"
	}
	return t + " is not supported"
}

// InvalidEncodeError is an error for unencodable arguments
type InvalidEncodeError struct {
	rt reflect.Type
}

func (e *InvalidEncodeError) Error() string {
	if e.rt == nil {
		return "nil is not supported"
	}
	// this error should only nil or nil-pointer errors
	return "nil " + e.rt.Kind().String() + " is not supported"
}

// InvalidEncodeError is an error for undecodable arguments
type InvalidDecodeError struct {
	rt reflect.Type
}

func (e *InvalidDecodeError) Error() string {
	if e.rt == nil {
		return "nil is not supported"
	}
	if e.rt.Kind() != reflect.Ptr {
		return "non-pointer is not supported"
	}
	return "nil " + e.rt.Kind().String() + " is not supported"
}

// ArrayIndexOutOfRangeDecodeError is an error
// when the capacity of the array is exceeded during decoding
type ArrayIndexOutOfRangeDecodeError struct {
	rt  reflect.Type
	len int
}

func (e *ArrayIndexOutOfRangeDecodeError) Error() string {
	return "index out of range [" + strconv.Itoa(e.len) + "] with " + e.rt.String()
}

// NoAssignableValueError is an error
// if the value cannot be assigned to a structure
type NoAssignableValueError struct {
	rt    reflect.Type
	value string
}

func (e *NoAssignableValueError) Error() string {
	return `"` + e.value + `" can not be assign to ` + e.rt.String()
}
