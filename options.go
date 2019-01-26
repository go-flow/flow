package flow

const (
	defaultEnv  = "development"
	defaultName = "FlowApp"
	defaultAddr = "0.0.0.0:3000"

	defaultLogLevel  = "debug"
	defaultLogFormat = "text"

	defaultRedirectTrailingSlash  = true
	defaultRedirectFixedPath      = false
	defaultHandleMethodNotAllowed = true
	defaultMaxMultipartMemory     = 32 << 20 // 32 MB

	default404Body = "404 page not found"
	default405Body = "405 method not allowed"

	defaultUseSession  = true
	defaultSessionName = "_flow_app_session"

	defaultUseI18n         = true
	defaultI18nLocalesRoot = "locales"
	defaultI18nDefaultLang = "en-US"

	defaultUseRequestLogger = true
	defaultUsePanicRecovery = true

	defaultUseViewEngine     = true
	defaultViewsRoot         = "views"
	defaultViewsExt          = ".tpl"
	defaultViewsMasterLayout = "layouts/master"
	defaultViewsPartialsRoot = "partials"

	defaultServeStatic = true
	defaultStaticPath  = "/static"
	defaultStaticDir   = "./public"
)

// Options holds flow configuration options
type Options struct {
	Env  string
	Name string
	Addr string

	LogLevel  string
	LogFormat string

	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	MaxMultipartMemory     int64

	Body404 string
	Body500 string

	UseSession  bool
	SessionName string

	UseI18n         bool
	I18nLocalesRoot string
	I18nDefaultLang string

	UseRequestLogger bool
	UsePanicRecovery bool

	UseViewEngine     bool
	ViewsRoot         string
	ViewsExt          string
	ViewsMasterLayout string
	ViewsPartialsRoot string

	ServeStatic bool
	StaticPath  string
	StaticDir   string
}

// NewOptions returns a new Options instance with default configuration
func NewOptions() Options {
	opts := Options{
		Env:                    defaultEnv,
		Name:                   defaultName,
		Addr:                   defaultAddr,
		LogLevel:               defaultLogLevel,
		LogFormat:              defaultLogFormat,
		RedirectTrailingSlash:  defaultRedirectTrailingSlash,
		RedirectFixedPath:      defaultRedirectFixedPath,
		HandleMethodNotAllowed: defaultHandleMethodNotAllowed,
		MaxMultipartMemory:     defaultMaxMultipartMemory,
		Body404:                default404Body,
		Body500:                default405Body,
		UseSession:             defaultUseSession,
		SessionName:            defaultSessionName,
		UseI18n:                defaultUseI18n,
		I18nLocalesRoot:        defaultI18nLocalesRoot,
		I18nDefaultLang:        defaultI18nDefaultLang,
		UseRequestLogger:       defaultUseRequestLogger,
		UsePanicRecovery:       defaultUsePanicRecovery,
		UseViewEngine:          defaultUseViewEngine,
		ViewsRoot:              defaultViewsRoot,
		ViewsExt:               defaultViewsExt,
		ViewsMasterLayout:      defaultViewsMasterLayout,
		ViewsPartialsRoot:      defaultViewsPartialsRoot,
		ServeStatic:            defaultServeStatic,
		StaticPath:             defaultStaticPath,
		StaticDir:              defaultStaticDir,
	}

	return opts
}
