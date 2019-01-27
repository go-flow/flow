package render

import (
	"html/template"
	"io"

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
func (r HTML) Render(out io.Writer) error {
	return r.Engine.Render(out, r.Name, r.Data, r.Helpers)
}

// ContentType returns contentType for renderer
func (HTML) ContentType() []string {
	return htmlContentType
}
