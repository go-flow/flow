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
		Container(values).CloneWithFieldsOf(s)...,
	)
}
