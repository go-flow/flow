package render

import (
	"encoding/xml"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo":  "bar",
		"num":  1,
		"html": "<p>html p</p>",
	}

	out := JSON{data}
	err := out.Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"\\u003cp\\u003ehtml p\\u003c/p\\u003e\",\"num\":1}", w.Body.String())
}

func TestJSONFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := make(chan int)

	assert.Error(t, (JSON{data}).Render(w))
}

type xmlmap map[string]interface{}

// Allows type H to be used with xml.Marshal
func (h xmlmap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func TestXML(t *testing.T) {
	w := httptest.NewRecorder()
	data := xmlmap{
		"foo": "bar",
	}

	err := XML{data}.Render(w)
	assert.NoError(t, err)
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "<map><foo>bar</foo></map>", w.Body.String())

}

func TestData(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte("#!Raw Data!!!")
	dr := Data{
		ContentType: "image/png",
		Data:        data,
	}

	err := dr.Render(w)
	assert.NoError(t, err)
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
	assert.Equal(t, "#!Raw Data!!!", w.Body.String())
}

func TestText(t *testing.T) {
	w := httptest.NewRecorder()

	txt1 := Text{"hello"}
	err := txt1.Render(w)
	assert.NoError(t, err)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "hello", w.Body.String())
}

func TestReader(t *testing.T) {
	w := httptest.NewRecorder()

	body := "#!PNG some raw data"
	headers := make(map[string]string)
	headers["Content-Disposition"] = `attachment; filename="filename.png"`

	err := (Reader{
		ContentLength: int64(len(body)),
		ContentType:   "image/png",
		Reader:        strings.NewReader(body),
		Headers:       headers,
	}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, body, w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
	assert.Equal(t, strconv.Itoa(len(body)), w.Header().Get("Content-Length"))
	assert.Equal(t, headers["Content-Disposition"], w.Header().Get("Content-Disposition"))
}
