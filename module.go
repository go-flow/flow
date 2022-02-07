package flow

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-flow/flow/v3/di"
)

type ImportsProvider interface {
	// ProvideImports returns list of instance providers for module dependecies
	// This method is used to register all module dependecies
	// eg. logging, db connection,....
	// all dependecies that are provided in this method
	// will be available to all modules imported by the factory
	ProvideImports() []Provider
}

type ExportsProvider interface {

	// ProvideExports returns list of instance providers for
	// functionalities that module will export.
	// Exported functionalities will be available to other modules that
	// import module created by the Factory
	ProvideExports() []Provider
}

type ModulesProvider interface {

	// ProvideModules returns list of instance providers
	// for modules that current module depends on
	ProvideModules() []Provider
}

type Module struct {
	name      string
	factory   interface{}
	container di.Container
	parent    *Module
}

// NewModule creates new module instance
func NewModule(factory interface{}, parent *Module) (*Module, error) {
	if factory == nil {
		return nil, errors.New("ModuleFactory can not be nil")
	}

	typ := reflect.TypeOf(factory)

	name := typ.String()

	if typ.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("ModuleFactory `%s` has to be pointer", name)
	}

	module := &Module{
		factory: factory,
		name:    name[1:],
		parent:  parent,
	}

	// create module container depending if parent it is root module or not
	if parent != nil {
		// clone di container from parent module
		module.container = parent.container.Clone()
	} else {
		// use global container
		module.container = di.Clone()
	}

	if v, ok := factory.(ImportsProvider); ok {
		// populate module dependecies
		for _, p := range v.ProvideImports() {
			dep, err := p.Provide(&module.container)
			if err != nil {
				return nil, fmt.Errorf("unable to provide dependecy for module  `%s`. Error: %w", module.name, err)
			}
			module.container.Add(dep)
		}

	}

	if v, ok := factory.(ModulesProvider); ok {
		// register provided modules
		for _, p := range v.ProvideModules() {

			dep, err := p.Provide(&module.container)
			if err != nil {
				return nil, fmt.Errorf("unable to provide dependecy module for `%s` module. Error: %w", module.name, err)
			}

			mod, err := NewModule(dep, module)
			if err != nil {
				return nil, fmt.Errorf("unable to provide dependecy module for `%s` module. Error: %w", module.name, err)
			}

			// check if imported module exports any functionality
			if v, ok := dep.(ExportsProvider); ok {
				for _, p := range v.ProvideExports() {
					dep, err := p.Provide(&mod.container)
					if err != nil {
						return nil, fmt.Errorf("unable to provide exported dependecy for module  `%s`. Error: %w", mod.name, err)
					}

					// add exported dependecy to the module container
					mod.container.Add(dep)
					// add exported dependency to parent module
					module.container.Add(dep)
				}
			}

		}
	}

	if v, ok := factory.(interface{ SetModule(module *Module) }); ok {
		v.SetModule(module)
	}

	if v, ok := factory.(interface{ SetInvoker(in Invoker) }); ok {
		v.SetInvoker(&module.container)
	}

	return module, nil
}

// Factory returns factory value
func (m *Module) Factory() interface{} {
	return m.factory
}
