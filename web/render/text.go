package render

import (
	"io"
)

var textContentType = []string{"text/plain; charset=utf-8"}

// Text renders content to Text format
type Text struct {
	Data string
}

// Render Plain Text io.Writer
func (r Text) Render(out io.Writer) error {
	_, err := io.WriteString(out, r.Data)
	return err
}

// ContentType returns contentType for renderer
func (Text) ContentType() []string {
	return textContentType
}
