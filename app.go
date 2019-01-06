package flow

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
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

	//initialize logger
	logger := NewLogger(opts.LogLevel)

	// create application router
	r := NewRouter()

	app := &App{
		Options: opts,
		Logger:  logger,
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
	Logger Logger

	router *Router

	//sessions
	//render engine
	pool sync.Pool
}

// Use appends one or more middlewares onto the Router stack.
func (a *App) Use(middleware ...HandlerFunc) {
	a.router.Use(middleware...)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (a *App) GET(path string, handler HandlerFunc) {
	a.router.GET(path, handler)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (a *App) HEAD(path string, handler HandlerFunc) {
	a.router.HEAD(path, handler)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (a *App) OPTIONS(path string, handler HandlerFunc) {
	a.router.OPTIONS(path, handler)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (a *App) POST(path string, handler HandlerFunc) {
	a.router.POST(path, handler)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (a *App) PUT(path string, handler HandlerFunc) {
	a.router.PUT(path, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (a *App) PATCH(path string, handler HandlerFunc) {
	a.router.PATCH(path, handler)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (a *App) DELETE(path string, handler HandlerFunc) {
	a.router.DELETE(path, handler)
}

// Serve the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully.
func (a *App) Serve() error {
	a.Logger.Infof("Starting Application at %s", a.Addr)

	// create http server
	srv := http.Server{
		Handler: a,
	}

	// make interupt channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt)
	// listen for interupt signal
	go func() {
		<-c
		a.Logger.Info("Shutting down application")
		if err := a.stop(); err != nil {
			a.Logger.Error(err)
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			a.Logger.Error(err)
		}
	}()

	if strings.HasPrefix(a.Addr, "unix:") {
		// create unix network listener
		lis, err := net.Listen("unix", a.Addr[5:])
		if err != nil {
			return err
		}
		// start accepting incomming requests on listener
		return srv.Serve(lis)
	}

	//assign address to server
	srv.Addr = a.Addr
	// start accepting incomming requests on listener
	return srv.ListenAndServe()

}

// ServeHTTP conforms to the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get context from pool
	c := a.pool.Get().(*Context)
	// reset response writer
	c.writermem.reset(w)
	// set request
	c.Request = r
	// reset context from previous use
	c.reset()

	// handle the request
	a.handleHTTPRequest(c)

	// put back context to pool
	a.pool.Put(c)
}

// HandleContext re-enter a context that has been rewritten.
// This can be done by setting c.Request.URL.Path to your new target.
// Disclaimer: You can loop yourself to death with this, use wisely.
func (a *App) HandleContext(c *Context) {
	c.reset()
	a.handleHTTPRequest(c)
}
func (a *App) stop() error {
	return nil
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
