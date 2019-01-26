package view

import "html/template"

// Config represents view configuration object
type Config struct {
	Root         string           //view root
	Ext          string           //template extension
	Master       string           //template master
	Partials     []string         //template partial, such as head, foot
	Funcs        template.FuncMap //template functions
	DisableCache bool             //disable cache, debug mode
	Delims       Delims           //delimeters
}

// Delims struct holds template delimeters info
type Delims struct {
	Left  string
	Right string
}
