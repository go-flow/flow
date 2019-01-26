package render

import (
	"io"
)

// Text renders content to Text format
type Text struct {
	Data string
}

var textContentType = []string{"text/plain; charset=utf-8"}

// Render Plain Text io.Writer
func (r Text) Render(out io.Writer) error {
	_, err := io.WriteString(out, r.Data)
	return err
}
