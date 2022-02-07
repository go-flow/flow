package response

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

type Download struct {
	code   int
	name   string
	reader io.Reader
}

// NewDownload creates file attachment Response with following headers:
//
//   Content-Type
//   Content-Length
//   Content-Disposition
//
// Content-Type is set using mime#TypeByExtension with the filename's extension. Content-Type will default to
// application/octet-stream if using a filename with an unknown extension.
func NewDownload(code int, name string, reader io.Reader) *Download {
	return &Download{
		code:   code,
		name:   name,
		reader: reader,
	}
}

func (rd *Download) Status() int {
	return rd.code
}

func (rd *Download) Handle(w http.ResponseWriter, r *http.Request) error {
	ext := filepath.Ext(rd.name)
	t := mime.TypeByExtension(ext)
	if t == "" {
		t = "application/octet-stream"
	}

	cd := fmt.Sprintf("attachment; filename=%s", rd.name)
	h := w.Header()
	h.Add("Content-Disposition", cd)
	h.Add("Content-Type", t)

	_, err := io.Copy(w, rd.reader)
	return err
}
