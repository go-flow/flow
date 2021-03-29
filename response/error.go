package response

import "net/http"

type Error struct {
	err  error
	code int
}

func NewError(code int, err error) *Error {
	return &Error{
		code: code,
		err:  err,
	}
}

func (re *Error) Status() int {
	return re.code
}

func (re *Error) Handle(w http.ResponseWriter, r *http.Request) error {
	ct := contentTypeFromString(r.Header.Get("Content-Type"))
	switch ct {
	case MIMEJSON:
		res := NewJSON(re.code, map[string]interface{}{"error": re.err.Error()})
		return res.Handle(w, r)
	case MIMEXML, MIMEXML2:
		type Error struct {
			Message string
		}
		res := NewXML(re.code, &Error{Message: re.err.Error()})
		return res.Handle(w, r)
	default:
		res := NewText(re.code, re.err.Error())
		return res.Handle(w, r)
	}
}

func contentTypeFromString(ct string) string {
	for i, char := range ct {
		if char == ' ' || char == ';' {
			return ct[:i]
		}
	}
	return ct
}
