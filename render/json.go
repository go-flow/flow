package render

import (
	"net/http"

	"github.com/go-flow/flow/internal/json"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

// JSON renders data as JSON content type
type JSON struct {
	Data interface{}
}

// Render JSON content with application/json content type
func (r JSON) Render(w http.ResponseWriter) error {
	data, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	writeContentType(w, jsonContentType)
	w.Write(data)
	return nil
}
