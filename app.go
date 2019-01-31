package flow

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"

	"github.com/go-flow/flow/di"
)

const (
	// ControllerPackage holds package name in which controllers can be registered
	ControllerPackage = "controllers"

	// ControllerIndex holds controller Index name
	ControllerIndex = "Index"

	// ControllerSuffix holds controller naming convention
	ControllerSuffix = "Controller"
)

// App holds fully working application setup
type App struct {
	Options

	router *Router
	pool   sync.Pool

	methodNotAllowedHandler HandlerFunc
	unauthorizedHandler     HandlerFunc
	notFoundHandler         HandlerFunc
	errorHandler            HandlerFunc

	container di.Container
}

// New returns an App instance with default configuration.
func New() *App {
	return NewWithOptions(NewOptions())
}

// NewWithOptions creates new application instance
// with given Application Options object
func NewWithOptions(opts Options) *App {

	opts = optionsWithDefault(opts)

	// create application router
	r := NewRouter()

	if opts.UseRequestLogger {
		r.Use(RequestLogger())
	}

	if opts.UsePanicRecovery {
		r.Use(PanicRecovery())
	}

	if opts.ServeStatic {
		r.Static(opts.StaticPath, opts.StaticDir)
	}

	app := &App{
		Options:   opts,
		router:    r,
		container: di.NewContainer(),
	}

	//context pool allocation
	app.pool.New = func() interface{} {
		return app.allocateContext()
	}

	return app
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
	if a.container.Len() == 0 {
		a.container.Add(value)
		return
	}

	// create injector
	injector := di.Struct(value, a.container...)

	// inject dependencies to value
	injector.Inject(value)

	a.container.Add(value)
}

// InjectDeps accepts a destination struct and any optional context value(s),
// and injects registered dependencies to the destination object
func (a *App) InjectDeps(dest interface{}, ctx ...reflect.Value) {
	injector := di.Struct(dest, a.container...)
	injector.Inject(dest, ctx...)
}

// RegisterController registers application controller
func (a *App) RegisterController(ctrl interface{}) {

	// set controller route prefix to default
	prefix := "/"

	// check naming convention
	typ := reflect.TypeOf(ctrl)

	// get full controller full name
	fullCtrlName := typ.String()

	// check if controller is pointer
	if typ.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("Controller `%s` has to be pointer", fullCtrlName))
	}
	// remove * from full name
	fullCtrlName = fullCtrlName[1:]

	// check if passed controller is in proper package
	if !strings.HasPrefix(fullCtrlName, ControllerPackage) {
		panic(fmt.Sprintf("Controller `%s` has to be in `%s` package", fullCtrlName, ControllerPackage))
	}

	//check if pased controller follows naming conventions
	if !strings.HasSuffix(fullCtrlName, ControllerSuffix) {
		panic(fmt.Sprintf("Controller `%s` does not follow naming convention", fullCtrlName))
	}

	// extract controller name from struct

	ctrlName := strings.Replace(fullCtrlName, ".", "", -1)
	ctrlName = strings.TrimPrefix(ctrlName, ControllerPackage)
	ctrlName = strings.TrimSuffix(ctrlName, ControllerSuffix)

	// assign controller Name to prefix if it is not Index controller
	if ctrlName != ControllerIndex {
		prefix = fmt.Sprintf("/%s", ctrlName)
		prefix = strings.ToLower(prefix)
	}

	// check if controller implements prefixer
	if p, ok := ctrl.(ControllerPrefixer); ok {
		prefix = p.Prefix()
	}

	// check if controller imlements initer
	if i, ok := ctrl.(ControllerIniter); ok {
		i.Init(a)
	}

	// log registration for debugging purposes
	a.Logger.Debug(fmt.Sprintf("Registering `%s` with Prefix: `%s`\n", fullCtrlName, prefix))

	ctrlRouter, ok := ctrl.(ControllerRouter)
	if !ok {
		panic(fmt.Sprintf("controller `%s` does not implement ControllerRouter interface", fullCtrlName))
	}

	routes := ctrlRouter.Routes()

	//check if we have any dependencies registered
	if a.container.Len() == 0 {
		// we dont have any dependencies defined
		a.router.Attach(prefix, routes)
		return
	}

	// get DI injector
	injector := di.Struct(ctrl, a.container...)

	// inject dependencies to controller
	injector.Inject(ctrl)

	a.router.Attach(prefix, routes)
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

// UnauthorizedHandler is handler which is triggered ServeError with 401 status code is called
func (a *App) UnauthorizedHandler(handler HandlerFunc) {
	a.unauthorizedHandler = handler
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
	a.Logger.Info(fmt.Sprintf("Starting Application at %s", a.Addr))
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
			a.Logger.Error(err.Error())
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			a.Logger.Error(err.Error())
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
			if a.methodNotAllowedHandler != nil {
				c.handlers = append(c.handlers, a.methodNotAllowedHandler)
				c.Next()
				return
			}
			c.ServeError(http.StatusMethodNotAllowed, errors.New(default405Body))
			return
		}
	}

	c.handlers = a.router.Handlers

	if a.notFoundHandler != nil {
		c.handlers = append(c.handlers, a.notFoundHandler)
		c.Next()
		return
	}

	c.ServeError(http.StatusNotFound, errors.New(default404Body))
}

func (a *App) allocateContext() *Context {
	return &Context{app: a}
}
