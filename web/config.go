package web

const (
	defaultEnv     = "development"
	defaultName    = "FlowWebApp"
	defaultAddr    = "0.0.0.0:5000"
	defaultVersion = "v0.0.0"

	defaultRedirectTrailingSlash  = true
	defaultRedirectFixedPath      = true
	defaultHandleMethodNotAllowed = true
	defaultHandleOptions          = true

	default404Body = "404 page not found"
	default405Body = "405 method not allowed"
)

//Configer is implemented by any value that has Config method which returns flow Web Config object
type Configer interface {
	Config() Config
}

// Config holds application configuration
type Config struct {
	RouterConfig
	Env              string
	Name             string
	Addr             string
	Version          string
	ModuleSuffix     string
	ControllerSuffix string
	ControllerIndex  string
}

// RouterConfig holds router configuration Options
type RouterConfig struct {
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	SaveMatchedRoutePath   bool
	HandleOptions          bool
	Body404                string
	Body405                string
}

// DefaultConfig creates New application Config instance
func DefaultConfig() Config {
	opts := Config{
		RouterConfig: RouterConfig{
			RedirectTrailingSlash:  defaultRedirectTrailingSlash,
			RedirectFixedPath:      defaultRedirectFixedPath,
			HandleMethodNotAllowed: defaultHandleMethodNotAllowed,
			HandleOptions:          defaultHandleOptions,
			Body404:                default404Body,
			Body405:                default405Body,
		},
		Env:     defaultEnv,
		Name:    defaultName,
		Addr:    defaultAddr,
		Version: defaultVersion,
	}

	return opts
}
