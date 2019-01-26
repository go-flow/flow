package flow

const (
	defaultEnv                    = "development"
	defaultName                   = "FlowApp"
	defaultAddr                   = "0.0.0.0:3000"
	defaultLogLevel               = "debug"
	defaultLogFormat              = "text"
	defaultRedirectTrailingSlash  = true
	defaultRedirectFixedPath      = false
	defaultHandleMethodNotAllowed = true
	defaultMultipartMemory        = 32 << 20 // 32 MB
	default404Body                = "404 page not found"
	default405Body                = "405 method not allowed"
	defaultUseLogger              = true
	defaultUseViewEngine          = true
	defaultUseSession             = true
	defaultSessionName            = "_flow_app_session"
	defaultUseTranslator          = true
	defaultTranslatorLocalesRoot  = "locales"
	defaultTranslatorDefaultLang  = "en-US"
	defaultUseRequestLogger       = true
	defaultUsePanicRecovery       = true
	defaultViewsRoot              = "views"
	defaultViewsExt               = ".tpl"
	defaultViewsMasterLayout      = "layouts/master"
	defaultViewsPartialsRoot      = "partials"
	defaultServeStatic            = true
	defaultStaticPath             = "/static"
	defaultStaticDir              = "./public"
	defaultUseI18n                = true
	defaultUseViewEngine          = true
)

// Options holds flow configuration options
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

	// LogFormat determines in which format our logs will be
	//
	// Default is `text`
	LogFormat string

	// UseSession determines if session management is enabled
	UseSession bool

	// SessionName is the name of the session cookie that is set.
	SessionName string

	// SessionSecret is used for encrypting sessions
	SessionSecret string

	// UseViewEngine determines if ViewEngine is enabled
	UseViewEngine bool

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

	// UseRequestLogger determines if RequestLogger will be used by default
	UseRequestLogger bool

	// UsePanicRecovery determines if PanicRecovery will be used by default
	UsePanicRecovery bool

	// ServeStatic determines if application will serve static files
	//
	// Default value is `true`
	ServeStatic bool

	//StaticPath defines on which path static content will be served
	//
	// Default value is `/`
	StaticPath string

	// StaticDir is path to static dir
	//
	// default value is `public`
	StaticDir string

	UseI18n bool
}

// NewOptions returns a new Options instance with default configuration
func NewOptions() Options {
	opts := Options{
		Name:                   defaultName,
		Addr:                   defaultAddr,
		Env:                    defaultEnv,
		LogLevel:               defaultLogLevel,
		LogFormat:              defaultLogFormat,
		RedirectTrailingSlash:  defaultRedirectTrailingSlash,
		RedirectFixedPath:      defaultRedirectFixedPath,
		MaxMultipartMemory:     defaultMaxMultipartMemory,
		HandleMethodNotAllowed: defaultHandleMethodNotAllowed,
		UseRequestLogger:       defaultUseRequestLogger,
		UsePanicRecovery:       defaultUsePanicRecovery,
		ServeStatic:            defaultServeStatic,
		StaticPath:             defaultStaticPath,
		StaticDir:              defaultStaticDir,
		SessionName:            defaultSessionName,
		UseI18n:                defaultUseI18n,
		UseViewEngine:          defaultUseViewEngine,
	}

	return opts
}
