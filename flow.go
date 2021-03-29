package flow

import (
	"net/http"

	"github.com/go-flow/flow/di"
)

// HandlerFunc is a function that is registered to a route to handle http requests
type HandlerFunc func(r *http.Request) Response

// Injector defines Dependency Injector interface
type Injector interface {
	Provide(constructor interface{}) (interface{}, error)
}

// Importer interface is used for importing module dependecies
type Importer interface {
	Imports() []Provider
}

// Exporter interface is used for exporting functionalities to modules that import the module
type Exporter interface {
	Exports() []Provider
}

// Moduler interface is used for importing modules as depedecies
type Moduler interface {
	Modules() []Provider
}

// Initer interface is used for Module Initialization
type Initer interface {
	Init() error
}

// Pather interface is used to define http path
type Pather interface {
	Path() string
}

// Middlewarer interface defines their routing middlewares
type Middlewarer interface {
	Middlewares() []MiddlewareHandlerFunc
}

// RouterProvider interface is used to define module http routers
type RouterProvider interface {
	Routers() []Provider
}

// ModuleIncluder interface is used on routers to determine
// if sub module routing should be included in routing
type ModuleIncluder interface {
	IncludeChildModules() bool
}

// ActionHandlerer interface is used to define http action handlers defined by module router
type ActionHandlerer interface {
	ActionHandlers() []Provider
}

// ActionHandler interface is used to define http action handlers defined by module router
type ActionHandler interface {
	Method() string
	Path() string
	Middlewares() []MiddlewareHandlerFunc
	Handle(r *http.Request) Response
}

// Bootstrap creates Flow Module instance for given factory object
func Bootstrap(moduleFactory ModuleFactory) (*Module, error) {
	rootModule, err := NewModule(moduleFactory, di.NewContainer(), nil)
	return rootModule, err
}
