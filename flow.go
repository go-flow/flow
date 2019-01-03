package flow

import "github.com/go-flow/flow/config"

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

// New creates new application instance
// New loads application from configuration file path provided
//
// function will panic if configuration can not be loaded
// and converted into map[string]interface{} type
func New(configFile string) *App {
	//load application config from file
	cfgData := map[string]interface{}{}
	err := config.LoadFromPath(configFile, cfgData)
	if err != nil {
		panic(err)
	}

	return NewWithConfig(cfgData)
}

// Default returns an App instance with default configuration.
func Default() *App {
	cfg := map[string]interface{}{}
	return NewWithConfig(cfg)
}

// NewWithConfig creates new application instance
// with given configuration object
func NewWithConfig(data map[string]interface{}) *App {

	// ensure we have minimum configuration set
	configWithDefaults(data)

	//create configuration object
	cfg := NewConfig(data)

	// create application router
	r := NewRouter()

	app := &App{
		config: cfg,
		router: r,
	}
	app.pool.New = func() interface{} {
		return app.allocateContext()
	}

	return app
}

func configWithDefaults(data map[string]interface{}) {
	// ensure we have env set
	if _, found := data["env"]; !found {
		data["env"] = defaultEnv
	}

	// ensure we have name set
	if _, found := data["name"]; !found {
		data["name"] = defaultName
	}

	// ensure we have addr set
	if _, found := data["addr"]; !found {
		data["addr"] = defaultAddr
	}

	// ensure we have LogLevel set
	if _, found := data["logLevel"]; !found {
		data["logLevel"] = defaultLogLevel
	}

	// ensure we have redirectTrailingSlash set
	if _, found := data["redirectTrailingSlash"]; !found {
		data["redirectTrailingSlash"] = defaultRedirectTrailingSlash
	}

	// ensure we have redirectFixedPath set
	if _, found := data["redirectFixedPath"]; !found {
		data["redirectFixedPath"] = defaultRedirectFixedPath
	}

	// ensure we have maxMultipartMemory set
	if _, found := data["maxMultipartMemory"]; !found {
		data["maxMultipartMemory"] = defaultMultipartMemory
	}

	if _, found := data["secureJSONPrefix"]; !found {
		data["secureJSONPrefix"] = defaultSecureJSONPrefix
	}

	if _, found := data["handleMethodNotAllowed"]; !found {
		data["handleMethodNotAllowed"] = defaultHandleMethodNotAllowed
	}

	if _, found := data["404Body"]; !found {
		data["404Body"] = default404Body
	}

	if _, found := data["405Body"]; !found {
		data["405Body"] = default405Body
	}

}
