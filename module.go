package flow

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/go-flow/flow/di"
)

// Module struct
type Module struct {
	factory   interface{}
	options   Options
	name      string
	path      string
	container di.Container
	parent    *Module
	router    *Router
	imports   []*Module
}

// NewModule creates new Module object
func NewModule(factory interface{}, container di.Container, parent *Module) (*Module, error) {
	if factory == nil {
		return nil, fmt.Errorf("factory object can not bi nil")
	}

	// get module factory type
	typ := reflect.TypeOf(factory)
	// get module name
	name := typ.String()

	// factory object has to be pointer
	if typ.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("Module factory object %s has to be pointer", name)
	}

	// create new module instance
	module := &Module{
		factory:   factory,
		name:      name[1:],
		container: container,
		parent:    parent,
	}

	// register all providers to module container
	if v, ok := factory.(ModuleProvider); ok {
		for _, provider := range v.Providers() {
			if err := module.container.ProvideAndRegister(provider); err != nil {
				return nil, fmt.Errorf("Unable to register provider for module  `%s`. Error: %v", module.name, err)
			}
		}
	}

	// register all dependecies (imported modules)
	if v, ok := factory.(ModuleImporter); ok {
		for _, dep := range v.Imports() {
			m, err := NewModule(dep, module.container.Clone(), module)
			if err != nil {
				return nil, fmt.Errorf("Unable to Import dependecy module `%s`. Error: %v", m.name, err)
			}

			// check if imported module exports any functionality
			if val, ok := dep.(ModuleExporter); ok {
				for _, exp := range val.Exports() {
					module.container.Register(exp)
				}
			}

			module.imports = append(module.imports, m)
		}
	}

	// initialize module
	if v, ok := factory.(ModuleIniter); ok {
		//inject dependecies to module factory
		// only if it is root module
		if module.parent == nil {
			module.container.InjectDeps(factory)
		}

		if err := v.Init(); err != nil {
			return nil, fmt.Errorf("Unable to initialize module %s. error : %v", module.name, err)
		}
	}

	// get module options
	if v, ok := factory.(ModuleOptioner); ok {
		module.options = v.Options()
	} else if module.parent == nil && !module.options.initialized {
		// ensure options are initialized at least for root module
		module.options = NewOptions()
	} else if module.parent != nil {
		// inherit parent options
		module.options = module.parent.options
	}

	return module, nil
}

// Serve the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully.
func (m *Module) Serve() error {
	r, err := m.Router()
	if err != nil {
		return err
	}

	// create http server
	srv := http.Server{
		Handler: r,
	}

	// make interrupt signal channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// listen for interrupt signals
	go func() {
		<-c
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(fmt.Errorf("Unable to gracefully shutdown HTTP.Server, error: %w", err))
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

// Path returns module path
func (m *Module) Path() string {
	if m.path != "" {
		return m.path
	}

	if m.parent == nil {
		// root module by default will have / path
		m.path = "/"
	} else {
		// try to guess module path based on the module name
		m.path = strings.Split(m.name, ".")[1]
		m.path = strings.TrimSuffix(m.path, m.options.ModuleSuffix)
		m.path = toSnakeCase(m.path)
		m.path = strings.ToLower(m.path)
		m.path = fmt.Sprintf("/%s", m.path)
	}

	// override module path if ModulePather interface is implemented
	if val, ok := m.factory.(ModulePather); ok {
		m.path = val.Path()
	}

	return m.path
}

// Router gets router, if router is not initialized router is initialized
func (m *Module) Router() (*Router, error) {
	if m.router != nil {
		return m.router, nil
	}
	// check if module provides router
	if v, ok := m.factory.(ModuleRouter); ok {
		m.router = v.Router()
	}

	// ensure default router
	if m.router == nil {
		m.router = NewRouterWithOptions(m.options.RouterOptions)
	}

	//register controllers for current module
	if err := m.registerControllers("", m.router); err != nil {
		return nil, err
	}

	// register controllers from all imported modules
	for _, module := range m.imports {
		if err := module.registerControllers(m.Path(), m.router); err != nil {
			return nil, err
		}
	}

	return m.router, nil
}

func (m *Module) registerControllers(parent string, r *Router) error {
	if v, ok := m.factory.(ModuleController); ok {
		for _, ctrlP := range v.Controllers() {

			ctrl, err := m.container.Provide(ctrlP)
			if err != nil {
				return fmt.Errorf("module %s can not invoke controller constructor %v", m.name, ctrlP)
			}

			// get controller type
			typ := reflect.TypeOf(ctrl)
			//get controller name
			name := typ.String()

			// check if controller is pointer
			if typ.Kind() != reflect.Ptr {
				return fmt.Errorf("controller %s in module %s has to be pointer", name, m.name)
			}
			// remove * from controller name
			name = name[1:]

			// check if controller follows naming convention
			if !strings.HasSuffix(name, m.options.ControllerSuffix) {
				return fmt.Errorf("controller %s in module %s does not follow naming convention", name, m.name)
			}

			// initialize controller
			if val, ok := ctrl.(ModuleIniter); ok {
				if err := val.Init(); err != nil {
					return fmt.Errorf("unable to initialize controller %s in module %s: %w", name, m.name, err)
				}
			}

			// define controller path
			path := "/"
			ctrlName := strings.Split(name, ".")[1]

			if ctrlName != m.options.ControllerIndex {
				path = toSnakeCase(ctrlName)
				path = fmt.Sprintf("/%s", path)
				path = strings.ToLower(path)
			}

			// assign custom path is controller implements Pather interface
			if val, ok := ctrl.(ModulePather); ok {
				path = val.Path()
			}

			if !strings.HasPrefix(path, "/") {
				return fmt.Errorf("unable to register controller %s in module %s: controller path has to start with `/` ", name, m.name)
			}

			if val, ok := ctrl.(ControllerRouter); ok {
				rPath := fmt.Sprintf("%s%s", parent, m.Path())
				gRouter := r.Group(rPath)
				val.Routes(gRouter)
			} else {
				return fmt.Errorf("unable to register controller %s in module %s: controller does not implement ControllerRouter interface", name, m.name)
			}
		}

	}
	return nil
}
