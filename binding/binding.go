package binding

import "net/http"

// Content-Type MIME for common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEYAML              = "application/x-yaml"
)

// Binder defines interface which needs to be implemented for request data binding
type Binder interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

// BodyBinder adds BindBody method to Binder interface.
type BodyBinder interface {
	Binder
	BindBody([]byte, interface{}) error
}

// URIBinder adds BindUri method to Binder interface.
type URIBinder interface {
	Name() string
	BindUri(map[string][]string, interface{}) error
}

// Default Binder implementations
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	YAML          = yamlBinding{}
	URI           = uriBinding{}
)

// Default returns the appropriate Binder instance based on the HTTP method
// and the content type.
func Default(method, contentType string) Binder {
	if method == "GET" {
		return Form
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEYAML:
		return YAML
	default: //case MIMEPOSTForm, MIMEMultipartPOSTForm:
		return Form
	}
}
