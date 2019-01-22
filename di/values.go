package di

import (
	"reflect"
)

// Values is a shortcut for []reflect.Value
type Values []reflect.Value

// NewValues returns new empty values
func NewValues() Values {
	return Values{}
}

// Clone returns a copy of the current value
func (v Values) Clone() Values {
	if n := len(v); n > 0 {
		values := make(Values, n, n)
		copy(values, v)
		return values
	}
	return NewValues()
}

// CloneWithFieldsOf will return a copy of the current values
// with provided struct fields that are filled(non-zero) by the caller
func (v Values) CloneWithFieldsOf(i interface{}) Values {
	values := v.Clone()

	// add the manual filled fields to the dependencies.
	filledFieldValues := LookupNonZeroFieldsValues(ValueOf(i), true)
	values = append(values, filledFieldValues...)
	return values
}

// Len returns Length of current Values slice
func (v Values) Len() int {
	return len(v)
}

// Add adds values as dependencies, if the struct's fields
// or the function's input arguments needs them, they will be defined as
// bindings (at build-time) and they will be used (at serve-time).
func (v *Values) Add(value interface{}) {
	v.AddValue(ValueOf(value))
}

// AddValue -
func (v *Values) AddValue(val reflect.Value) {
	if !goodVal(val) {
		return
	}
	*v = append(*v, val)
}

// Remove unbinds a binding value based on the type,
// it returns true if at least one field is not binded anymore.
//
// The "n" indicates the number of elements to remove, if <=0 then it's 1,
// this is useful because you may have bind more than one value to two or more fields
// with the same type.
func (v *Values) Remove(value interface{}, n int) bool {
	return v.remove(reflect.TypeOf(value), n)
}

func (v *Values) remove(typ reflect.Type, n int) (ok bool) {
	input := *v
	for i, in := range input {
		if equalTypes(in.Type(), typ) {
			ok = true
			input = input[:i+copy(input[i:], input[i+1:])]
			if n > 1 {
				continue
			}
			break
		}
	}

	*v = input

	return
}

// Has returns true if a binder responsible to
// bind and return a type of "typ" is already registered to this controller.
func (v Values) Has(value interface{}) bool {
	return v.valueTypeExists(reflect.TypeOf(value))
}

func (v Values) valueTypeExists(typ reflect.Type) bool {
	for _, in := range v {
		if equalTypes(in.Type(), typ) {
			return true
		}
	}
	return false
}

// AddOnce binds a value to the controller's field with the same type,
// if it's not binded already.
//
// Returns false if binded already or the value is not the proper one for binding,
// otherwise true.
func (v *Values) AddOnce(value interface{}) bool {
	return v.addIfNotExists(reflect.ValueOf(value))
}

func (v *Values) addIfNotExists(val reflect.Value) bool {
	var (
		typ = val.Type() // no element, raw things here.
	)

	if !goodVal(val) {
		return false
	}

	if v.valueTypeExists(typ) {
		return false
	}

	v.Add(val)
	return true
}
