package response

import (
	"bytes"
	"io"
	"net/http"

	"github.com/go-flow/flow/render"
)

type Render struct {
	render.Renderer
	code int
}

func (re *Render) Status() int {
	return re.code
}

func (rr *Render) Handle(w http.ResponseWriter, r *http.Request) error {
	var res bytes.Buffer
	if err := rr.Render(&res); err != nil {
		return err
	}

	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = rr.ContentType()
	}

	w.WriteHeader(rr.code)

	_, err := w.Write(res.Bytes())
	return err

}

// ResponseJSON creates JSON Response
func NewJSON(code int, data interface{}) *Render {
	return &Render{
		Renderer: render.JSON{Data: data},
		code:     code,
	}
}

// NewText creates Text Response
func NewText(code int, text string) *Render {
	return &Render{
		Renderer: render.Text{Data: text},
		code:     code,
	}
}

// NewXML creates XML Response
func NewXML(code int, data interface{}) *Render {
	return &Render{
		Renderer: render.XML{Data: data},
		code:     code,
	}
}

//NeweData creates []byte Response
func NewData(code int, data []byte, contentType []string) *Render {
	return &Render{
		Renderer: render.Data{
			Data:  data,
			CType: contentType,
		},
		code: code,
	}
}

// NewReader creates io.Reader render Response for given http code, reader and content type
func NewReader(code int, reader io.Reader, contentType []string) *Render {
	return &Render{
		Renderer: render.Reader{
			Reader: reader,
			CType:  contentType,
		},
		code: code,
	}
}
