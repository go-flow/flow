package render

import (
	"io"
	"net/http"
)

// Text renders content to Text format
type Text struct {
	Data string
}

var textContentType = []string{"text/plain; charset=utf-8"}

// Render content as Text Plain format to ResponseWriter
func (r Text) Render(w http.ResponseWriter) error {
	writeContentType(w, textContentType)
	_, err := io.WriteString(w, r.Data)
	return err
}
