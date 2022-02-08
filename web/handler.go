package web

import "net/http"

// HandlerFunc is a function that is registered to a route to handle http requests
type HandlerFunc func(*http.Request, Params) Response
