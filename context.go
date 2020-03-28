package flow

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-flow/flow/binding"
	"github.com/go-flow/flow/i18n"
	"github.com/go-flow/flow/log"
	"github.com/go-flow/flow/render"
)

// Context is request scoped application context
// It manages application request flow
type Context struct {
	app *App

	writermem Response

	Request  *http.Request
	Response ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int8

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors errorMsgs

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

	logger log.Logger
}

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

func (c *Context) reset() {
	c.Response = &c.writermem
	c.Params = c.Params[0:0]
	c.handlers = nil
	c.index = -1
	c.Keys = nil
	c.Errors = c.Errors[0:0]
	c.Accepted = nil
	c.logger = nil
}

// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This has to be used when the context has to be passed to a goroutine.
func (c *Context) Copy() *Context {
	var cp = *c
	cp.writermem.ResponseWriter = nil
	cp.Response = &cp.writermem
	cp.index = abortIndex
	cp.handlers = nil
	return &cp
}

// Handler returns the main handler.
func (c *Context) Handler() HandlerFunc {
	return c.handlers.Last()
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
//
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	for s := int8(len(c.handlers)); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Response.WriteHeaderNow()
	c.Abort()
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

/************************************/
/*********  APP MANAGEMENT  *********/
/************************************/

// Logger gets application Logger instance
func (c *Context) Logger() log.Logger {
	if c.logger != nil {
		return c.logger
	}
	return c.app.Logger
}

// LogFields adds adds structured context to context Logger
//
// This allows you to easily add things
// like metrics (think DB times) to your request.
func (c *Context) LogFields(fields log.Fields) {
	c.logger = c.Logger().WithFields(fields)
}

// AppOptions returns copy of application Options object
func (c *Context) AppOptions() Options {
	return c.app.Options
}

/************************************/
/********* ERROR MANAGEMENT *********/
/************************************/

// Error attaches an error to the current context. The error is pushed to a list of errors.
//
// Error will panic if err is nil.
func (c *Context) Error(err error) {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err: err,
		}
	}

	c.Errors = append(c.Errors, parsedError)
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

/************************************/
/************ INPUT DATA ************/
/************************************/

// Param returns the value of the URL param.
//
// It is a shortcut for c.Params.ByName(key)
func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}

// Query returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
//
// It is shortcut for `c.Request.URL.Query().Get(key)`
func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
//
// See: Query() and GetQuery() for further information.
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`.
//
// It is shortcut for `c.Request.URL.Query().Get(key)`
func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// QueryArray returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

// GetQueryArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetQueryArray(key string) ([]string, bool) {
	if values, ok := c.Request.URL.Query()[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// QueryMap returns a map for a given query key.
func (c *Context) QueryMap(key string) map[string]string {
	dicts, _ := c.GetQueryMap(key)
	return dicts
}

// GetQueryMap returns a map for a given query key, plus a boolean value
// whether at least one value exists for the given key.
func (c *Context) GetQueryMap(key string) (map[string]string, bool) {
	return c.get(c.Request.URL.Query(), key)
}

// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns an empty string `("")`.
func (c *Context) PostForm(key string) string {
	value, _ := c.GetPostForm(key)
	return value
}

// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the specified defaultValue string.
// See: PostForm() and GetPostForm() for further information.
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value, ok := c.GetPostForm(key); ok {
		return value
	}
	return defaultValue
}

// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns ("", false).
func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// PostFormArray returns a slice of strings for a given form key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) PostFormArray(key string) []string {
	values, _ := c.GetPostFormArray(key)
	return values
}

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	req := c.Request
	_ = req.ParseMultipartForm(c.app.MaxMultipartMemory)

	if values := req.PostForm[key]; len(values) > 0 {
		return values, true
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return values, true
		}
	}
	return []string{}, false
}

// PostFormMap returns a map for a given form key.
func (c *Context) PostFormMap(key string) map[string]string {
	dicts, _ := c.GetPostFormMap(key)
	return dicts
}

// GetPostFormMap returns a map for a given form key, plus a boolean value
// whether at least one value exists for the given key.
func (c *Context) GetPostFormMap(key string) (map[string]string, bool) {
	req := c.Request
	err := req.ParseMultipartForm(c.app.MaxMultipartMemory)
	if err != nil {
		panic(err)
	}

	dicts, exist := c.get(req.PostForm, key)

	if !exist && req.MultipartForm != nil && req.MultipartForm.File != nil {
		dicts, exist = c.get(req.MultipartForm.Value, key)
	}

	return dicts, exist
}

// get is an internal method and returns a map which satisfy conditions.
func (c *Context) get(m map[string][]string, key string) (map[string]string, bool) {
	dicts := make(map[string]string)
	exist := false
	for k, v := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				dicts[k[i+1:][:j]] = v[0]
			}
		}
	}
	return dicts, exist
}

// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(c.app.MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	_, fh, err := c.Request.FormFile(name)
	return fh, err
}

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(c.app.MaxMultipartMemory)
	return c.Request.MultipartForm, err
}

// SaveUploadedFile uploads the form file to specific dest.
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dest string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(out, src)
	return nil
}

// Bind checks the Content-Type to select a binding engine automatically,
func (c *Context) Bind(obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.BindWith(obj, b)
}

// BindJSON binds the passed struct pointer using JSON binding engine.
func (c *Context) BindJSON(obj interface{}) error {
	return c.BindWith(obj, binding.JSON)
}

// BindXML binds the passed struct pointer using XML binding engine.
func (c *Context) BindXML(obj interface{}) error {
	return c.BindWith(obj, binding.XML)
}

// BindQuery binds the passed struct pointer using Query binding engine.
func (c *Context) BindQuery(obj interface{}) error {
	return c.BindWith(obj, binding.Query)
}

// BindURI binds the passed struct pointer using URI binding engine.
func (c *Context) BindURI(obj interface{}) error {
	m := make(map[string][]string)
	for _, v := range c.Params {
		m[v.Key] = []string{v.Value}
	}
	return binding.URI.BindURI(m, obj)
}

// BindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func (c *Context) BindWith(obj interface{}, b binding.Binder) error {
	return b.Bind(c.Request, obj)
}

// ClientIP implements a best effort algorithm to return the real client IP
//
// it parses X-Real-IP and X-Forwarded-For in order to work properly
// with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
func (c *Context) ClientIP() string {

	clientIP := c.requestHeader("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(c.requestHeader("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if addr := c.requestHeader("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// RequestID implements a best effort algorithm to return tracing request ID for current request
//
// it parses X-Request-ID which is ment to be an application level tracing id
// and X-Amzn-Trace-Id which is automatically added by Amazon loadbalancers
func (c *Context) RequestID() string {
	// check if request ID exists in headers
	requestID := c.requestHeader("X-Request-ID")

	if requestID == "" {
		//check if  X-Amzn-Trace-Id exists
		requestID = c.requestHeader("X-Amzn-Trace-Id")
	}
	return requestID
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	return filterFlags(c.requestHeader("Content-Type"))
}

// SetContentType sets Content-Type header to response
func (c *Context) SetContentType(value []string) {
	header := c.Response.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// IsWebsocket returns true if the request headers indicate that a websocket
// handshake is being initiated by the client.
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.requestHeader("Connection")), "upgrade") &&
		strings.ToLower(c.requestHeader("Upgrade")) == "websocket" {
		return true
	}
	return false
}

func (c *Context) requestHeader(key string) string {
	return c.Request.Header.Get(key)
}

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.Response.WriteHeader(code)
}

// SetHeader is a intelligent shortcut for c.Response.Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `c.Response.Header().Del(key)`
func (c *Context) SetHeader(key, value string) {
	if value == "" {
		c.Response.Header().Del(key)
		return
	}
	c.Response.Header().Set(key, value)
}

// Header returns value from request headers.
func (c *Context) Header(key string) string {
	return c.requestHeader(key)
}

// Render writes the response headers and calls render.Render to render data.
func (c *Context) Render(code int, r render.Renderer) {
	var tpl bytes.Buffer
	if err := r.Render(&tpl); err != nil {
		c.ServeError(http.StatusInternalServerError, err)
		return
	}

	c.Status(code)
	_, err := c.Response.Write(tpl.Bytes())
	if err != nil {
		c.ServeError(http.StatusInternalServerError, err)
	}

}

// HTML renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (c *Context) HTML(code int, name string, obj interface{}) {
	if c.app.ViewEngine == nil {
		c.ServeError(http.StatusInternalServerError, errors.New("application view engine not enabled"))
		return
	}
	// request scoped view Helpers
	helpers := make(template.FuncMap)

	// request scoped data
	data := make(map[string]interface{})

	if c.app.SessionStore != nil {
		s := c.Session()
		// pass session
		data["session"] = s.Values()
		// save session
		s.Save()
	}

	// pass all context keys to view
	data["context"] = c.Keys
	// pass all context Errors
	data["errors"] = c.Errors.Errors()
	// object from action is passed to View as model
	data["model"] = obj

	// check if we use translations
	translator := c.app.Translator
	if translator != nil {
		// reload translations during development
		if c.AppOptions().Env == "development" {
			err := translator.Load()
			if err != nil {
				panic(err)
			}
		}

		// get languages
		langs := translator.ExtractLanguage(c)
		// define translation function
		transFunc, err := i18n.Tfunc(langs[0], langs[1:]...)
		if err != nil {
			c.Logger().Warn(err.Error())
			c.Logger().Warn("Your locale files are probably empty or missing")
		}

		// create viewHelper function
		helpers[translator.HelperName] = func(translationID string, args ...interface{}) string {
			return transFunc(translationID, args...)
		}
	}

	r := render.HTML{
		Engine:  c.app.ViewEngine,
		Name:    name,
		Data:    data,
		Helpers: helpers,
	}

	// set render contentType
	c.SetContentType(r.ContentType())

	// render
	c.Render(code, r)
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface{}) {
	r := render.JSON{Data: obj}
	c.SetContentType(r.ContentType())
	c.Render(code, r)
}

// XML serializes the given struct as XML into the response body.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(code int, obj interface{}) {
	r := render.XML{Data: obj}
	c.SetContentType(r.ContentType())
	c.Render(code, r)
}

// String writes the given string into the response body.
func (c *Context) String(code int, data string) {
	r := render.Text{Data: data}
	c.SetContentType(r.ContentType())
	c.Render(code, r)
}

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	if (code < 300 || code > 308) && code != 201 {
		panic(fmt.Errorf("can not redirect with status code %d", code))
	}
	http.Redirect(c.Response, c.Request, location, code)
}

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, data []byte) {
	c.Render(code, render.Data{
		Data: data,
	})
}

// GetRawData return stream data.
func (c *Context) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(c.Request.Body)
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Response, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

// File writes the specified file into the body stream in a efficient way.
func (c *Context) File(filepath string) {
	http.ServeFile(c.Response, c.Request, filepath)
}

// Stream sends a streaming response.
func (c *Context) Stream(step func(w io.Writer) bool) {
	w := c.Response
	clientGone := w.CloseNotify()
	for {
		select {
		case <-clientGone:
			return
		default:
			keepOpen := step(w)
			w.Flush()
			if !keepOpen {
				return
			}
		}
	}
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}

// ServeError serves error message with given code and message
// the error is served with text/plain mime type
func (c *Context) ServeError(code int, err error) {

	// store error in context error stack
	c.Error(err)

	// set response status
	c.Status(code)

	// execute all handlers in context
	c.Next()

	if c.Response.Written() {
		return
	}

	if c.app.methodNotAllowedHandler != nil && code == http.StatusMethodNotAllowed {
		c.app.methodNotAllowedHandler(c)
	} else if c.app.notFoundHandler != nil && code == http.StatusNotFound {
		c.app.notFoundHandler(c)
	} else if c.app.unauthorizedHandler != nil && code == http.StatusUnauthorized {
		c.app.unauthorizedHandler(c)
	} else if c.app.errorHandler != nil {
		c.app.errorHandler(c)
	} else {
		c.SetContentType([]string{"text/plain"})
		_, _ = c.Response.Write([]byte(err.Error()))
		return
	}
	c.Response.WriteHeaderNow()
}

/************************************/
/*****    SESSION MANAGEMENT    *****/
/************************************/

// Session gets session object for current request
func (c *Context) Session() *Session {
	if c.app.SessionStore == nil {
		c.Logger().Error("Session is not enabled in configuration")
		return nil
	}

	session, _ := c.app.SessionStore.Get(c.Request, c.app.SessionName)
	return &Session{
		Session: session,
		req:     c.Request,
		res:     c.Response,
	}
}

/************************************/
/*****  GOLANG.ORG/NET/CONTEXT  *****/
/************************************/

// Deadline returns the time when work done on behalf of this context
// should be canceled.
//
// Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *Context) Deadline() (time.Time, bool) {
	return c.Request.Context().Deadline()
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled.
//
// Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
func (c *Context) Done() <-chan struct{} {
	return c.Request.Context().Done()
}

// Err returns a non-nil error value after Done is closed,
// successive calls to Err return the same error.
//
// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
func (c *Context) Err() error {
	return c.Request.Context().Err()
}
