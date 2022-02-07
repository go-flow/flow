package render

import "io"

// Data renders byte array
type Data struct {
	Data  []byte
	CType []string
}

// Render writes []byte to io.Writer
func (r Data) Render(out io.Writer) (err error) {
	_, err = out.Write(r.Data)
	return
}

// ContentType returns contentType for renderer
func (r Data) ContentType() []string {
	return r.CType
}
