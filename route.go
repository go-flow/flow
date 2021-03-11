package flow

import (
	"net/http"
)

// Route structure
type Route struct {
	Method string
	Path   string

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

	var resp Response

	chain := rt.Mws.Clone((func(next MiddlewareFunc) MiddlewareFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			resp = rt.Handler(r)
			return nil
		}
	}))

	if err := chain.handle(w, r); err != nil {
		return ResponseError(http.StatusInternalServerError, err)
	}
	return resp
}
