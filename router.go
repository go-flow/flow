package flow

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	parent     *Router
	root       bool
	basePath   string
	trees      map[string]*node
	paramsPool sync.Pool
	maxParams  uint16
	mws        *MiddlewareStack

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 308 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 308 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOptions bool

	// Body404 string to be displayed when route is not found
	Body404 string

	// Body405 string to be displayed when route is not allowed
	Body405 string
}

// NewRouter creates new Router instance with default options
func NewRouter() *Router {
	opts := NewOptions()
	return NewRouterWithOptions(opts.RouterOptions)
}

// NewRouterWithOptions creates new Router instance for given options
func NewRouterWithOptions(opts RouterOptions) *Router {
	return &Router{
		root:                   true,
		basePath:               "/",
		mws:                    new(MiddlewareStack),
		trees:                  make(map[string]*node),
		RedirectTrailingSlash:  opts.RedirectTrailingSlash,
		RedirectFixedPath:      opts.RedirectFixedPath,
		HandleMethodNotAllowed: opts.HandleMethodNotAllowed,
		HandleOptions:          opts.HandleOptions,
		Body404:                opts.Body404,
		Body405:                opts.Body405,
	}
}

func (r *Router) getParams() *Params {
	if r.parent != nil {
		return r.parent.getParams()
	}

	ps, _ := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

func (r *Router) putParams(ps *Params) {
	if r.parent != nil {
		r.parent.putParams(ps)
		return
	}
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}

// GET is a shortcut for router.Handle(http.MethodGet, path, handler)
func (r *Router) GET(path string, handler HandlerFunc, middlewares ...MiddlewareHandlerFunc) {
	r.Handle(http.MethodGet, path, handler, middlewares...)
}

// HEAD is a shortcut for router.Handle(http.MethodHead, path, handler)
func (r *Router) HEAD(path string, handler HandlerFunc, middlewares ...MiddlewareHandlerFunc) {
	r.Handle(http.MethodHead, path, handler, middlewares...)
}

// OPTIONS is a shortcut for router.Handle(http.MethodOptions, path, handler)
func (r *Router) OPTIONS(path string, handler HandlerFunc, middlewares ...MiddlewareHandlerFunc) {
	r.Handle(http.MethodOptions, path, handler, middlewares...)
}

// POST is a shortcut for router.Handle(http.MethodPost, path, handler)
func (r *Router) POST(path string, handler HandlerFunc, middlewares ...MiddlewareHandlerFunc) {
	r.Handle(http.MethodPost, path, handler, middlewares...)
}

// PUT is a shortcut for router.Handle(http.MethodPut, path, handler)
func (r *Router) PUT(path string, handler HandlerFunc, middlewares ...MiddlewareHandlerFunc) {
	r.Handle(http.MethodPut, path, handler, middlewares...)
}

// PATCH is a shortcut for router.Handle(http.MethodPatch, path, handler)
func (r *Router) PATCH(path string, handler HandlerFunc, middlewares ...MiddlewareHandlerFunc) {
	r.Handle(http.MethodPatch, path, handler, middlewares...)
}

// DELETE is a shortcut for router.Handle(http.MethodDelete, path, handler)
func (r *Router) DELETE(path string, handler HandlerFunc, middlewares ...MiddlewareHandlerFunc) {
	r.Handle(http.MethodDelete, path, handler, middlewares...)
}

// Use appends one or more middlewares to middleware stack.
func (r *Router) Use(mw ...MiddlewareHandlerFunc) {
	r.mws.Append(mw...)
}

// Group creates a new router group.
//
// You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (r *Router) Group(path string, middlewares ...MiddlewareHandlerFunc) *Router {
	return &Router{
		RedirectTrailingSlash:  r.RedirectTrailingSlash,
		RedirectFixedPath:      r.RedirectFixedPath,
		HandleMethodNotAllowed: r.HandleMethodNotAllowed,
		HandleOptions:          r.HandleOptions,
		parent:                 r,
		root:                   false,
		basePath:               joinPaths(r.basePath, path),
		mws:                    r.mws.Clone(middlewares...),
		trees:                  r.trees,
		Body404:                r.Body404,
		Body405:                r.Body405,
	}
}

// Attach another router to current one
func (r *Router) Attach(prefix string, router *Router) {
	for _, route := range router.Routes() {
		path := joinPaths(prefix, route.Path)
		r.Handle(route.Method, path, route.Handler, route.Mws.stack...)
	}
}

// AttachRoutes to current routes
func (r *Router) AttachRoutes(prefix string, routes Routes) {
	for _, route := range routes {
		path := joinPaths(prefix, route.Path)
		mws := r.mws.Clone(route.Mws.stack...)
		r.Handle(route.Method, path, route.Handler, mws.stack...)
	}
}

// Routes returns a slice of registered routes
func (r *Router) Routes() (routes Routes) {
	for method, tree := range r.trees {
		routes = iterate("", method, routes, tree)
	}
	return routes
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handler HandlerFunc, middlewares ...MiddlewareHandlerFunc) {

	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if handler == nil {
		panic("handler must not be nil")
	}

	path = joinPaths(r.basePath, path)

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	route := &Route{
		router:  r,
		Method:  method,
		Path:    path,
		Mws:     r.mws.Clone(middlewares...),
		Handler: handler,
	}

	root.addRoute(path, route)

	r.updateParams(path)

}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (*Route, Params, bool) {
	if root := r.trees[method]; root != nil {
		handle, ps, tsr := root.getValue(path, r.getParams)
		if handle == nil {
			r.putParams(ps)
			return nil, nil, tsr
		}
		if ps == nil {
			return handle, nil, tsr
		}
		return handle, *ps, tsr
	}
	return nil, nil, false
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	res := r.dispatchRequest(w, req)
	if res == nil {
		return
	}

	if err := res.Handle(w, req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Response error: %v", err)))
		return
	}
}

func (r *Router) dispatchRequest(w http.ResponseWriter, req *http.Request) Response {
	path := req.URL.Path
	if root := r.trees[req.Method]; root != nil {
		if route, ps, tsr := root.getValue(path, r.getParams); route != nil {

			if ps != nil {
				ctx := req.Context()
				ctx = context.WithValue(ctx, ParamsKey, *ps)
				req = req.WithContext(ctx)
				r.putParams(ps)
			}
			return route.HandleRequest(w, req)
		} else if req.Method != http.MethodConnect && path != "/" {
			code := http.StatusMovedPermanently
			if req.Method != http.MethodGet {
				code = http.StatusPermanentRedirect
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				return ResponseRedirect(code, req.URL.String())
			}
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = fixedPath
					return ResponseRedirect(code, req.URL.String())
				}
			}
		}
	}

	if req.Method == http.MethodOptions && r.HandleOptions {
		if allow := r.allowed(path, http.MethodOptions); allow != "" {

			allowed := strings.Split(allow, ",")

			route, _, _ := r.trees[allowed[0]].getValue(path, nil)

			mws := r.findRoot().mws
			if route != nil {
				mws = route.router.findRoot().mws
			}

			// Add request method to list of allowed methods
			allowed = append(allowed, http.MethodOptions)

			headers := map[string]string{"Allow": strings.Join(allowed, ", ")}

			r := &Route{
				Method: req.Method,
				Path:   path,
				Mws:    mws,
				Handler: func(r *http.Request) Response {
					return ResponseHeader(http.StatusNoContent, headers)
				},
			}

			return r.HandleRequest(w, req)
		}

	} else if r.HandleMethodNotAllowed {
		if allow := r.allowed(path, req.Method); allow != "" {
			w.Header().Set("Allow", allow)
			return ResponseError(http.StatusMethodNotAllowed, errors.New(r.Body405))
		}
	}

	return ResponseError(http.StatusNotFound, errors.New(r.Body404))
}

func (r *Router) findRoot() *Router {
	if r.root {
		return r
	}
	if r.parent != nil {
		return r.parent
	}

	return r
}

func (r *Router) updateParams(path string) {
	if r.parent != nil {
		r.parent.updateParams(path)
		return
	}
	// Update maxParams
	if paramsCount := countParams(path); paramsCount > r.maxParams {
		r.maxParams = paramsCount
	}

	// Lazy-init paramsPool alloc func
	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}

}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		for method := range r.trees {
			if method == http.MethodOptions {
				continue
			}
			// Add request method to list of allowed methods
			allowed = append(allowed, method)
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == http.MethodOptions {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path, nil)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}
	return
}
