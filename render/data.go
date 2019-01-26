package render

import "io"

// Data renders byte array
type Data struct {
	Data []byte
}

// Render writes []byte to io.Writer
func (r Data) Render(out io.Writer) (err error) {
	_, err = out.Write(r.Data)
	return
}
