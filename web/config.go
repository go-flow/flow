package web

const (
	defaultEnv     = "development"
	defaultName    = "FlowWebApp"
	defaultAddr    = "0.0.0.0:5000"
	defaultVersion = "v0.0.0"

	defaultHandleOptions = true

	default404Body = "404 page not found"
)

// Config holds application configuration
type Config struct {
	RouterConfig
	Env     string
	Name    string
	Addr    string
	Version string
}

// RouterConfig holds router configuration Options
type RouterConfig struct {
	HandleOptions bool
	Body404       string
}

// DefaultConfig creates New application Config instance
func DefaultConfig() Config {
	opts := Config{
		RouterConfig: RouterConfig{
			HandleOptions: defaultHandleOptions,
			Body404:       default404Body,
		},
		Env:     defaultEnv,
		Name:    defaultName,
		Addr:    defaultAddr,
		Version: defaultVersion,
	}

	return opts
}
