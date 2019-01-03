package flow

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/go-flow/flow/render"
)

// App -
type App struct {
	config *Config
	router *Router
	pool   sync.Pool

	delims     render.Delims
	HTMLRender render.HTMLRender
	FuncMap    template.FuncMap
}

// ServeHTTP conforms to the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := a.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = r
	c.reset()

	a.handleHTTPRequest(c)
	a.pool.Put(c)
}

// HandleContext re-enter a context that has been rewritten.
// This can be done by setting c.Request.URL.Path to your new target.
// Disclaimer: You can loop yourself to death with this, use wisely.
func (a *App) HandleContext(c *Context) {
	c.reset()
	a.handleHTTPRequest(c)
}

// Delims sets template left and right delims and returns a Engine instance.
func (a *App) Delims(left, right string) {
	a.delims = render.Delims{Left: left, Right: right}
}

// LoadHTMLGlob loads HTML files identified by glob pattern
// and associates the result with HTML renderer.
func (a *App) LoadHTMLGlob(pattern string) {
	left := a.delims.Left
	right := a.delims.Right
	templ := template.Must(template.New("").Delims(left, right).Funcs(a.FuncMap).ParseGlob(pattern))

	a.SetHTMLTemplate(templ)
}

// LoadHTMLFiles loads a slice of HTML files
// and associates the result with HTML renderer.
func (a *App) LoadHTMLFiles(files ...string) {

	templ := template.Must(template.New("").Delims(a.delims.Left, a.delims.Right).Funcs(a.FuncMap).ParseFiles(files...))
	a.SetHTMLTemplate(templ)
}

// SetHTMLTemplate associate a template with HTML renderer.
func (a *App) SetHTMLTemplate(templ *template.Template) {
	a.HTMLRender = render.HTMLProduction{Template: templ.Funcs(a.FuncMap)}
}

// SetFuncMap sets the FuncMap used for template.FuncMap.
func (a *App) SetFuncMap(funcMap template.FuncMap) {
	a.FuncMap = funcMap
}

func (a *App) handleHTTPRequest(c *Context) {
	req := c.Request
	httpMethod := req.Method
	path := req.URL.Path

	if root := a.router.trees[httpMethod]; root != nil {
		if handlers, ps, tsr := root.getValue(path); handlers != nil {
			c.handlers = handlers
			c.Params = ps
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		} else if httpMethod != "CONNECT" && path != "/" {
			redirectTS := a.config.GetBool("redirectTrailingSlash")
			redirectFP := a.config.GetBool("redirectFixedPath")
			code := http.StatusMovedPermanently // Permanent redirect, request with GET method
			if httpMethod != "GET" {
				code = http.StatusTemporaryRedirect
			}
			if tsr && redirectTS {
				req.URL.Path = path + "/"
				if length := len(path); length > 1 && path[length-1] == '/' {
					req.URL.Path = path[:length-1]
				}
				// logger here
				http.Redirect(c.Writer, req, req.URL.String(), code)
				c.writermem.WriteHeaderNow()
				return
			}

			if redirectFP {
				fixedPath, found := root.findCaseInsensitivePath(CleanPath(path), redirectTS)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(c.Writer, req, req.URL.String(), code)
					c.writermem.WriteHeaderNow()
					return
				}
			}
		}
	}

	if a.config.GetBool("handleMethodNotAllowed") {
		if allow := a.router.allowed(path, httpMethod); len(allow) > 0 {
			c.handlers = a.router.Middlewares
			c.ServeError(http.StatusMethodNotAllowed, []byte(a.config.GetString("405Body")))
			return
		}
	}

	c.handlers = a.router.Middlewares
	c.ServeError(http.StatusNotFound, []byte(a.config.GetString("404Body")))
}

func (a *App) allocateContext() *Context {
	return &Context{app: a}
}
