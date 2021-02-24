package flow

import "net/http"

// Content-Type MIME of the most common data formats.
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

// HandlerFunc is a function that is registered to a route to handle http requests
type HandlerFunc func(r *http.Request) Response
