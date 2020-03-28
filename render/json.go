package render

import (
	"encoding/json"
	"io"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

// JSON renders data as JSON content type
type JSON struct {
	Data interface{}
}

// Render JSON content to io.Writer
func (r JSON) Render(out io.Writer) error {
	data, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	out.Write(data)
	return nil
}

// ContentType returns contentType for renderer
func (JSON) ContentType() []string {
	return jsonContentType
}
