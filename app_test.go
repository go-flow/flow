package flow

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppServeHTTPDefault(t *testing.T) {

	app := New()
	app.GET("/ok", func(c *Context) {
		c.Status(http.StatusOK)
	})

	app.GET("/nok", func(c *Context) {
		c.AbortWithError(http.StatusInternalServerError, errors.New("nok error"))
	})

	req, err := http.NewRequest("GET", "/ok", nil)
	if err != nil {
		t.Errorf("An error occured. %v", err)
	}

	rr := httptest.NewRecorder()

	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

}

func TestAppServeHTTPError(t *testing.T) {

	app := New()

	app.GET("/nok", func(c *Context) {
		c.AbortWithError(http.StatusInternalServerError, errors.New("nok error"))
	})

	req, err := http.NewRequest("GET", "/nok", nil)
	if err != nil {
		t.Errorf("An error occured. %v", err)
	}

	rr := httptest.NewRecorder()

	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	}
}

func TestAppRoutes(t *testing.T) {

	// create test cases
	tt := []struct {
		Method string
		Path   string
	}{
		{Method: "GET", Path: "/get"},
		{Method: "HEAD", Path: "/head"},
		{Method: "OPTIONS", Path: "/options"},
		{Method: "POST", Path: "/post"},
		{Method: "PUT", Path: "/put"},
		{Method: "PATCH", Path: "/patch"},
		{Method: "DELETE", Path: "/delete"},
	}

	// create app
	app := New()

	//create handler
	handler := func(c *Context) { c.Status(http.StatusOK) }

	// register routes
	for _, r := range tt {
		switch r.Method {
		case "GET":
			app.GET(r.Path, handler)
		case "HEAD":
			app.HEAD(r.Path, handler)
		case "OPTIONS":
			app.OPTIONS(r.Path, handler)
		case "POST":
			app.POST(r.Path, handler)
		case "PUT":
			app.PUT(r.Path, handler)
		case "PATCH":
			app.PATCH(r.Path, handler)
		case "DELETE":
			app.DELETE(r.Path, handler)

		}
	}

	// run test cases
	for _, tc := range tt {
		t.Run(tc.Method, func(t *testing.T) {
			req, err := http.NewRequest(tc.Method, tc.Path, nil)
			if err != nil {
				t.Errorf("An error occured. %v", err)
			}

			rr := httptest.NewRecorder()

			app.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
			}
		})
	}
}
