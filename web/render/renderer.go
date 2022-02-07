package render

import "io"

// Renderer defines rendering interface for different content types
type Renderer interface {
	// Render writes data to io.Writer
	Render(io.Writer) error

	//ContentType returns renderer content type
	ContentType() []string
}
