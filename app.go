package flow

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/go-flow/flow/di"
)

// New returns an App instance with default configuration.
func New() *App {
	cfg := Config{}
	return NewWithConfig(cfg)
}

// NewWithConfig creates new application instance
// with given configuration object
func NewWithConfig(cfg Config) *App {

	opts := NewOptions(cfg)
	return NewWithOptions(opts)
}

// NewWithOptions creates new application instance
// with given Application Options object
func NewWithOptions(opts Options) *App {
	// create application router
	r := NewRouter()

	if opts.UsePanicRecovery {
		r.Use(PanicRecovery())
	}

	if opts.UseRequestLogger {
		r.Use(RequestLogger())
	}

	if opts.ServeStatic {
		r.Static(opts.StaticPath, opts.StaticDir)
	}

	app := &App{
		Options:      opts,
		router:       r,
		dependencies: di.NewValues(),
	}

	//context pool allocation
	app.pool.New = func() interface{} {
		return app.allocateContext()
	}

	return app
}

// App -
type App struct {
	Options
	router *Router
	pool   sync.Pool

	methodNotAllowedHandler HandlerFunc
	notFoundHandler         HandlerFunc
	errorHandler            HandlerFunc

	dependencies di.Values
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

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (a *App) Any(relativePath string, handler HandlerFunc) {
	a.router.Any(relativePath, handler)
}

// Attach another router to current one
func (a *App) Attach(prefix string, router *Router) {
	a.router.Attach(prefix, router)
}

// Register appends one or more values as dependecies
func (a *App) Register(value interface{}) {
	if a.dependencies.Len() == 0 {
		a.dependencies.Add(value)
		return
	}

	// create injector
	injector := di.Struct(value, a.dependencies...)

	// inject dependencies to value
	injector.Inject(value)

	a.dependencies.Add(value)

}

// RegisterController registers application controller
func (a *App) RegisterController(ctrl Controller) {
	// check if controller imlements initer
	if i, ok := ctrl.(ControllerIniter); ok {
		i.Init(a)
	}

	// set prefix to default
	prefix := "/"
	// check if controller implements prefixer
	if p, ok := ctrl.(ControllerPrefixer); ok {
		prefix = p.Prefix()
	}

	//check if we have any dependencies registered
	if a.dependencies.Len() == 0 {
		// we dont have any dependencies defined
		a.router.Attach(prefix, ctrl.Routes())
		return
	}

	// get DI injector
	injector := di.Struct(ctrl, a.dependencies...)

	// inject dependencies to controller
	injector.Inject(ctrl)

	a.router.Attach(prefix, ctrl.Routes())
}

// MethodNotAllowedHandler is Handler where message and error can be personalized
// to be in line with application design and logic
func (a *App) MethodNotAllowedHandler(handler HandlerFunc) {
	a.methodNotAllowedHandler = handler
}

// NotFoundHandler is Handler where message and error can be personalized
// to be in line with application design and logic
func (a *App) NotFoundHandler(handler HandlerFunc) {
	a.notFoundHandler = handler
}

// ErrorHandler is Handler where message and error can be personalized
// to be in line with application design and logic
func (a *App) ErrorHandler(handler HandlerFunc) {
	a.errorHandler = handler
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

// Router returns application router instance
func (a *App) Router() *Router {
	return a.router
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

func (a *App) stop() error {
	return nil
}

// Stop issues interupt signal
func (a *App) Stop() error {
	// get current process
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	a.Logger.Debug("Stopping....")
	// issue interupt signal
	return proc.Signal(os.Interrupt)
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
			c.handlers = a.router.Handlers
			c.ServeError(http.StatusMethodNotAllowed, errors.New(a.Config.StringDefault("405Body", default405Body)))
			return
		}
	}

	c.handlers = a.router.Handlers
	c.ServeError(http.StatusNotFound, errors.New(a.Config.StringDefault("404Body", default404Body)))
}

func (a *App) allocateContext() *Context {
	return &Context{app: a}
}
