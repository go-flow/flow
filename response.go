package flow

import (
	"net/http"

	"github.com/go-flow/flow/response"
)

// Response defines interface for HTTP action responses
type Response interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

// ResponseRedirect creates http Response redirect for given code and url
func ResponseRedirect(code int, url string) Response {
	return response.NewRedirect(code, url)
}
