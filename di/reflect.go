package di

import "reflect"

// EmptyIn is just an empty slice of reflect.Value.
var EmptyIn = []reflect.Value{}

// ValueOf returns the reflect.Value of "o".
// If "o" is already a reflect.Value returns "o".
func ValueOf(o interface{}) reflect.Value {
	if v, ok := o.(reflect.Value); ok {
		return v
	}

	return reflect.ValueOf(o)
}

func goodVal(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		if v.IsNil() {
			return false
		}
	}

	return v.IsValid()
}

func equalTypes(got reflect.Type, expected reflect.Type) bool {
	if got == expected {
		return true
	}

	// if accepts an interface, check if the given "got" type does
	// implement this "expected" user handler's input argument.
	if expected.Kind() == reflect.Interface {
		return got.Implements(expected)
	}

	return false
}
