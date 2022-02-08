package flow

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-flow/flow/v3/web"
)

//Configer is implemented by any value that has Config method which returns flow Web Config object
type Configer interface {
	Config() web.Config
}

// RouterFactory interface responsible for creating module routers
type RouterFactory interface {
	Path() string
	Middlewares() []web.MiddlewareHandlerFunc
	ProvideHandlers() []Provider
	RegisterSubRouters() bool
}

// ActionHandler interface is used to define http action handlers defined by module router
type ActionHandler interface {
	Method() string
	Path() string
	Middlewares() []web.MiddlewareHandlerFunc
	Handle(r *http.Request, ps web.Params) web.Response
}

type WebApp struct {
	web.Config
	root   *Module
	router *web.Router
}

// BuildWebApp Creates web Application instance
func BuildWebApp(factory interface{}) (*WebApp, error) {

	root, err := Bootstrap(factory)
	if err != nil {
		return nil, err
	}

	app := &WebApp{
		root: root,
	}

	f := root.Factory()

	//check if app implements Configer interface
	if c, ok := f.(Configer); ok {
		app.Config = c.Config()
	} else {
		app.Config = web.DefaultConfig()
	}

	app.router = web.NewRouterWithOptions(app.RouterConfig)

	if err := root.registerRouters(app.router); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *WebApp) Serve() error {
	if a.router == nil {
		return fmt.Errorf("unable to serve app. Error: %w", errors.New("http router is not initialized"))
	}

	if s, ok := a.root.factory.(Starter); ok {
		if err := s.Start(); err != nil {
			return fmt.Errorf("unable to start module `%s`. Error: %w", a.root.name, err)
		}
	}

	// create http server
	srv := http.Server{
		Handler: a.router,
	}

	// make interrupt signal channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// listen for interrupt signals
	go func() {
		<-c

		if s, ok := a.root.factory.(Stopper); ok {
			s.Stop()
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			panic(fmt.Errorf("unable to gracefully shutdown HTTP.Server. Error: %w", err))
		}
	}()

	// get listen address from options
	addr := a.Addr

	if strings.HasPrefix(addr, "unix:") {
		// create unix network listener
		lis, err := net.Listen(addr, addr[5:])
		if err != nil {
			return err
		}
		// start accepting incomming requests on listener
		return srv.Serve(lis)
	}

	// assign address to http server
	srv.Addr = addr

	//start accepting incomming requests
	return srv.ListenAndServe()
}
