package flow

import "net/http"

// MiddlewareFunc defines middleware handler function
type MiddlewareFunc func(w http.ResponseWriter, r *http.Request) error

// MiddlewareHandlerFunc defines middleware interface
//
// func DoSomething(next MiddlewareFunc) MiddlewareFunc {
// 	return func(w http.ResponseWriter, r *http.Request) error {
// 		// do something before calling the next handler
// 		err := next(w, r)
// 		// do something after call the handler
// 		return err
// 	}
// }
type MiddlewareHandlerFunc func(MiddlewareFunc) MiddlewareFunc

// MiddlewareStack holds middlewares applied to router
type MiddlewareStack struct {
	stack []MiddlewareHandlerFunc
}

// Append new Middlewares to stack
func (mws *MiddlewareStack) Append(mw ...MiddlewareHandlerFunc) {
	mws.stack = append(mws.stack, mw...)
}

// Clear current middleware stack
func (mws *MiddlewareStack) Clear() {
	mws.stack = []MiddlewareHandlerFunc{}
}

// Clone current stack to new one abd apply new middlewares
func (mws *MiddlewareStack) Clone(mw ...MiddlewareHandlerFunc) *MiddlewareStack {
	n := &MiddlewareStack{}
	n.Append(mws.stack...)
	n.Append(mw...)
	return n
}

func (mws *MiddlewareStack) handle(w http.ResponseWriter, r *http.Request, actionFn func(w http.ResponseWriter, r *http.Request) error) error {

	// define last handler in chain
	h := actionFn

	// loop through middlewares and chain calls
	for i := len(mws.stack) - 1; i >= 0; i-- {
		h = mws.stack[i](h)
	}

	return h(w, r)
}
