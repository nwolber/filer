// Copyright (c) 2014 Niklas Wolber

package filer

import (
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestNew(t *testing.T) {
	s, _ := New("")
	if s == nil {
		t.Error("A static server should be creatable at all times.")
	}
}

var tests = []struct {
	dir  string
	url  string
	code int
}{
	{"tests", "../filer_test.go", http.StatusForbidden},
	{"tests", "invalid.js", http.StatusNotFound},
	{"tests", "test.js", http.StatusOK},
}

func TestServeHTTP(t *testing.T) {
	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.url, nil)
		s, _ := New(tt.dir)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != tt.code {
			t.Errorf("url: %s, want: %d, got: %d", tt.url, tt.code, rec.Code)
		}
	}
}
