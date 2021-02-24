package response

import "net/http"

type Redirect struct {
	url  string
	code int
}

func (rr *Redirect) Handle(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, rr.url, rr.code)
	return nil
}

func NewRedirect(code int, url string) *Redirect {
	return &Redirect{
		code: code,
		url:  url,
	}
}
