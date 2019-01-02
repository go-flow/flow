package flow

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-flow/flow/defaults"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

// Options type is  used to configure and
// define how the application should run.
type Options struct {
	Name string

	// Addr is the bind address provided to http.Server. Default is "127.0.0.1:3000"
	// Can be set using ENV vars "ADDR" and "PORT".
	Addr string

	// Host that this application will be available at.
	// Default is "http://127.0.0.1:[$PORT|3000]".
	Host string

	// Env is the "environment" in which the App is running.
	// Default is "development".
	Env string

	// LogLevel defaults to "debug".
	LogLevel string

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

	// Value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64
}

// NewOptions returns a new Options instance with sensible defaults
func NewOptions() Options {
	return optionsWithDefaults(Options{
		RedirectTrailingSlash: true,
		RedirectFixedPath:     false,
		MaxMultipartMemory:    defaultMultipartMemory,
	})
}

func optionsWithDefaults(opts Options) Options {

	opts.Env = defaults.String(opts.Env, GetEnv("GO_ENV", "development"))
	opts.LogLevel = defaults.String(opts.LogLevel, GetEnv("LOG_LEVEL", "debug"))
	opts.Name = defaults.String(opts.Name, "/")
	addr := "0.0.0.0"
	if opts.Env == "development" {
		addr = "127.0.0.1"
	}
	envAddr := GetEnv("ADDR", addr)

	if strings.HasPrefix(envAddr, "unix:") {
		// UNIX domain socket doesn't have a port
		opts.Addr = envAddr
	} else {
		// TCP case
		opts.Addr = defaults.String(opts.Addr, fmt.Sprintf("%s:%s", envAddr, GetEnv("PORT", "3000")))
	}

	opts.Host = defaults.String(opts.Host, GetEnv("HOST", fmt.Sprintf("http://127.0.0.1:%s", GetEnv("PORT", "3000"))))

	return opts
}

// GetEnv returns environment variable value for a given key
// if value is not found defaultValue param will be returned
func GetEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue

}
