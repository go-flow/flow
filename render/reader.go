package render

import (
	"io"
)

// Reader renders content from io.Reader to response writer
type Reader struct {
	Reader io.Reader
	CType  []string
}

// Render renders io.Writer content to responseWriter
func (r Reader) Render(out io.Writer) error {
	_, err := io.Copy(out, r.Reader)
	return err
}

// ContentType returns contentType for renderer
func (r Reader) ContentType() []string {
	return r.CType
}
