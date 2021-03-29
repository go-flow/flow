package flow

import (
	"io"
	"net/http"

	"github.com/go-flow/flow/response"
)

// Response defines interface for HTTP action responses
type Response interface {
	Status() int
	Handle(w http.ResponseWriter, r *http.Request) error
}

// ResponseError creates new Error response for given http code and error
func ResponseError(code int, err error) Response {
	return response.NewError(code, err)
}

// ResponseRedirect creates Redirect ewsponse for given http code and destination URL
func ResponseRedirect(code int, url string) Response {
	return response.NewRedirect(code, url)
}

func ResponseHeader(code int, headers map[string]string) Response {
	return response.NewHeaders(code, headers)
}

// ResponseFile serves content from given file
func ResponseFile(filepath string) Response {
	return response.NewFile(filepath)
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
	return response.NewDownload(name, reader)
}

// ResponseJSON creates JSON rendered Response
func ResponseJSON(code int, data interface{}) Response {
	return response.NewJSON(code, data)
}

// ResponseText creates Text rendered Response
func ResponseText(code int, text string) Response {
	return response.NewText(code, text)
}

// ResponseXML creates XML rendered Response
func ResponseXML(code int, data interface{}) Response {
	return response.NewXML(code, data)
}

//ResponseData creates []byte render Response
func ResponseData(code int, data []byte, contentType []string) Response {
	return response.NewData(code, data, contentType)
}

// ResponseReader creates io.Reader render Response for given http code, reader and content type
func ResponseReader(code int, reader io.Reader, contentType []string) Response {
	return response.NewReader(code, reader, contentType)
}
