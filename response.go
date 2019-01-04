package flow

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
//
// A ResponseWriter may not be used after the Handler.ServeHTTP method
// has returned.
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// Returns true if the response body was already written.
	Written() bool

	// Forces to write the http header (status code + headers).
	WriteHeaderNow()

	// get the http.Pusher for server push
	Pusher() http.Pusher
}

// A Response implements ResponseWriter interface and it
// is used by flow.Context to construct an HTTP response.
type Response struct {
	http.ResponseWriter
	size   int
	status int
}

func (w *Response) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (w *Response) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			fmt.Fprintf(os.Stderr, "[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

// WriteHeaderNow forces to write the http header (status code + headers).
func (w *Response) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *Response) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

// WriteString writes the string into the response body.
func (w *Response) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return
}

// Status returns the HTTP response status code of the current request.
func (w *Response) Status() int {
	return w.status
}

// Size returns the number of bytes already written into the response http body.
// See Written()
func (w *Response) Size() int {
	return w.size
}

// Written retunrs true if response is sent to http body
func (w *Response) Written() bool {
	return w.size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (w *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify interface.
func (w *Response) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flush interface.
func (w *Response) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

// Pusher returns http.Pusher object if underlying Response writer
// implements http.Pusher interface
func (w *Response) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
