package flow

import (
	"net/http"
	"sync"
)

// New returns an App instance with default configuration.
func New() *App {
	cfg := map[string]interface{}{}
	return NewWithConfig(cfg)
}

// NewWithConfig creates new application instance
// with given configuration object
func NewWithConfig(data map[string]interface{}) *App {

	opts := NewOptions(data)

	// create application router
	r := NewRouter()

	app := &App{
		Options: opts,
		router:  r,
	}
	app.pool.New = func() interface{} {
		return app.allocateContext()
	}

	return app
}

// App -
type App struct {
	Options

	router *Router
	//logger
	//sessions
	//render engine
	pool sync.Pool
}

// ServeHTTP conforms to the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := a.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = r
	c.reset()

	a.handleHTTPRequest(c)
	a.pool.Put(c)
}

// HandleContext re-enter a context that has been rewritten.
// This can be done by setting c.Request.URL.Path to your new target.
// Disclaimer: You can loop yourself to death with this, use wisely.
func (a *App) HandleContext(c *Context) {
	c.reset()
	a.handleHTTPRequest(c)
}

func (a *App) handleHTTPRequest(c *Context) {
	req := c.Request
	httpMethod := req.Method
	path := req.URL.Path

	if root := a.router.trees[httpMethod]; root != nil {
		if handlers, ps, tsr := root.getValue(path); handlers != nil {
			c.handlers = handlers
			c.Params = ps
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		} else if httpMethod != "CONNECT" && path != "/" {
			code := http.StatusMovedPermanently // Permanent redirect, request with GET method
			if httpMethod != "GET" {
				code = http.StatusTemporaryRedirect
			}
			if tsr && a.RedirectTrailingSlash {
				req.URL.Path = path + "/"
				if length := len(path); length > 1 && path[length-1] == '/' {
					req.URL.Path = path[:length-1]
				}
				// logger here
				http.Redirect(c.Response, req, req.URL.String(), code)
				c.writermem.WriteHeaderNow()
				return
			}

			if a.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(CleanPath(path), a.RedirectTrailingSlash)
				if found {
					req.URL.Path = string(fixedPath)
					// logger here
					http.Redirect(c.Response, req, req.URL.String(), code)
					c.writermem.WriteHeaderNow()
					return
				}
			}
		}
	}

	if a.HandleMethodNotAllowed {
		if allow := a.router.allowed(path, httpMethod); len(allow) > 0 {
			c.handlers = a.router.Middlewares
			c.ServeError(http.StatusMethodNotAllowed, []byte(a.AppConfig.GetStringD("405Body", default405Body)))
			return
		}
	}

	c.handlers = a.router.Middlewares
	c.ServeError(http.StatusNotFound, []byte(a.AppConfig.GetStringD("404Body", default404Body)))
}

func (a *App) allocateContext() *Context {
	return &Context{app: a}
}
