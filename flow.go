package flow

import (
	"net/http"

	"github.com/go-flow/flow/di"
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEYAML              = "application/x-yaml"
)

// HandlerFunc is a function that is registered to a route to handle http requests
type HandlerFunc func(r *http.Request) Response

// ModuleProvider is interface used for injecting dependecies that can be used in the module
type ModuleProvider interface {
	Providers() []interface{}
}

// ModulePather is interface used to define http path for Module
type ModulePather interface {
	Path() string
}

// ModuleImporter interface is used to for providing list of imported modules
// that export providers that arerequired in this module
type ModuleImporter interface {
	Imports() []interface{}
}

// ModuleExporter interface is used for exporting functionalities to modules that import the module
type ModuleExporter interface {
	Exports() []interface{}
}

// ModuleIniter is interface used for Module Initialization
type ModuleIniter interface {
	Init() error
}

// ModuleOptioner is interface used for providing Application Options
// This interface is used only for root module or AppModule
type ModuleOptioner interface {
	Options() Options
}

// ModuleRouter interface alows module to define custom Routing
type ModuleRouter interface {
	Router() *Router
}

// ModuleController interface allows module to define Controllers
type ModuleController interface {
	Controllers() []interface{}
}

// ModuleMiddleware interface allows module to define their routing middlewares
type ModuleMiddleware interface {
	Middlewares() []MiddlewareHandlerFunc
}

// ControllerRouter interface allows controllers to define their routing logic
type ControllerRouter interface {
	Routes(*Router)
}

// Bootstrap creates Flow Module instance for given factory object
func Bootstrap(moduleFactory interface{}) (*Module, error) {
	rootModule, err := NewModule(moduleFactory, di.NewContainer(), nil)

	return rootModule, err
}
