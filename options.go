package flow

import (
	"html/template"

	"github.com/go-flow/flow/log"
	"github.com/go-flow/flow/render/view"
	"github.com/go-flow/flow/sessions"
)

const (
	defaultEnv  = "development"
	defaultName = "FlowApp"
	defaultAddr = "0.0.0.0:3000"

	defaultLogLevel = "debug"

	defaultRedirectTrailingSlash  = true
	defaultRedirectFixedPath      = false
	defaultHandleMethodNotAllowed = true
	defaultMaxMultipartMemory     = 32 << 20 // 32 MB

	default404Body = "404 page not found"
	default405Body = "405 method not allowed"

	defaultUseSession  = true
	defaultSessionName = "_flow_app_session"

	defaultUseTranslator         = true
	defaultTranslatorLocalesRoot = "locales"
	defaultTranslatorDefaultLang = "en-US"

	defaultUseRequestLogger = true
	defaultUsePanicRecovery = true

	defaultUseViewEngine     = true
	defaultViewsRoot         = "views"
	defaultViewsExt          = ".tpl"
	defaultViewsMasterLayout = "layouts/master"
	defaultViewsPartialsRoot = "partials"
	defaultViewsDisableCache = false

	defaultServeStatic = true
	defaultStaticPath  = "/static"
	defaultStaticDir   = "./public"
)

// Options holds flow configuration options
type Options struct {
	Env  string
	Name string
	Addr string

	LogLevel string

	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	MaxMultipartMemory     int64

	Body404 string
	Body500 string

	UseSession    bool
	SessionName   string
	SessionSecret string

	UseTranslator         bool
	TranslatorLocalesRoot string
	TranslatorDefaultLang string

	UseRequestLogger bool
	UsePanicRecovery bool

	UseViewEngine     bool
	ViewsRoot         string
	ViewsExt          string
	ViewsMasterLayout string
	ViewsPartialsRoot string
	ViewsDisableCache bool

	ServeStatic bool
	StaticPath  string
	StaticDir   string

	Logger       log.Logger
	SessionStore sessions.Store
	ViewEngine   view.Engine
	Translator   *Translator

	AppConfig interface{}
}

// NewOptions returns a new Options instance with default configuration
func NewOptions() Options {
	opts := Options{
		Env:                    defaultEnv,
		Name:                   defaultName,
		Addr:                   defaultAddr,
		LogLevel:               defaultLogLevel,
		RedirectTrailingSlash:  defaultRedirectTrailingSlash,
		RedirectFixedPath:      defaultRedirectFixedPath,
		HandleMethodNotAllowed: defaultHandleMethodNotAllowed,
		MaxMultipartMemory:     defaultMaxMultipartMemory,
		Body404:                default404Body,
		Body500:                default405Body,
		UseSession:             defaultUseSession,
		SessionName:            defaultSessionName,
		UseTranslator:          defaultUseTranslator,
		TranslatorLocalesRoot:  defaultTranslatorLocalesRoot,
		TranslatorDefaultLang:  defaultTranslatorDefaultLang,
		UseRequestLogger:       defaultUseRequestLogger,
		UsePanicRecovery:       defaultUsePanicRecovery,
		UseViewEngine:          defaultUseViewEngine,
		ViewsRoot:              defaultViewsRoot,
		ViewsExt:               defaultViewsExt,
		ViewsMasterLayout:      defaultViewsMasterLayout,
		ViewsPartialsRoot:      defaultViewsPartialsRoot,
		ViewsDisableCache:      defaultViewsDisableCache,
		ServeStatic:            defaultServeStatic,
		StaticPath:             defaultStaticPath,
		StaticDir:              defaultStaticDir,
	}

	return opts
}

func optionsWithDefault(opts Options) Options {
	//configure logger
	if opts.Logger == nil {
		opts.Logger = log.New(opts.LogLevel, opts.Env)
	}

	//configure session store
	if opts.UseSession && opts.SessionStore == nil {
		if opts.SessionSecret == "" {
			opts.Logger.Warn("SessionSecret configuration key is not set. Your sessions are not safe!")
		}
		opts.SessionStore = sessions.NewCookieStore([]byte(opts.SessionSecret))
	}
	//configure ViewEngine
	if opts.UseViewEngine && opts.ViewEngine == nil {
		partials, err := loadPartials(opts.ViewsRoot, opts.ViewsPartialsRoot, opts.ViewsExt)
		if err != nil {
			opts.Logger.Fatal(err.Error())
		}
		opts.ViewEngine = view.NewHTMLEngine(view.Config{
			Root:         opts.ViewsRoot,
			Ext:          opts.ViewsExt,
			Master:       opts.ViewsMasterLayout,
			Partials:     partials,
			Funcs:        make(template.FuncMap),
			DisableCache: opts.ViewsDisableCache,
			Delims:       view.Delims{Left: "{{", Right: "}}"},
		})
	}

	// configure translator
	if opts.UseTranslator && opts.Translator == nil {
		t, err := NewTranslator(opts.TranslatorLocalesRoot, opts.TranslatorDefaultLang)
		if err != nil {
			opts.Logger.Fatal(err.Error())
		}
		opts.Translator = t
	}

	return opts
}
