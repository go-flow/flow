package flow

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-flow/flow/binding"
	"github.com/stretchr/testify/assert"
)

func TestContextFormFile(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		_, err = w.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()
	c, _ := createTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.FormFile("file")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", f.Filename)
	}

	assert.NoError(t, c.SaveUploadedFile(f, "test"))
}

func TestContextFormFileFailed(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.Close()
	c, _ := createTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	c.app.MaxMultipartMemory = 8 << 20
	f, err := c.FormFile("file")
	assert.Error(t, err)
	assert.Nil(t, f)
}

func TestContextMultipartForm(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	assert.NoError(t, mw.WriteField("foo", "bar"))
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		_, err = w.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()

	c, _ := createTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.MultipartForm()
	if assert.NoError(t, err) {
		assert.NotNil(t, f)
	}

	assert.NoError(t, c.SaveUploadedFile(f.File["file"][0], "test"))
}

func TestSaveUploadedOpenFailed(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.Close()

	c, _ := createTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	f := &multipart.FileHeader{
		Filename: "file",
	}
	assert.Error(t, c.SaveUploadedFile(f, "test"))
}

func TestSaveUploadedCreateFailed(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		_, err = w.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()
	c, _ := createTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.FormFile("file")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", f.Filename)
	}

	assert.Error(t, c.SaveUploadedFile(f, "/"))
}

func TestContextReset(t *testing.T) {
	app := newTestAppInstance()
	c := app.allocateContext()

	assert.Equal(t, c.app, app)

	c.index = 2
	c.Response = &Response{ResponseWriter: httptest.NewRecorder()}
	c.Params = Params{Param{}}
	c.Error(errors.New("test")) // nolint: errcheck
	c.Set("foo", "bar")
	c.reset()

	assert.False(t, c.IsAborted())
	assert.Nil(t, c.Keys)
	assert.Nil(t, c.Accepted)
	assert.Len(t, c.Errors, 0)
	assert.Empty(t, c.Errors.Errors())
	assert.Len(t, c.Params, 0)
	assert.EqualValues(t, c.index, -1)
	assert.Equal(t, c.Response.(*Response), &c.writermem)
}

func TestContextHandlers(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	assert.Nil(t, c.handlers)
	assert.Nil(t, c.handlers.Last())

	c.handlers = HandlersChain{}
	assert.NotNil(t, c.handlers)
	assert.Nil(t, c.handlers.Last())

	f := func(c *Context) {}
	g := func(c *Context) {}

	c.handlers = HandlersChain{f}
	compareFunc(t, f, c.handlers.Last())

	c.handlers = HandlersChain{f, g}
	compareFunc(t, g, c.handlers.Last())
}

func TestContextSetGet(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("foo", "bar")

	value, err := c.Get("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, err)

	value, err = c.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, err)

	assert.Equal(t, "bar", c.MustGet("foo"))
	assert.Panics(t, func() { c.MustGet("no_exist") })
}

func TestContextSetGetValues(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("string", "this is a string")
	c.Set("int32", int32(-42))
	c.Set("int64", int64(42424242424242))
	c.Set("uint64", uint64(42))
	c.Set("float32", float32(4.2))
	c.Set("float64", 4.2)
	var a interface{} = 1
	c.Set("intInterface", a)

	assert.Exactly(t, c.MustGet("string").(string), "this is a string")
	assert.Exactly(t, c.MustGet("int32").(int32), int32(-42))
	assert.Exactly(t, c.MustGet("int64").(int64), int64(42424242424242))
	assert.Exactly(t, c.MustGet("uint64").(uint64), uint64(42))
	assert.Exactly(t, c.MustGet("float32").(float32), float32(4.2))
	assert.Exactly(t, c.MustGet("float64").(float64), 4.2)
	assert.Exactly(t, c.MustGet("intInterface").(int), 1)

}

func TestContextGetString(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("string", "this is a string")
	assert.Equal(t, "this is a string", c.GetString("string"))
}

func TestContextSetGetBool(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("bool", true)
	assert.True(t, c.GetBool("bool"))
}

func TestContextGetInt(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("int", 1)
	assert.Equal(t, 1, c.GetInt("int"))
}

func TestContextGetInt64(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("int64", int64(42424242424242))
	assert.Equal(t, int64(42424242424242), c.GetInt64("int64"))
}

func TestContextGetFloat64(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("float64", 4.2)
	assert.Equal(t, 4.2, c.GetFloat64("float64"))
}

func TestContextGetTime(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	t1, _ := time.Parse("1/2/2006 15:04:05", "01/01/2017 12:00:00")
	c.Set("time", t1)
	assert.Equal(t, t1, c.GetTime("time"))
}

func TestContextGetDuration(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("duration", time.Second)
	assert.Equal(t, time.Second, c.GetDuration("duration"))
}

func TestContextGetStringSlice(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Set("slice", []string{"foo"})
	assert.Equal(t, []string{"foo"}, c.GetStringSlice("slice"))
}

func TestContextGetStringMap(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	var m = make(map[string]interface{})
	m["foo"] = 1
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMap("map"))
	assert.Equal(t, 1, c.GetStringMap("map")["foo"])
}

func TestContextGetStringMapString(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	var m = make(map[string]string)
	m["foo"] = "bar"
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMapString("map"))
	assert.Equal(t, "bar", c.GetStringMapString("map")["foo"])
}

func TestContextGetStringMapStringSlice(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	var m = make(map[string][]string)
	m["foo"] = []string{"foo"}
	c.Set("map", m)

	assert.Equal(t, m, c.GetStringMapStringSlice("map"))
	assert.Equal(t, []string{"foo"}, c.GetStringMapStringSlice("map")["foo"])
}

func TestContextCopy(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.index = 2
	c.Request, _ = http.NewRequest("POST", "/hola", nil)
	c.handlers = HandlersChain{func(c *Context) {}}
	c.Params = Params{Param{Key: "foo", Value: "bar"}}
	c.Set("foo", "bar")

	cp := c.Copy()
	assert.Nil(t, cp.handlers)
	assert.Nil(t, cp.writermem.ResponseWriter)
	assert.Equal(t, &cp.writermem, cp.Response.(*Response))
	assert.Equal(t, cp.Request, c.Request)
	assert.Equal(t, cp.index, abortIndex)
	assert.Equal(t, cp.Keys, c.Keys)
	assert.Equal(t, cp.app, c.app)
	assert.Equal(t, cp.Params, c.Params)
}

var handlerTest HandlerFunc = func(c *Context) {

}

func TestContextHandler(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.handlers = HandlersChain{func(c *Context) {}, handlerTest}

	assert.Equal(t, reflect.ValueOf(handlerTest).Pointer(), reflect.ValueOf(c.Handler()).Pointer())
}

func TestContextQuery(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "http://example.com/?foo=bar&page=10&id=", nil)

	value, ok := c.GetQuery("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", value)
	assert.Equal(t, "bar", c.DefaultQuery("foo", "none"))
	assert.Equal(t, "bar", c.Query("foo"))

	value, ok = c.GetQuery("page")
	assert.True(t, ok)
	assert.Equal(t, "10", value)
	assert.Equal(t, "10", c.DefaultQuery("page", "0"))
	assert.Equal(t, "10", c.Query("page"))

	value, ok = c.GetQuery("id")
	assert.True(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.DefaultQuery("id", "nada"))
	assert.Empty(t, c.Query("id"))

	value, ok = c.GetQuery("NoKey")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Equal(t, "nada", c.DefaultQuery("NoKey", "nada"))
	assert.Empty(t, c.Query("NoKey"))

	// postform should not mess
	value, ok = c.GetPostForm("page")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.PostForm("foo"))
}

func TestContextQueryAndPostForm(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	body := bytes.NewBufferString("foo=bar&page=11&both=&foo=second")
	c.Request, _ = http.NewRequest("POST",
		"/?both=GET&id=main&array[]=first&array[]=second&ids[a]=hi&ids[b]=3.14", body)

	c.Request.Header.Add("Content-Type", binding.MIMEPOSTForm)

	assert.Equal(t, "application/x-www-form-urlencoded", c.Header("Content-Type"))

	assert.Equal(t, "bar", c.DefaultPostForm("foo", "none"))
	assert.Equal(t, "bar", c.PostForm("foo"))
	assert.Empty(t, c.Query("foo"))

	value, ok := c.GetPostForm("page")
	assert.True(t, ok)
	assert.Equal(t, "11", value)
	assert.Equal(t, "11", c.DefaultPostForm("page", "0"))
	assert.Equal(t, "11", c.PostForm("page"))
	assert.Empty(t, c.Query("page"))

	value, ok = c.GetPostForm("both")
	assert.True(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.PostForm("both"))
	assert.Empty(t, c.DefaultPostForm("both", "nothing"))
	assert.Equal(t, "GET", c.Query("both"), "GET")

	value, ok = c.GetQuery("id")
	assert.True(t, ok)
	assert.Equal(t, "main", value)
	assert.Equal(t, "000", c.DefaultPostForm("id", "000"))
	assert.Equal(t, "main", c.Query("id"))
	assert.Empty(t, c.PostForm("id"))

	value, ok = c.GetQuery("NoKey")
	assert.False(t, ok)
	assert.Empty(t, value)
	value, ok = c.GetPostForm("NoKey")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Equal(t, "nada", c.DefaultPostForm("NoKey", "nada"))
	assert.Equal(t, "nothing", c.DefaultQuery("NoKey", "nothing"))
	assert.Empty(t, c.PostForm("NoKey"))
	assert.Empty(t, c.Query("NoKey"))

	var obj struct {
		Foo   string   `form:"foo"`
		ID    string   `form:"id"`
		Page  int      `form:"page"`
		Both  string   `form:"both"`
		Array []string `form:"array[]"`
	}

	assert.NoError(t, c.Bind(&obj))
	assert.Equal(t, "bar", obj.Foo, "bar")
	assert.Equal(t, "main", obj.ID, "main")
	assert.Equal(t, 11, obj.Page, 11)
	assert.Empty(t, obj.Both)
	assert.Equal(t, []string{"first", "second"}, obj.Array)

	values, ok := c.GetQueryArray("array[]")
	assert.True(t, ok)
	assert.Equal(t, "first", values[0])
	assert.Equal(t, "second", values[1])

	values = c.QueryArray("array[]")
	assert.Equal(t, "first", values[0])
	assert.Equal(t, "second", values[1])

	values = c.QueryArray("nokey")
	assert.Equal(t, 0, len(values))

	values = c.QueryArray("both")
	assert.Equal(t, 1, len(values))
	assert.Equal(t, "GET", values[0])

	dicts, ok := c.GetQueryMap("ids")
	assert.True(t, ok)
	assert.Equal(t, "hi", dicts["a"])
	assert.Equal(t, "3.14", dicts["b"])

	dicts, ok = c.GetQueryMap("nokey")
	assert.False(t, ok)
	assert.Equal(t, 0, len(dicts))

	dicts, ok = c.GetQueryMap("both")
	assert.False(t, ok)
	assert.Equal(t, 0, len(dicts))

	dicts, ok = c.GetQueryMap("array")
	assert.False(t, ok)
	assert.Equal(t, 0, len(dicts))

	dicts = c.QueryMap("ids")
	assert.Equal(t, "hi", dicts["a"])
	assert.Equal(t, "3.14", dicts["b"])

	dicts = c.QueryMap("nokey")
	assert.Equal(t, 0, len(dicts))
}

func TestContextPostFormMultipart(t *testing.T) {
	c, _ := createTestContext(httptest.NewRecorder())
	c.Request = createMultipartRequest()

	var obj struct {
		Foo          string    `form:"foo"`
		Bar          string    `form:"bar"`
		BarAsInt     int       `form:"bar"`
		Array        []string  `form:"array"`
		ID           string    `form:"id"`
		TimeLocal    time.Time `form:"time_local" time_format:"02/01/2006 15:04"`
		TimeUTC      time.Time `form:"time_utc" time_format:"02/01/2006 15:04" time_utc:"1"`
		TimeLocation time.Time `form:"time_location" time_format:"02/01/2006 15:04" time_location:"Asia/Tokyo"`
		BlankTime    time.Time `form:"blank_time" time_format:"02/01/2006 15:04"`
	}
	assert.NoError(t, c.Bind(&obj))
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "10", obj.Bar)
	assert.Equal(t, 10, obj.BarAsInt)
	assert.Equal(t, []string{"first", "second"}, obj.Array)
	assert.Empty(t, obj.ID)
	assert.Equal(t, "31/12/2016 14:55", obj.TimeLocal.Format("02/01/2006 15:04"))
	assert.Equal(t, time.Local, obj.TimeLocal.Location())
	assert.Equal(t, "31/12/2016 14:55", obj.TimeUTC.Format("02/01/2006 15:04"))
	assert.Equal(t, time.UTC, obj.TimeUTC.Location())
	loc, _ := time.LoadLocation("Asia/Tokyo")
	assert.Equal(t, "31/12/2016 14:55", obj.TimeLocation.Format("02/01/2006 15:04"))
	assert.Equal(t, loc, obj.TimeLocation.Location())
	assert.True(t, obj.BlankTime.IsZero())

	value, ok := c.GetQuery("foo")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.Query("bar"))
	assert.Equal(t, "nothing", c.DefaultQuery("id", "nothing"))

	value, ok = c.GetPostForm("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", value)
	assert.Equal(t, "bar", c.PostForm("foo"))

	value, ok = c.GetPostForm("array")
	assert.True(t, ok)
	assert.Equal(t, "first", value)
	assert.Equal(t, "first", c.PostForm("array"))

	assert.Equal(t, "10", c.DefaultPostForm("bar", "nothing"))

	value, ok = c.GetPostForm("id")
	assert.True(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.PostForm("id"))
	assert.Empty(t, c.DefaultPostForm("id", "nothing"))

	value, ok = c.GetPostForm("nokey")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Equal(t, "nothing", c.DefaultPostForm("nokey", "nothing"))

	values, ok := c.GetPostFormArray("array")
	assert.True(t, ok)
	assert.Equal(t, "first", values[0])
	assert.Equal(t, "second", values[1])

	values = c.PostFormArray("array")
	assert.Equal(t, "first", values[0])
	assert.Equal(t, "second", values[1])

	values = c.PostFormArray("nokey")
	assert.Equal(t, 0, len(values))

	values = c.PostFormArray("foo")
	assert.Equal(t, 1, len(values))
	assert.Equal(t, "bar", values[0])

	dicts, ok := c.GetPostFormMap("names")
	assert.True(t, ok)
	assert.Equal(t, "thinkerou", dicts["a"])
	assert.Equal(t, "tianou", dicts["b"])

	dicts, ok = c.GetPostFormMap("nokey")
	assert.False(t, ok)
	assert.Equal(t, 0, len(dicts))

	dicts = c.PostFormMap("names")
	assert.Equal(t, "thinkerou", dicts["a"])
	assert.Equal(t, "tianou", dicts["b"])

	dicts = c.PostFormMap("nokey")
	assert.Equal(t, 0, len(dicts))
}
