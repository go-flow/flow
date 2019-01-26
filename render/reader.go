package render

import (
	"io"
)

// Reader renders content from io.Reader to response writer
type Reader struct {
	Reader io.Reader
}

// Render renders io.Writer content to responseWriter
func (r Reader) Render(out io.Writer) error {
	_, err := io.Copy(out, r.Reader)
	return err
}
