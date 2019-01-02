// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package flow

import (
	"math"
	"reflect"
	"runtime"
)

const abortIndex int8 = math.MaxInt8 / 2

// Mux is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Mux struct {
	trees map[string]*node

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

	// Middlewares represents list of middlewares that will be executed in chain
	Middlewares HandlersChain
}

var _ Router = NewMux()

// NewMux returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func NewMux() *Mux {
	return &Mux{
		RedirectTrailingSlash: true,
		RedirectFixedPath:     true,
	}
}

// Use appends one or more middlewares onto the Router stack.
func (m *Mux) Use(middleware ...HandlerFunc) {
	m.Middlewares = append(m.Middlewares, middleware...)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (m *Mux) GET(path string, handle HandlerFunc) {
	m.Handle("GET", path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (m *Mux) HEAD(path string, handle HandlerFunc) {
	m.Handle("HEAD", path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (m *Mux) OPTIONS(path string, handle HandlerFunc) {
	m.Handle("OPTIONS", path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (m *Mux) POST(path string, handle HandlerFunc) {
	m.Handle("POST", path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (m *Mux) PUT(path string, handle HandlerFunc) {
	m.Handle("PUT", path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (m *Mux) PATCH(path string, handle HandlerFunc) {
	m.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (m *Mux) DELETE(path string, handle HandlerFunc) {
	m.Handle("DELETE", path, handle)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (m *Mux) Handle(method, path string, handle HandlerFunc) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if m.trees == nil {
		m.trees = make(map[string]*node)
	}

	root := m.trees[method]
	if root == nil {
		root = new(node)
		m.trees[method] = root
	}

	chained := m.prepareChainHandler(handle)

	root.addRoute(path, chained)
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (m *Mux) Lookup(method, path string) (HandlersChain, Params, bool) {
	if root := m.trees[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (m *Mux) Routes() (routes Routes) {
	for method, tree := range m.trees {
		routes = iterate("", method, routes, tree)
	}
	return routes
}

func (m *Mux) prepareChainHandler(handler HandlerFunc) HandlersChain {
	finalSize := len(m.Middlewares) + 1
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, m.Middlewares)
	mergedHandlers = append(mergedHandlers, handler)
	return mergedHandlers
}

func (m *Mux) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range m.trees {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range m.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _, _ := m.trees[method].getValue(path)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

func iterate(path, method string, routes Routes, root *node) Routes {
	path += root.path
	if len(root.handler) > 0 {
		handlerFunc := root.handler.Last()
		routes = append(routes, Route{
			Method:      method,
			Path:        path,
			Handler:     nameOfFunction(handlerFunc),
			HandlerFunc: handlerFunc,
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
