package di

import (
	"errors"
	"fmt"
	"reflect"
)

// Container is a shortcut for []reflect.Value
type Container []reflect.Value

// NewContainer returns new empty Container
func NewContainer() Container {
	return Container{}
}

// Clone returns a copy of the current value
func (c Container) Clone() Container {
	if n := len(c); n > 0 {
		values := make(Container, n, n)
		copy(values, c)
		return values
	}
	return NewContainer()
}

// CloneWithFieldsOf will return a copy of the current container
// with provided struct fields that are filled(non-zero) by the caller
func (c Container) CloneWithFieldsOf(i interface{}) Container {
	values := c.Clone()

	// add the manual filled fields to the dependencies.
	filledFieldValues := LookupNonZeroFieldsValues(ValueOf(i), true)
	values = append(values, filledFieldValues...)
	return values
}

// Len returns Length of current Container slice
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

// Remove unbinds a binding value based on the type,
// it returns true if at least one field is not binded anymore.
//
// The "n" indicates the number of elements to remove, if <=0 then it's 1,
// this is useful because you may have bind more than one value to two or more fields
// with the same type.
func (c *Container) Remove(value interface{}, n int) bool {
	return c.remove(reflect.TypeOf(value), n)
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

// Has returns true if a binder responsible to
// bind and return a type of "typ" is already registered to this controller.
func (c Container) Has(value interface{}) bool {
	return c.valueTypeExists(reflect.TypeOf(value))
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

// AddOnce binds a value to the controller's field with the same type,
// if it's not binded already.
//
// Returns false if binded already or the value is not the proper one for binding,
// otherwise true.
func (c *Container) AddOnce(value interface{}) bool {
	return c.addIfNotExists(reflect.ValueOf(value))
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

// Register appends one or more values as dependecies
func (c *Container) Register(value interface{}) {
	if c.Len() == 0 {
		c.Add(value)
		return
	}

	// create injector
	injector := Struct(value, *c...)

	// inject dependencies to value
	injector.Inject(value)

	c.AddOnce(value)
}

// InjectDeps accepts a destination struct and any optional context value(s),
// and injects registered dependencies to the destination object
func (c *Container) InjectDeps(dest interface{}, ctx ...reflect.Value) {
	injector := Struct(dest, *c...)
	injector.Inject(dest, ctx...)
}

// Provide registers constructor function and
// invokes it with injected values from container to constructor
func (c *Container) Provide(constructor interface{}) (interface{}, error) {

	typ := reflect.TypeOf(constructor)
	if typ == nil {
		return nil, errors.New("can't provide an untyped nil")
	}

	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("must provide constructor function, got %v (type %v)", constructor, typ)
	}

	in := make([]reflect.Value, typ.NumIn())

	for i := 0; i < typ.NumIn(); i++ {
		t := typ.In(i)
		if v, ok := c.getTypeVal(t); ok {
			in[i] = v
		} else {
			return nil, fmt.Errorf("Unable to provide constructor, parameter %s is nil for constructor %#v", t.Name(), constructor)
		}
	}
	val := reflect.ValueOf(constructor).Call(in)

	return val[0].Interface(), nil
}

// ProvideAndRegister registers constructor function and
// creates instances and register them to container
func (c *Container) ProvideAndRegister(constructor interface{}) error {
	v, err := c.Provide(constructor)
	if err != nil {
		return err
	}

	c.Add(v)
	return nil
}
