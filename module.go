package flow

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-flow/flow/v3/di"
	"github.com/go-flow/flow/v3/web"
)

type Module struct {
	name      string
	factory   interface{}
	container di.Container
	parent    *Module
	modules   []*Module
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

			module.modules = append(module.modules, mod)
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

func (m *Module) registerRouters(parent *web.Router) error {
	rp, ok := m.factory.(RouterProvider)
	if !ok {
		return nil
	}

	for _, provider := range rp.ProvideRouters() {
		val, err := provider.Provide(&m.container)
		if err != nil {
			return fmt.Errorf("unable to register router provider for module  `%s`. Error: %w", m.name, err)
		}

		rf, ok := val.(RouterFactory)
		if !ok {
			return fmt.Errorf("unable to register router provider for module `%s`. Error: %w", m.name, errors.New("provided constructor did not create instance of RouterFactory interface"))
		}

		router := parent.Group(rf.Path(), rf.Middlewares()...)

		for _, p := range rf.ProvideHandlers() {
			val, err = p.Provide(&m.container)
			if err != nil {
				return fmt.Errorf("unable to register action handler for module  `%s`. Error: %w", m.name, err)
			}

			handler, ok := val.(ActionHandler)
			if !ok {
				return fmt.Errorf("unable to register action handler for module  `%s`. Error: %w", m.name, errors.New("provided constructor did not create instance of ActionHandler interface"))
			}

			router.Handle(handler.Method(), handler.Path(), handler.Handle, handler.Middlewares()...)
		}

		// check if sub routers should be registered for given router
		if rf.RegisterSubRouters() {
			for _, module := range m.modules {
				if err := module.registerRouters(router); err != nil {
					return fmt.Errorf("unable to register routers of imported module `%s`. Error: %w", module.name, err)
				}
			}
		}
	}

	return nil

}
