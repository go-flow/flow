package flow

import "github.com/go-flow/flow/sessions"

const (
	defaultEnv                    = "development"
	defaultName                   = "FlowApp"
	defaultAddr                   = "0.0.0.0:3000"
	defaultLogLevel               = "debug"
	defaultRedirectTrailingSlash  = true
	defaultRedirectFixedPath      = false
	defaultHandleMethodNotAllowed = true
	defaultMultipartMemory        = 32 << 20 // 32 MB
	defaultSecureJSONPrefix       = "while(1);"
	default404Body                = "404 page not found"
	default405Body                = "405 method not allowed"
)

// Options are used to configure and define how your application should run.
type Options struct {
	// Name is application name
	Name string

	// Addr is the bind address provided to http.Server. Default is "127.0.0.1:3000"
	// Can be set using ENV vars "ADDR" and "PORT".
	Addr string

	// Env is the "environment" in which the App is running. Default is "development".
	Env string

	// LogLevel defaults to "debug".
	LogLevel string

	// Logger to be used with the application. A default one is provided.
	Logger Logger

	// SessionStore is used to back the session.
	SesionStore sessions.Store

	// SessionName is the name of the session cookie that is set.
	SessionName string

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// Value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64

	// Application speccific configuration object
	AppConfig Config
}

// NewOptions returns a new Options instance with sensible defaults
func NewOptions(data Config) Options {
	return optionsWithDefaults(data)
}

func optionsWithDefaults(cfg Config) Options {
	opts := Options{}

	opts.Env = cfg.StringDefault("env", defaultEnv)

	opts.Name = cfg.StringDefault("name", defaultName)

	opts.Addr = cfg.StringDefault("addr", defaultAddr)

	opts.LogLevel = cfg.StringDefault("logLevel", defaultLogLevel)

	opts.RedirectTrailingSlash = cfg.BoolDefault("redirectTrailingSlash", defaultRedirectTrailingSlash)

	opts.RedirectFixedPath = cfg.BoolDefault("redirectFixedPath", defaultRedirectFixedPath)

	opts.MaxMultipartMemory = cfg.Int64Default("maxMultipartMemory", defaultMultipartMemory)

	opts.HandleMethodNotAllowed = cfg.BoolDefault("handleMethodNotAllowed", defaultHandleMethodNotAllowed)

	if _, found := cfg["404Body"]; !found {
		cfg["404Body"] = default404Body
	}

	if _, found := cfg["405Body"]; !found {
		cfg["405Body"] = default405Body
	}

	opts.AppConfig = cfg
	return opts
}
