package render

import (
	"encoding/xml"
	"io"
)

var xmlContentType = []string{"application/xml; charset=utf-8"}

// XML renders XML
type XML struct {
	Data interface{}
}

// Render XML to io.Writer
func (r XML) Render(out io.Writer) error {
	return xml.NewEncoder(out).Encode(r.Data)
}

// ContentType returns contentType for renderer
func (XML) ContentType() []string {
	return xmlContentType
}
