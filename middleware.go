package flow

import "net/http"

// MiddlewareFunc defines middleware handler function
type MiddlewareFunc func(w http.ResponseWriter, r *http.Request) Response

// MiddlewareHandlerFunc defines middleware interface
//
// func DoSomething(next MiddlewareFunc) MiddlewareFunc {
// 	return func(w http.ResponseWriter, r *http.Request) Response {
// 		// do something before calling the next handler
// 		resp := next(w, r)
// 		// do something after call the handler
// 		return resp
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
