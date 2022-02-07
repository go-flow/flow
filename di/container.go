package di

import (
	"errors"
	"fmt"
	"reflect"
)

// Container is a shortcut for []reflect.Value
type Container []reflect.Value

// NewContainer creates Empty Container
func NewContainer() Container {
	return Container{}
}

// Clone returns a copy of the current container value
func (c Container) Clone() Container {
	if n := len(c); n > 0 {
		values := make(Container, n)
		copy(values, c)
		return values
	}
	return NewContainer()
}

// Len returns Length of the current Container slice
func (c Container) Len() int {
	return len(c)
}

// Add adds values as dependencies, if the struct's fields
// or the function's input arguments needs them, they will be defined as
// bindings (at build-time) and they will be used (at serve-time).
func (c *Container) Add(value interface{}) {
	c.AddValue(ValueOf(value))
}

// AddValue adds values as dependencies, if the struct's fields
// or the function's input arguments needs them, they will be defined as
// bindings (at build-time) and they will be used (at serve-time).
func (c *Container) AddValue(val reflect.Value) {
	if !goodVal(val) {
		return
	}
	*c = append(*c, val)
}

// AddOnce binds a value to the controller's field with the same type,
// if it's not binded already.
//
// Returns false if binded already or the value is not the proper one for binding,
// otherwise true.
func (c *Container) AddOnce(value interface{}) bool {
	return c.addIfNotExists(reflect.ValueOf(value))
}

// Remove unbinds a binding value based on the type,
// it returns true if at least one field is not binded anymore.
//
// The "n" indicates the number of elements to remove, if <=0 then it's 1,
// this is useful because you may have bind more than one value to two or more fields
// with the same type.
func (c *Container) Remove(value interface{}, n int) bool {
	return c.remove(reflect.TypeOf(value), n)
}

// Has returns true if a binder responsible to
// bind and return a type of "typ" is already registered to this controller.
func (c Container) Has(value interface{}) bool {
	return c.valueTypeExists(reflect.TypeOf(value))
}

// Invoke calls constructor function and invokes it with
// injected values from container to constructor parameters
func (c *Container) Invoke(constructor interface{}) (interface{}, error) {

	typ := reflect.TypeOf(constructor)
	if typ == nil {
		return nil, errors.New("can not provide an untyped nil")
	}

	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("constructor function not provided, got %v type(%v)", constructor, typ.Kind())
	}

	in := make([]reflect.Value, typ.NumIn())

	for i := 0; i < typ.NumIn(); i++ {
		t := typ.In(i)

		if v, ok := c.getTypeVal(t); ok {
			in[i] = v
		} else {
			return nil, fmt.Errorf("unable to provide constructor, parameter %s is nill for constructor %#v", t.Name(), constructor)
		}
	}

	val := reflect.ValueOf(constructor).Call(in)

	return val[0].Interface(), nil
}

func (c *Container) remove(typ reflect.Type, n int) (ok bool) {
	input := *c
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

	*c = input

	return
}

func (c Container) valueTypeExists(typ reflect.Type) bool {
	for _, in := range c {
		if equalTypes(in.Type(), typ) {
			return true
		}
	}
	return false
}

func (c Container) getTypeVal(typ reflect.Type) (reflect.Value, bool) {
	for _, in := range c {
		if equalTypes(in.Type(), typ) {
			return in, true
		}
	}
	return reflect.Value{}, false
}

func (c *Container) addIfNotExists(val reflect.Value) bool {
	var (
		typ = val.Type() // no element, raw things here.
	)

	if !goodVal(val) {
		return false
	}

	if c.valueTypeExists(typ) {
		return false
	}

	c.Add(val)
	return true
}
