package flow

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/go-flow/flow/render"
)

// Response defines interface for HTTP action responses
type Response interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

type responseError struct {
	err  error
	code int
}

func (re *responseError) Handle(w http.ResponseWriter, r *http.Request) error {
	var res Response
	ct := contentTypeFromString(r.Header.Get("Content-Type"))
	switch ct {
	case MIMEJSON:
		res = ResponseJSON(re.code, Map{"error": re.err.Error()})
	case MIMEXML, MIMEXML2:
		type Error struct {
			Message string
		}
		res = ResponseXML(re.code, &Error{Message: re.err.Error()})
	default:
		res = ResponseText(re.code, re.err.Error())
	}

	return res.Handle(w, r)

}

// ResponseError creates new Error response for given http code and error
func ResponseError(code int, err error) Response {
	return &responseError{
		code: code,
		err:  err,
	}
}

type responseRedirect struct {
	url  string
	code int
}

func (rr *responseRedirect) Handle(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, rr.url, rr.code)
	return nil
}

// ResponseRedirect creates Redirect ewsponse for given http code and destination URL
func ResponseRedirect(code int, url string) Response {
	return &responseRedirect{
		code: code,
		url:  url,
	}
}

type responseFile struct {
	filepath string
}

func (rf *responseFile) Handle(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, rf.filepath)
	return nil
}

// ResponseFile serves content from given file
func ResponseFile(filepath string) Response {
	return &responseFile{
		filepath: filepath,
	}
}

type responseDownload struct {
	name   string
	reader io.Reader
}

func (rd *responseDownload) Handle(w http.ResponseWriter, r *http.Request) error {
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

// ResponseDownload creates file attachment ActionResult with following headers:
//
//   Content-Type
//   Content-Length
//   Content-Disposition
//
// Content-Type is set using mime#TypeByExtension with the filename's extension. Content-Type will default to
// application/octet-stream if using a filename with an unknown extension.
func ResponseDownload(name string, reader io.Reader) Response {
	return &responseDownload{
		name:   name,
		reader: reader,
	}
}

type responseRender struct {
	render.Renderer
	code int
}

func (rr *responseRender) Handle(w http.ResponseWriter, r *http.Request) error {
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

// ResponseJSON creates JSON rendered Response
func ResponseJSON(code int, data interface{}) Response {
	return &responseRender{
		Renderer: render.JSON{Data: data},
		code:     code,
	}
}

// ResponseText creates Text rendered Response
func ResponseText(code int, text string) Response {
	return &responseRender{
		Renderer: render.Text{Data: text},
		code:     code,
	}
}

// ResponseXML creates XML rendered Response
func ResponseXML(code int, data interface{}) Response {
	return &responseRender{
		Renderer: render.XML{Data: data},
		code:     code,
	}
}

//ResponseData creates []byte render Response
func ResponseData(code int, data []byte, contentType []string) Response {
	return &responseRender{
		Renderer: render.Data{
			Data:  data,
			CType: contentType,
		},
		code: code,
	}
}

// ResponseReader creates io.Reader render Response for given http code, reader and content type
func ResponseReader(code int, reader io.Reader, contentType []string) Response {
	return &responseRender{
		Renderer: render.Reader{
			Reader: reader,
			CType:  contentType,
		},
		code: code,
	}
}
