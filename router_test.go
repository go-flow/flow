// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

package flow

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func TestParams(t *testing.T) {
	ps := Params{
		Param{"param1", "value1"},
		Param{"param2", "value2"},
		Param{"param3", "value3"},
	}
	for i := range ps {
		if val := ps.ByName(ps[i].Key); val != ps[i].Value {
			t.Errorf("Wrong value for %s: Got %s; Want %s", ps[i].Key, val, ps[i].Value)
		}
	}
	if val := ps.ByName("noKey"); val != "" {
		t.Errorf("Expected empty string for not found key; got: %s", val)
	}
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	routed := false
	router.Handle(http.MethodGet, "/user/:name", func(r *http.Request) Response {
		routed = true
		want := Params{Param{"name", "gopher"}}
		ps := ParamsFromContext(r.Context())
		if !reflect.DeepEqual(ps, want) {
			t.Fatalf("wrong wildcard values: want %v, got %v", want, ps)
		}
		return ResponseText(http.StatusOK, "ok")
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest(http.MethodGet, "/user/gopher", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing failed")
	}
}

func TestRouterAPI(t *testing.T) {
	var get, head, options, post, put, patch, delete bool

	router := NewRouter()
	router.GET("/GET", func(r *http.Request) Response {
		get = true
		return ResponseText(http.StatusOK, "ok")
	})
	router.HEAD("/GET", func(r *http.Request) Response {
		head = true
		return ResponseText(http.StatusOK, "ok")
	})
	router.OPTIONS("/GET", func(r *http.Request) Response {
		options = true
		return ResponseText(http.StatusOK, "ok")
	})
	router.POST("/POST", func(r *http.Request) Response {
		post = true
		return ResponseText(http.StatusOK, "ok")
	})
	router.PUT("/PUT", func(r *http.Request) Response {
		put = true
		return ResponseText(http.StatusOK, "ok")
	})
	router.PATCH("/PATCH", func(r *http.Request) Response {
		patch = true
		return ResponseText(http.StatusOK, "ok")
	})
	router.DELETE("/DELETE", func(r *http.Request) Response {
		delete = true
		return ResponseText(http.StatusOK, "ok")
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest(http.MethodGet, "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest(http.MethodHead, "/GET", nil)
	router.ServeHTTP(w, r)
	if !head {
		t.Error("routing HEAD failed")
	}

	r, _ = http.NewRequest(http.MethodOptions, "/GET", nil)
	router.ServeHTTP(w, r)
	if !options {
		t.Error("routing OPTIONS failed")
	}

	r, _ = http.NewRequest(http.MethodPost, "/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest(http.MethodPut, "/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest(http.MethodPatch, "/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest(http.MethodDelete, "/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}
}

func TestRouterInvalidInput(t *testing.T) {
	router := NewRouter()

	handle := func(_ *http.Request) Response {
		return ResponseText(http.StatusOK, "OK")
	}

	recv := catchPanic(func() {
		router.Handle("", "/", handle)
	})
	if recv == nil {
		t.Fatal("registering empty method did not panic")
	}

	recv = catchPanic(func() {
		router.GET("", handle)
	})
	if recv == nil {
		t.Fatal("registering empty path did not panic")
	}

	recv = catchPanic(func() {
		router.GET("noSlashRoot", handle)
	})
	if recv == nil {
		t.Fatal("registering path not beginning with '/' did not panic")
	}

	recv = catchPanic(func() {
		router.GET("/", nil)
	})
	if recv == nil {
		t.Fatal("registering nil handler did not panic")
	}
}

func BenchmarkAllowed(b *testing.B) {
	handlerFunc := func(_ *http.Request) Response {
		return ResponseText(http.StatusOK, "OK")
	}

	router := NewRouter()
	router.POST("/path", handlerFunc)
	router.GET("/path", handlerFunc)

	b.Run("Global", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = router.allowed("*", http.MethodOptions)
		}
	})
	b.Run("Path", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = router.allowed("/path", http.MethodOptions)
		}
	})
}

func TestRouterOPTIONS(t *testing.T) {
	handlerFunc := func(_ *http.Request) Response {
		return ResponseText(http.StatusOK, "OK")
	}

	router := NewRouter()
	router.POST("/path", handlerFunc)

	// test not allowed
	// * (server)
	r, _ := http.NewRequest(http.MethodOptions, "*", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusNoContent) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "OPTIONS, POST" {
		t.Error("unexpected Allow header value: " + allow)
	}

	// path
	r, _ = http.NewRequest(http.MethodOptions, "/path", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusNoContent) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "OPTIONS, POST" {
		t.Error("unexpected Allow header value: " + allow)
	}

	r, _ = http.NewRequest(http.MethodOptions, "/doesnotexist", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusNotFound) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	// custom handler
	var custom bool
	router.OPTIONS("/path", func(r *http.Request) Response {
		custom = true
		return ResponseText(http.StatusOK, "OK")
	})

	// path
	r, _ = http.NewRequest(http.MethodOptions, "/path", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusOK) {
		t.Errorf("OPTIONS handling failed: Code=%d, Header=%v", w.Code, w.Header())
	}
	if !custom {
		t.Error("custom handler not called")
	}
}

func TestRouterNotAllowed(t *testing.T) {
	handlerFunc := func(_ *http.Request) Response {
		return ResponseText(http.StatusOK, "OK")
	}

	router := NewRouter()
	router.POST("/path", handlerFunc)

	// test not allowed
	r, _ := http.NewRequest(http.MethodGet, "/path", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusMethodNotAllowed) {
		t.Errorf("NotAllowed handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "OPTIONS, POST" {
		t.Error("unexpected Allow header value: " + allow)
	}

	// add another method
	router.DELETE("/path", handlerFunc)
	router.OPTIONS("/path", handlerFunc) // must be ignored

	// test again
	r, _ = http.NewRequest(http.MethodGet, "/path", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusMethodNotAllowed) {
		t.Errorf("NotAllowed handling failed: Code=%d, Header=%v", w.Code, w.Header())
	} else if allow := w.Header().Get("Allow"); allow != "DELETE, OPTIONS, POST" {
		t.Error("unexpected Allow header value: " + allow)
	}

}

func TestRouterNotFound(t *testing.T) {
	handlerFunc := func(_ *http.Request) Response {
		return ResponseText(http.StatusOK, "OK")
	}

	router := NewRouter()
	router.GET("/path", handlerFunc)
	router.GET("/dir/", handlerFunc)
	router.GET("/", handlerFunc)

	testRoutes := []struct {
		route    string
		code     int
		location string
	}{
		{"/path/", http.StatusMovedPermanently, "/path"},   // TSR -/
		{"/dir", http.StatusMovedPermanently, "/dir/"},     // TSR +/
		{"", http.StatusMovedPermanently, "/"},             // TSR +/
		{"/PATH", http.StatusMovedPermanently, "/path"},    // Fixed Case
		{"/DIR/", http.StatusMovedPermanently, "/dir/"},    // Fixed Case
		{"/PATH/", http.StatusMovedPermanently, "/path"},   // Fixed Case -/
		{"/DIR", http.StatusMovedPermanently, "/dir/"},     // Fixed Case +/
		{"/../path", http.StatusMovedPermanently, "/path"}, // CleanPath
		{"/nope", http.StatusNotFound, ""},                 // NotFound
	}
	for _, tr := range testRoutes {
		r, _ := http.NewRequest(http.MethodGet, tr.route, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		if !(w.Code == tr.code && (w.Code == http.StatusNotFound || fmt.Sprint(w.Header().Get("Location")) == tr.location)) {
			t.Errorf("NotFound handling route %s failed: Code=%d, Header=%v", tr.route, w.Code, w.Header().Get("Location"))
		}
	}

	// Test special case where no node for the prefix "/" exists
	router = NewRouter()
	router.GET("/a", handlerFunc)
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusNotFound) {
		t.Errorf("NotFound handling route / failed: Code=%d", w.Code)
	}
}

func TestRouterLookup(t *testing.T) {
	routed := false
	wantHandle := func(_ *http.Request) Response {
		routed = true
		return ResponseText(http.StatusOK, "OK")
	}

	wantParams := Params{Param{"name", "gopher"}}

	router := NewRouter()

	// try empty router first
	route, _, tsr := router.Lookup(http.MethodGet, "/nope")
	if route != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", route)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.GET("/user/:name", wantHandle)
	route, params, _ := router.Lookup(http.MethodGet, "/user/gopher")
	if route == nil {
		t.Fatal("Got no handle!")
	} else {
		route.HandleRequest(nil, nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}
	if !reflect.DeepEqual(params, wantParams) {
		t.Fatalf("Wrong parameter values: want %v, got %v", wantParams, params)
	}
	routed = false

	// route without param
	router.GET("/user", wantHandle)
	route, params, _ = router.Lookup(http.MethodGet, "/user")
	if route == nil {
		t.Fatal("Got no handle!")
	} else {
		route.HandleRequest(nil, nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}
	if params != nil {
		t.Fatalf("Wrong parameter values: want %v, got %v", nil, params)
	}

	route, _, tsr = router.Lookup(http.MethodGet, "/user/gopher/")
	if route != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", route)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	route, _, tsr = router.Lookup(http.MethodGet, "/nope")
	if route != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", route)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}

func TestRouterParamsFromContext(t *testing.T) {
	routed := false

	wantParams := Params{Param{"name", "gopher"}}
	handlerFunc := func(req *http.Request) Response {
		// get params from request context
		params := ParamsFromContext(req.Context())

		if !reflect.DeepEqual(params, wantParams) {
			t.Fatalf("Wrong parameter values: want %v, got %v", wantParams, params)
		}

		routed = true
		return ResponseText(http.StatusOK, "OK")
	}

	var nilParams Params
	handlerFuncNil := func(req *http.Request) Response {
		// get params from request context
		params := ParamsFromContext(req.Context())

		if !reflect.DeepEqual(params, nilParams) {
			t.Fatalf("Wrong parameter values: want %v, got %v", nilParams, params)
		}

		routed = true
		return ResponseText(http.StatusOK, "OK")
	}
	router := NewRouter()
	router.Handle(http.MethodGet, "/user", handlerFuncNil)
	router.Handle(http.MethodGet, "/user/:name", handlerFunc)

	w := new(mockResponseWriter)
	r, _ := http.NewRequest(http.MethodGet, "/user/gopher", nil)
	router.ServeHTTP(w, r)
	if !routed {
		t.Fatal("Routing failed!")
	}

	routed = false
	r, _ = http.NewRequest(http.MethodGet, "/user", nil)
	router.ServeHTTP(w, r)
	if !routed {
		t.Fatal("Routing failed!")
	}
}
