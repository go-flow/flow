package di

import (
	"reflect"
	"sync"
)

var once sync.Once

var gc Container

// Global returns Global Container instance
func Global() Container {
	once.Do(func() {
		gc = NewContainer()
	})

	return gc
}

// Clone returns a copy of the current container value
func Clone() Container {
	return Global().Clone()
}

// Len returns Length of the current Container slice
func Len() int {
	return Global().Len()
}

// Add adds values as dependencies, if the struct's fields
// or the function's input arguments needs them, they will be defined as
// bindings (at build-time) and they will be used (at serve-time).
func Add(value interface{}) {
	g := Global()
	g.Add(value)
}

// AddValue adds values as dependencies, if the struct's fields
// or the function's input arguments needs them, they will be defined as
// bindings (at build-time) and they will be used (at serve-time).
func AddValue(val reflect.Value) {
	g := Global()
	g.AddValue(val)
}

// AddOnce binds a value to the controller's field with the same type,
// if it's not binded already.
//
// Returns false if binded already or the value is not the proper one for binding,
// otherwise true.
func AddOnce(value interface{}) bool {
	g := Global()
	return g.AddOnce(value)
}

// Remove unbinds a binding value based on the type,
// it returns true if at least one field is not binded anymore.
//
// The "n" indicates the number of elements to remove, if <=0 then it's 1,
// this is useful because you may have bind more than one value to two or more fields
// with the same type.
func Remove(value interface{}, n int) bool {
	g := Global()
	return g.Remove(value, n)
}

// Has returns true if a binder responsible to
// bind and return a type of "typ" is already registered to this controller.
func Has(value interface{}) bool {
	return Global().Has(value)
}

// Invoke calls constructor function and invokes it with
// injected values from container to constructor parameters
func Invoke(c interface{}) (interface{}, error) {
	g := Global()
	return g.Invoke(c)
}
