package di

import (
	"reflect"
)

// Struct is being used to return a new injector based on
// a struct value instance, if it contains fields that the types of those
// are matching with one or more of the `Values` then they are binded
// with the injector's `Inject` and `InjectElem` methods.
func Struct(s interface{}, values ...reflect.Value) *StructInjector {
	if s == nil {
		return &StructInjector{Has: false}
	}

	return MakeStructInjector(
		ValueOf(s),
		Values(values).CloneWithFieldsOf(s)...,
	)
}

// D is the Dependency Injection container,
// it contains the Values that can be changed before the injectors.
// `Struct` and the `Func` methods returns an injector for specific
// struct instance-value or function.
type D struct {
	Values
}

// New creates and returns a new Dependency Injection container.
// See `Values` field and `Func` and `Struct` methods for more.
func New() *D {
	return &D{}
}

// Clone returns a new Dependency Injection container, it adopts the
// parent's (current "D") hijacker, good func type checker and all dependencies values.
func (d *D) Clone() *D {
	return &D{
		Values: d.Values.Clone(),
	}
}

// Struct is being used to return a new injector based on
// a struct value instance, if it contains fields that the types of those
// are matching with one or more of the `Values` then they are binded
// with the injector's `Inject` and `InjectElem` methods.
func (d *D) Struct(s interface{}) *StructInjector {
	if s == nil {
		return &StructInjector{Has: false}
	}

	return MakeStructInjector(
		ValueOf(s),
		d.Values.CloneWithFieldsOf(s)...,
	)
}
