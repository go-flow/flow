package flow

import (
	"html/template"
	"io/ioutil"
	"path"
	"strings"

	"github.com/go-flow/flow/sessions"
	"github.com/go-flow/flow/view"
)

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
	defaultSessionSecret          = "c8OMa61enGu9Nt1fS13RkmUz17W7SRt8"
	defaultUseRequestLogger       = true
	defaultUsePanicRecovery       = true
	defaultViewsRoot              = "views"
	defaultViewsExt               = ".tpl"
	defaultViewsMasterLayout      = "layouts/master"
	defaultViewsPartialsRoot      = "partials"
	defaultServeStatic            = true
	defaultStaticPath             = "/static"
	defaultStaticDir              = "./public"
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

	// LogFormat determines in which format our logs will be
	//
	// Default is `text`
	LogFormat string

	// Logger to be used with the application. A default one is provided.
	Logger Logger

	// SessionStore is used to back the session.
	SessionStore sessions.Store

	//ViewEngine is used to render HTML
	ViewEngine *view.Engine

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
	Config Config

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
	opts.LogFormat = cfg.StringDefault("logFormat", defaultLogFormat)

	opts.RedirectTrailingSlash = cfg.BoolDefault("redirectTrailingSlash", defaultRedirectTrailingSlash)

	opts.RedirectFixedPath = cfg.BoolDefault("redirectFixedPath", defaultRedirectFixedPath)

	opts.MaxMultipartMemory = cfg.Int64Default("maxMultipartMemory", defaultMultipartMemory)

	opts.HandleMethodNotAllowed = cfg.BoolDefault("handleMethodNotAllowed", defaultHandleMethodNotAllowed)

	opts.UseRequestLogger = cfg.BoolDefault("useRequestLogger", defaultUseRequestLogger)

	opts.UsePanicRecovery = cfg.BoolDefault("usePanicRecovery", defaultUsePanicRecovery)

	opts.ServeStatic = cfg.BoolDefault("serveStatic", defaultServeStatic)
	opts.StaticPath = cfg.StringDefault("staticPath", defaultStaticPath)
	opts.StaticDir = cfg.StringDefault("staticDir", defaultStaticDir)

	if opts.Logger == nil && cfg.BoolDefault("useLogger", defaultUseLogger) == true {
		opts.Logger = NewLoggerWithFormatter(opts.LogLevel, opts.LogFormat)
	}

	viewsRoot := cfg.StringDefault("viewsRoot", defaultViewsRoot)
	partialsRoot := cfg.StringDefault("viewsPartialsRoot", defaultViewsPartialsRoot)
	ext := cfg.StringDefault("viewsExt", defaultViewsExt)
	partials, err := loadPartials(viewsRoot, partialsRoot, ext)
	if err != nil {
		if opts.Logger != nil {
			opts.Logger.Error(err)
		} else {
			panic(err)
		}

		return opts
	}

	if opts.ViewEngine == nil && cfg.BoolDefault("useViewEngine", defaultUseViewEngine) == true {
		opts.ViewEngine = view.New(view.Config{
			Root:         viewsRoot,
			Extension:    ext,
			Master:       cfg.StringDefault("viewsMasterLayout", defaultViewsMasterLayout),
			Partials:     partials,
			Funcs:        make(template.FuncMap),
			DisableCache: cfg.BoolDefault("viewsDisableCache", opts.Env == "development"),
			Delims:       view.Delims{Left: "{{", Right: "}}"},
		})
	}

	if opts.SessionStore == nil && cfg.BoolDefault("useSession", defaultUseSession) {
		secret := cfg.String("sessionSecret")
		if secret == "" {
			if opts.Env != defaultEnv && opts.Logger != nil {
				opts.Logger.Warn("SessionSecret configuration key is not set. Your sessions are not safe!")
			} else {
				secret = defaultSessionSecret
			}
		}
		opts.SessionStore = sessions.NewCookieStore([]byte(secret))
	}

	if _, found := cfg["404Body"]; !found {
		cfg["404Body"] = default404Body
	}

	if _, found := cfg["405Body"]; !found {
		cfg["405Body"] = default405Body
	}

	opts.Config = cfg
	return opts
}

func loadPartials(viewsRoot, partialsRoot, ext string) ([]string, error) {
	dirname := path.Join(viewsRoot, partialsRoot)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	partials := []string{}
	for _, f := range files {
		partial := f.Name()
		if strings.HasSuffix(partial, ext) {
			// remove ext from file
			partial = strings.TrimRight(partial, ext)

			// join file with folder name
			partial = path.Join(partialsRoot, partial)

			// add to partials
			partials = append(partials, partial)
		}
	}
	return partials, nil
}
