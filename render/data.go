package render

import "net/http"

// Data renders byte array
type Data struct {
	ContentType string
	Data        []byte
}

// Render writes []byte with custom ContentType
func (r Data) Render(w http.ResponseWriter) (err error) {
	writeContentType(w, []string{r.ContentType})
	_, err = w.Write(r.Data)
	return
}
