package view

import (
	"html/template"
	"io"
)

// Engine defines View engine functionality
type Engine interface {
	Render(out io.Writer, name string, data map[string]interface{}, viewFuncs template.FuncMap) error
	SetViewHelpers(viewFuncs template.FuncMap)
}
