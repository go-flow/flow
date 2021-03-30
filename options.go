package flow

const (
	defaultEnv     = "development"
	defaultName    = "HiveApp"
	defaultAddr    = "0.0.0.0:5000"
	defaultVersion = "v0.0.0"

	defaultRedirectTrailingSlash  = true
	defaultRedirectFixedPath      = true
	defaultHandleMethodNotAllowed = true
	defaultHandleOptions          = true

	default404Body = "404 page not found"
	default405Body = "405 method not allowed"
)

// Options holds application configuration Options
type Options struct {
	RouterOptions
	initialized      bool
	Env              string
	Name             string
	Addr             string
	Version          string
	ModuleSuffix     string
	ControllerSuffix string
	ControllerIndex  string
}

// RouterOptions holds router configuration Options
type RouterOptions struct {
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	SaveMatchedRoutePath   bool
	HandleOptions          bool
	Body404                string
	Body405                string
}

// NewOptions creates New application Options instance
func NewOptions() Options {
	opts := Options{
		RouterOptions: RouterOptions{
			RedirectTrailingSlash:  defaultRedirectTrailingSlash,
			RedirectFixedPath:      defaultRedirectFixedPath,
			HandleMethodNotAllowed: defaultHandleMethodNotAllowed,
			HandleOptions:          defaultHandleOptions,
			Body404:                default404Body,
			Body405:                default405Body,
		},
		initialized: true,
		Env:         defaultEnv,
		Name:        defaultName,
		Addr:        defaultAddr,
		Version:     defaultVersion,
	}

	return opts
}
