package flow

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-flow/flow/binding"
)

func newTestAppInstance() *App {
	opts := NewOptions()

	opts.UseViewEngine = false
	opts.UseRequestLogger = false
	opts.UseSession = false
	opts.UseTranslator = false

	return NewWithOptions(opts)
}

func createTestContext(w http.ResponseWriter) (ctx *Context, app *App) {
	app = newTestAppInstance()
	ctx = app.allocateContext()
	ctx.reset()
	ctx.writermem.reset(w)
	return
}

func compareFunc(t *testing.T, a, b interface{}) {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	if sf1.Pointer() != sf2.Pointer() {
		t.Error("different functions")
	}
}

func createMultipartRequest() *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	must(mw.SetBoundary(boundary))
	must(mw.WriteField("foo", "bar"))
	must(mw.WriteField("bar", "10"))
	must(mw.WriteField("bar", "foo2"))
	must(mw.WriteField("array", "first"))
	must(mw.WriteField("array", "second"))
	must(mw.WriteField("id", ""))
	must(mw.WriteField("time_local", "31/12/2016 14:55"))
	must(mw.WriteField("time_utc", "31/12/2016 14:55"))
	must(mw.WriteField("time_location", "31/12/2016 14:55"))
	must(mw.WriteField("names[a]", "thinkerou"))
	must(mw.WriteField("names[b]", "tianou"))
	req, err := http.NewRequest("POST", "/", body)
	must(err)
	req.Header.Set("Content-Type", binding.MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
