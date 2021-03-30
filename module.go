package flow

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/go-flow/flow/v2/di"
)

// Module struct
type Module struct {
	factory   ModuleFactory
	name      string
	options   Options
	container di.Container
	parent    *Module
	modules   []*Module
	router    *Router
}

// NewModule creates new Module object
func NewModule(factory ModuleFactory, container di.Container, parent *Module) (*Module, error) {
	if factory == nil {
		return nil, fmt.Errorf("factory object can not be nil")
	}

	// get module factory type
	typ := reflect.TypeOf(factory)
	// get module name
	name := typ.String()

	// factory object has to be pointer
	if typ.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("module factory object %s has to be pointer", name)
	}

	// create new module instance
	module := &Module{
		factory:   factory,
		name:      name[1:],
		container: container,
		parent:    parent,
	}

	// root module provides dependecies for all child modules.
	if parent == nil {
		// register all imports to module container
		for _, p := range factory.ProvideImports() {
			obj, err := p.Provide(&module.container)
			if err != nil {
				return nil, fmt.Errorf("unable to provide dependecy for module  `%s`. Error: %w", module.name, err)
			}
			module.container.Add(obj)
		}
	}

	// register all dependecies (imported modules)
	for _, p := range factory.ProvideModules() {

		// provide module factory object
		dep, err := p.Provide(&module.container)
		if err != nil {
			return nil, fmt.Errorf("unable to provide dependecy module for module `%s`. Error: %w", module.name, err)
		}

		depFac, ok := dep.(ModuleFactory)
		if !ok {
			return nil, fmt.Errorf("unable to provide dependecy module for module `%s`. Error: %w", module.name, errors.New("provided constructor did not create instance of ModuleFactory interface"))
		}
		// create module object
		m, err := NewModule(depFac, module.container.Clone(), module)
		if err != nil {
			return nil, fmt.Errorf("unable to provide dependecy module for module `%s`. Error: %w", module.name, err)
		}

		// check if imported module exports any functionality
		for _, p := range depFac.ProvideExports() {
			exp, err := p.Provide(&m.container)
			if err != nil {
				return nil, fmt.Errorf("unable to provide exported dependecy for module `%s`. Error: %w", m.name, err)
			}
			// add feature to module
			m.container.Add(exp)
			// add feature to parent module
			module.container.Add(exp)
		}

		module.modules = append(module.modules, m)
	}

	// import all dependecies for child modules
	// child modules first register their child modules,
	// and then they provide functionality internally which can depend on child modules
	if parent != nil {
		for _, p := range factory.ProvideImports() {
			dep, err := p.Provide(&module.container)
			if err != nil {
				return nil, fmt.Errorf("unable to provide dependecy for module  `%s`. Error: %w", module.name, err)
			}

			module.container.Add(dep)
		}
	}

	if module.IsRoot() {
		module.container.InjectDeps(factory)

		// get module options
		if v, ok := factory.(ModuleOptioner); ok {
			module.options = v.Options()
		} else if !module.options.initialized {
			// ensure options are initialized at least for root module
			module.options = NewOptions()
		} else if module.parent != nil {
			// inherit parent options
			module.options = module.parent.options
		}

		module.router = NewRouterWithOptions(module.options.RouterOptions)
		if err := module.registerRouters(module.router); err != nil {
			return nil, fmt.Errorf("unable to register routers for module `%s`. Error: %w", module.name, err)
		}
	}

	return module, nil
}

func (m *Module) registerRouters(parent *Router) error {

	// initialize module routers
	for _, p := range m.factory.ProvideRouters() {
		// get router provider
		rp, err := p.Provide(&m.container)
		if err != nil {
			return fmt.Errorf("unable to register router provider for module  `%s`. Error: %w", m.name, err)
		}

		// check if provided object is RouterFactory interface
		rf, ok := rp.(RouterFactory)
		if !ok {
			return fmt.Errorf("unable to register router provider for module `%s`. Error: %w", m.name, errors.New("provided constructor did not create instance of RouterFactory interface"))
		}

		group := parent.Group(rf.Path(), rf.Middlewares()...)

		// for root modules create root routers with shared tree
		if m.IsRoot() {
			group.parent = nil
			group.root = true
		}

		// get all action handlers
		for _, p := range rf.ProvideHandlers() {
			// provide action handler
			ah, err := p.Provide(&m.container)
			if err != nil {
				return fmt.Errorf("unable to register action handler for module  `%s`. Error: %w", m.name, err)
			}

			// check if provided object is ActionHandler interface
			handler, ok := ah.(ActionHandler)
			if !ok {
				return fmt.Errorf("unable to register action handler for module  `%s`. Error: %w", m.name, errors.New("provided constructor did not create instance of ActionHandler interface"))
			}

			group.Handle(handler.Method(), handler.Path(), handler.Handle, handler.Middlewares()...)
		}

		// check if sub routers should be registered for given router
		if rf.RegisterSubRouters() {
			for _, module := range m.modules {
				if err := module.registerRouters(group); err != nil {
					return fmt.Errorf("unable to register routers of imported module `%s`. Error: %w", module.name, err)
				}
			}
		}

	}
	return nil
}

func (m *Module) IsRoot() bool {
	return m.parent == nil
}

// Serve the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully.
func (m *Module) Serve() error {

	if m.router == nil {
		return fmt.Errorf("unable to serve module `%s`. Error: %w", m.name, errors.New("http router is not initialized"))
	}

	if s, ok := m.factory.(ModuleStarter); ok {
		if err := s.Start(); err != nil {
			return fmt.Errorf("unable to start module `%s`. Error: %w", m.name, err)
		}
	}

	// create http server
	srv := http.Server{
		Handler: m.router,
	}

	// make interrupt signal channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// listen for interrupt signals
	go func() {
		<-c

		if s, ok := m.factory.(ModuleStopper); ok {
			s.Stop()
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			panic(fmt.Errorf("unable to gracefully shutdown HTTP.Server. Error: %w", err))
		}
	}()

	// get listen address from options
	addr := m.options.Addr

	if strings.HasPrefix(addr, "unix:") {
		// create unix network listener
		lis, err := net.Listen(addr, addr[5:])
		if err != nil {
			return err
		}
		// start accepting incomming requests on listener
		return srv.Serve(lis)
	}

	// assign address to http server
	srv.Addr = addr

	//start accepting incomming requests
	return srv.ListenAndServe()

}
