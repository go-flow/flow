package render

import (
	"encoding/xml"
	"io"
)

// XML renders XML
type XML struct {
	Data interface{}
}

var xmlContentType = []string{"application/xml; charset=utf-8"}

// Render XML to io.Writer
func (r XML) Render(out io.Writer) error {
	return xml.NewEncoder(out).Encode(r.Data)
}
