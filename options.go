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

type Options struct {
	RouterOptions
	Env     string
	Name    string
	Addr    string
	Version string
}

type RouterOptions struct {
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	SaveMatchedRoutePath   bool
	HandleOptions          bool
	Body404                string
	Body405                string
}

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

		Env:     defaultEnv,
		Name:    defaultName,
		Addr:    defaultAddr,
		Version: defaultVersion,
	}

	return opts
}
