package response

import "net/http"

type Headers struct {
	code    int
	headers map[string]string
}

func NewHeaders(code int, headers map[string]string) *Headers {
	return &Headers{
		code:    code,
		headers: headers,
	}
}

func (rh *Headers) Status() int {
	return rh.code
}

func (rh *Headers) Handle(w http.ResponseWriter, r *http.Request) error {
	for k, v := range rh.headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(rh.code)
	return nil
}
