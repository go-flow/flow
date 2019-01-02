package flow

// Router defines core routing methods
type Router interface {

	// Use appends one or more middlewares onto the Router stack.
	Use(middleware ...HandlerFunc)

	// GET is a shortcut for router.Handle("GET", path, handle)
	GET(path string, handle HandlerFunc)

	// HEAD is a shortcut for router.Handle("HEAD", path, handle)
	HEAD(path string, handle HandlerFunc)

	// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
	OPTIONS(path string, handle HandlerFunc)

	// POST is a shortcut for router.Handle("POST", path, handle)
	POST(path string, handle HandlerFunc)

	// PUT is a shortcut for router.Handle("PUT", path, handle)
	PUT(path string, handle HandlerFunc)

	// PATCH is a shortcut for router.Handle("PATCH", path, handle)
	PATCH(path string, handle HandlerFunc)

	// DELETE is a shortcut for router.Handle("DELETE", path, handle)
	DELETE(path string, handle HandlerFunc)

	// Handle registers a new request handle with the given path and method.
	//
	// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
	// functions can be used.
	//
	// This function is intended for bulk loading and to allow the usage of less
	// frequently used, non-standardized or custom methods (e.g. for internal
	// communication with a proxy).
	Handle(method, path string, handle HandlerFunc)

	// Lookup allows the manual lookup of a method + path combo.
	// This is e.g. useful to build a framework around this router.
	// If the path was found, it returns the handle function and the path parameter
	// values. Otherwise the third return value indicates whether a redirection to
	// the same path with an extra / without the trailing slash should be performed.
	Lookup(method, path string) (HandlersChain, Params, bool)

	// Routes returns a slice of registered routes, including some useful information, such as:
	// the http method, path and the handler name.
	Routes() Routes
}
