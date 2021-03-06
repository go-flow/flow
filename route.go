package flow

import (
	"net/http"
)

// Route structure
type Route struct {
	router  *Router
	Method  string
	Path    string
	Mws     *MiddlewareStack
	Handler HandlerFunc
}

// Routes is Route collection
type Routes []Route

// HandleRequest handles http request. It executes all route middlewares and action handler
func (rt *Route) HandleRequest(w http.ResponseWriter, r *http.Request) Response {
	if rt.Mws == nil {
		return rt.Handler(r)
	}
	// define last handler in chain
	h := func(_ http.ResponseWriter, r *http.Request) Response {
		return rt.Handler(r)
	}

	// loop through middlewares and chain calls
	for i := len(rt.Mws.stack) - 1; i >= 0; i-- {
		h = rt.Mws.stack[i](h)
	}

	return h(w, r)
}
