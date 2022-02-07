package response

import "net/http"

type Redirect struct {
	url  string
	code int
}

// NewRedirect creates new Response Rediredt object
func NewRedirect(code int, url string) *Redirect {
	return &Redirect{
		code: code,
		url:  url,
	}
}

func (re *Redirect) Status() int {
	return re.code
}

func (rr *Redirect) Handle(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, rr.url, rr.code)
	return nil
}
