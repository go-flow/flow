package web

import (
	"errors"
	"net/http"

	"github.com/go-flow/flow/v3"
)

// HandlerFunc is a function that is registered to a route to handle http requests
type HandlerFunc func(r *http.Request) Response

// Serve serves HTTP traffic for given Module
func Serve(root *flow.Module) error {

	app := &App{
		root: root,
	}

	f := root.Factory()

	//check if app implements Configer interface
	if c, ok := f.(Configer); ok {
		app.Config = c.Config()
	} else {
		app.Config = DefaultConfig()
	}

	return errors.New("not implemented")
}
