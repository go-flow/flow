package view

import (
	"html/template"
	"net/http"
)

// Render object
type Render struct {
	Engine  *Engine
	Name    string
	Data    interface{}
	Helpers template.FuncMap
}

// Render executes rendering on responseWriter
func (r Render) Render(w http.ResponseWriter) error {
	return r.Engine.executeRender(w, r.Name, r.Data, r.Helpers)
}

// WriteContentType writes Renderes content type to responseWriter
func (r Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = htmlContentType
	}
}
