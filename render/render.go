package render

import "net/http"

// Renderer defines for rendering different content types
type Renderer interface {
	// Render writes data with custom ContentType
	Render(http.ResponseWriter) error
}

// writeContentType is helper function for writing content-type
// to provided ResponseWriter

// if responseWriter has "Content-Type" already set provided values will be ignored
func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
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
