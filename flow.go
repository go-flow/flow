package flow

import "github.com/go-flow/flow/config"

const (
	defaultEnv                   = "development"
	defaultName                  = "FlowApp"
	defaultAddr                  = "0.0.0.0:3000"
	defaultLogLevel              = "debug"
	defaultRedirectTrailingSlash = true
	defaultRedirectFixedPath     = false
	defaultMultipartMemory       = 32 << 20 // 32 MB
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

// NewWithConfig creates new application instance
// with given configuration object
func NewWithConfig(data map[string]interface{}) *App {

	// ensure we have minimum configuration set
	configWithDefaults(data)

	cfg := NewConfig(data)

	r := NewMux()
	r.RedirectTrailingSlash = cfg.GetBool("redirectTrailingSlash")
	r.RedirectFixedPath = cfg.GetBool("redirectFixedPath")

	return nil
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

}
