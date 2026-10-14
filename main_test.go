package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	h := newMux(false)
	cases := []struct {
		path string
		code int
		want string
	}{
		{"/", 200, `href="/a/2"`},
		{"/a/2", 200, `href="/a/3"`},
		{"/a/3", 200, `Text me the link`},
		{"/b/4", 200, `href="/b/3"`},
		{"/b/5", 404, ""},
		{"/zz/2", 404, ""},
		{"/a/x", 404, ""},
		{"/a/1", 404, ""},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", c.path, nil))
		if rec.Code != c.code {
			t.Errorf("%s: got %d, want %d", c.path, rec.Code, c.code)
		}
		if !strings.Contains(rec.Body.String(), c.want) {
			t.Errorf("%s: missing %q", c.path, c.want)
		}
	}
	// Last page must not offer a Next link.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/a/3", nil))
	if strings.Contains(rec.Body.String(), `rel="next"`) {
		t.Error("/a/3: last page should have no next link")
	}
}
