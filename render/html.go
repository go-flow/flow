package render

import (
	"html/template"
	"net/http"

	"github.com/go-flow/flow/render/view"
)

var htmlContentType = []string{"text/html; charset=utf-8"}

// HTML renders html content
type HTML struct {
	Engine  view.Engine
	Name    string
	Data    map[string]interface{}
	Helpers template.FuncMap
}

// Render writes HTML content to response writer
func (r HTML) Render(w http.ResponseWriter) error {
	writeContentType(w, htmlContentType)
	return r.Engine.Render(w, r.Name, r.Data, r.Helpers)
}
