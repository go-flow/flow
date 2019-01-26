package render

import (
	"io"
	"net/http"
	"strconv"
)

// Reader renders content from io.Reader to response writer
type Reader struct {
	ContentType   string
	ContentLength int64
	Reader        io.Reader
	Headers       map[string]string
}

// Render renders io.Reader content to responseWriter
func (r Reader) Render(w http.ResponseWriter) error {
	writeContentType(w, []string{r.ContentType})
	r.Headers["Content-Length"] = strconv.FormatInt(r.ContentLength, 10)
	writeHeaders(w, r.Headers)
	_, err := io.Copy(w, r.Reader)
	return err
}
