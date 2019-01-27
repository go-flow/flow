package render

import (
	"io"
	"net/http"
)

// Renderer defines for rendering different content types
type Renderer interface {
	// Render writes data to io.Writer
	Render(io.Writer) error

	// ContentType returns contentType for renderer
	ContentType() []string
}

// writeHeaders is helper function for writing headers
// to provided ResponseWriter
func writeHeaders(w http.ResponseWriter, headers map[string]string) {
	header := w.Header()
	for k, v := range headers {
		if val := header[k]; len(val) == 0 {
			header[k] = []string{v}
		}
	}
}
