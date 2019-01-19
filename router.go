// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package flow

import (
	"fmt"
	"math"
	"net/http"
	"path"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
//
// Router is associated with a prefix and an array of handlers(middlewares)
type Router struct {
	// routing tree nodes
	trees map[string]*node

	// Handlers represents list of middlewares that will be executed in chain
	Handlers HandlersChain

	// base path for router
	basePath string

	// root determines if the router is root router
	root bool
}

// NewRouter returns a new initialized Router.
func NewRouter() *Router {
	return &Router{
		root:     true,
		basePath: "/",
		trees:    make(map[string]*node),
		Handlers: nil,
	}
}

// BasePath returns the base path of router group.
func (r *Router) BasePath() string {
	return r.basePath
}

// Group creates a new router group.
//
// You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (r *Router) Group(relativePath string, handlers ...HandlerFunc) *Router {
	return &Router{
		root:     false,
		basePath: r.calculateAbsolutePath(relativePath),
		trees:    r.trees,
		Handlers: r.combineHandlers(handlers),
	}
}

// Use appends one or more middlewares onto the Router stack.
func (r *Router) Use(middleware ...HandlerFunc) {
	r.Handlers = append(r.Handlers, middleware...)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handlers ...HandlerFunc) {

	path = r.calculateAbsolutePath(path)
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if method == "" {
		panic("HTTP method can not be empty")
	}

	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	if r.trees == nil {
		panic("Router tree not initialized")
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	chained := r.combineHandlers(handlers)

	root.addRoute(path, chained)
}

// GET is a shortcut for router.Handle("GET", path, handler)
func (r *Router) GET(path string, handler HandlerFunc) {
	r.Handle("GET", path, handler)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handler)
func (r *Router) HEAD(path string, handler HandlerFunc) {
	r.Handle("HEAD", path, handler)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handler)
func (r *Router) OPTIONS(path string, handler HandlerFunc) {
	r.Handle("OPTIONS", path, handler)
}

// POST is a shortcut for router.Handle("POST", path, handler)
func (r *Router) POST(path string, handler HandlerFunc) {
	r.Handle("POST", path, handler)
}

// PUT is a shortcut for router.Handle("PUT", path, handler)
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.Handle("PUT", path, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handler)
func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.Handle("PATCH", path, handler)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.Handle("DELETE", path, handler)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (r *Router) Any(relativePath string, handler HandlerFunc) {
	r.Handle("GET", relativePath, handler)
	r.Handle("POST", relativePath, handler)
	r.Handle("PUT", relativePath, handler)
	r.Handle("PATCH", relativePath, handler)
	r.Handle("HEAD", relativePath, handler)
	r.Handle("OPTIONS", relativePath, handler)
	r.Handle("DELETE", relativePath, handler)
	r.Handle("CONNECT", relativePath, handler)
	r.Handle("TRACE", relativePath, handler)
}

// Attach another router to current one
func (r *Router) Attach(prefix string, router *Router) {

	for _, route := range router.Routes() {
		path := joinPaths(prefix, route.Path)
		r.Handle(route.Method, path, route.HandlersChain...)
	}
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
//
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func (r *Router) StaticFile(relativePath, filePath string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := func(c *Context) {
		c.File(filePath)
	}

	r.GET(relativePath, handler)
	r.HEAD(relativePath, handler)
}

// StaticFS serves files from the given file system root with a custom `http.FileSystem` can be used instead.
func (r *Router) StaticFS(relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	handler := r.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "/*filepath")

	r.GET(urlPattern, handler)
	r.HEAD(urlPattern, handler)
}

// Static serves files from the given file system root.
//
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
//
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (r *Router) Static(relativePath, root string) {
	r.StaticFS(relativePath, Dir(root, false))
}

// Lookup allows the manual lookup of a method + path combo.
//
// If the path was found, it returns the handler chain and the path parameter
// values. The third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (HandlersChain, Params, bool) {
	if root := r.trees[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

// Routes returns a slice of registered routes
func (r *Router) Routes() (routes Routes) {
	for method, tree := range r.trees {
		routes = iterate("", method, routes, tree)
	}
	return routes
}

func (r *Router) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := r.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	// create handler
	handler := func(c *Context) {
		if _, nolisting := fs.(onlyFilesFS); nolisting {
			c.Response.WriteHeader(http.StatusNotFound)
		}

		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			fmt.Println(err)
			c.ServeError(http.StatusNotFound, []byte(c.app.Config.StringDefault("404Body", default404Body)))
			return
		}

		fileServer.ServeHTTP(c.Response, c.Request)

	}

	return handler
}

func (r *Router) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(r.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, r.Handlers)
	copy(mergedHandlers[len(r.Handlers):], handlers)
	return mergedHandlers
}

func (r *Router) calculateAbsolutePath(relativePath string) string {
	return joinPaths(r.basePath, relativePath)
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range r.trees {
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
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path)
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
