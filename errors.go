package qstring

import (
	"reflect"
	"strconv"
)

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

type ArrayIndexOutOfRangeDecodeError struct {
	rt  reflect.Type
	len int
}

func (e *ArrayIndexOutOfRangeDecodeError) Error() string {
	return "index out of range [" + strconv.Itoa(e.len) + "] with " + e.rt.String()
}
