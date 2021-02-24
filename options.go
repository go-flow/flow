package flow

const (
	defaultEnv     = "development"
	defaultName    = "HiveApp"
	defaultAddr    = "0.0.0.0:5000"
	defaultVersion = "v0.0.0"

	defaultRedirectTrailingSlash  = true
	defaultRedirectFixedPath      = true
	defaultHandleMethodNotAllowed = true

	default404Body = "404 page not found"
	default405Body = "405 method not allowed"
)

type Options struct {
	RouterOptions
	Env     string
	Name    string
	Addr    string
	Version string

	Body404 string
	Body405 string
}

type RouterOptions struct {
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	SaveMatchedRoutePath   bool
}

func NewOptions() Options {
	opts := Options{
		RouterOptions: RouterOptions{
			RedirectTrailingSlash:  defaultRedirectTrailingSlash,
			RedirectFixedPath:      defaultRedirectFixedPath,
			HandleMethodNotAllowed: defaultHandleMethodNotAllowed,
		},

		Env:     defaultEnv,
		Name:    defaultName,
		Addr:    defaultAddr,
		Version: defaultVersion,

		Body404: default404Body,
		Body405: default405Body,
	}

	return opts
}
