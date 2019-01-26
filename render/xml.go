package render

import (
	"encoding/xml"
	"net/http"
)

// XML renders XML
type XML struct {
	Data interface{}
}

var xmlContentType = []string{"application/xml; charset=utf-8"}

// Render XML content with application/xml content type
func (r XML) Render(w http.ResponseWriter) error {
	writeContentType(w, xmlContentType)
	return xml.NewEncoder(w).Encode(r.Data)
}
