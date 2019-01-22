package di

import (
	"errors"
	"reflect"
)

// BindType is the type of a binded object/value, it's being used to
// check if the value is accessible after a function call with a  "ctx" when needed ( Dynamic type)
// or it's just a struct value (a service | Static type).
type BindType uint32

const (
	// Static is the simple assignable value, a static value.
	Static BindType = iota
	// Dynamic returns a value but it depends on some input arguments from the caller,
	// on serve time.
	Dynamic
)

func bindTypeString(typ BindType) string {
	switch typ {
	case Dynamic:
		return "Dynamic"
	default:
		return "Static"
	}
}

// BindObject contains the dependency value's read-only information.
//
// StructInjector keeps information about their
// input arguments/or fields, these properties contain a `BindObject` inside them.
type BindObject struct {
	Type  reflect.Type // the Type of 'Value' or the type of the returned 'ReturnValue' .
	Value reflect.Value

	BindType    BindType
	ReturnValue func([]reflect.Value) reflect.Value
}

// MakeBindObject accepts any "v" value, struct, pointer or a function
func MakeBindObject(v reflect.Value) (b BindObject) {
	b.BindType = Static
	b.Type = v.Type()
	b.Value = v
	b.ReturnValue = func(c []reflect.Value) reflect.Value {
		return c[0]
	}
	return
}

var errBad = errors.New("bad")

// IsAssignable checks if "to" type can be used as "b.Value/ReturnValue".
func (b *BindObject) IsAssignable(to reflect.Type) bool {
	return equalTypes(b.Type, to)
}

// Assign sets the values to a setter, "toSetter" contains the setter, so the caller
// can use it for multiple and different structs/functions as well.
func (b *BindObject) Assign(ctx []reflect.Value, toSetter func(reflect.Value)) {
	if b.BindType == Dynamic {
		toSetter(b.ReturnValue(ctx))
		return
	}
	toSetter(b.Value)
}
