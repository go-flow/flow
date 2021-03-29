package flow

import (
	"net/http"

	"github.com/go-flow/flow/v2/di"
)

// HandlerFunc is a function that is registered to a route to handle http requests
type HandlerFunc func(r *http.Request) Response

// Injector defines Dependency Injector interface
type Injector interface {
	Provide(constructor interface{}) (interface{}, error)
}

// ActionHandler interface is used to define http action handlers defined by module router
type ActionHandler interface {
	Method() string
	Path() string
	Middlewares() []MiddlewareHandlerFunc
	Handle(r *http.Request) Response
}

// ModuleFactory interface for creating flow.Module
type ModuleFactory interface {

	// ProvideImports returns list of instance providers for module dependecies
	// This method is used to register all module dependecies
	// eg. logging, db connection,....
	// all dependecies that are provided in this method
	// will be available to all modules imported by the factory
	ProvideImports() []Provider

	// ProvideExports returns list of instance providers for
	// functionalities that module will export.
	// Exported functionalities will be available to other modules that
	// import module created by the Factory
	ProvideExports() []Provider

	// ProvideModules returns list of instance providers
	// for modules that current module depends on
	ProvideModules() []Provider

	// ProvideRouters returns list of instance providers for module routers.
	// Module routers are used for http routing
	ProvideRouters() []Provider
}

// ModuleOptioner interface is used for providing Application Options
// This interface is used only for root module or AppModule
type ModuleOptioner interface {
	Options() Options
}

// RouterFactory interface responsible for creating module routers
type RouterFactory interface {
	Path() string
	Middlewares() []MiddlewareHandlerFunc
	ProvideHandlers() []Provider
	RegisterSubRouters() bool
}

// ModuleStarter interface used when http application is served
// Start method is invoked if module implements the interface
type ModuleStarter interface {
	Start() error
}

// ModuleStopper interface used when http application is stopped
// Stop method is invoked during shutdown process if module implements the interface
type ModuleStopper interface {
	Stop()
}

// Bootstrap creates Flow Module instance for given factory object
func Bootstrap(moduleFactory ModuleFactory) (*Module, error) {
	rootModule, err := NewModule(moduleFactory, di.NewContainer(), nil)
	return rootModule, err
}
