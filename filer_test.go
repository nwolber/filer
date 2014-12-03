// Copyright (c) 2014 Niklas Wolber
// This file is licensed under the MIT license.
// See the LICENSE file for more information.

package filer

import (
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestNew(t *testing.T) {
	s, _ := NewFileSystemFiler("")
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
	// should redirect to index.html
	{"tests", "/", http.StatusOK},
	{"tests", "", http.StatusOK},
}

func TestServeHTTP(t *testing.T) {
	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.url, nil)
		s, _ := NewFileSystemFiler(tt.dir)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != tt.code {
			t.Errorf("url: %s, want: %d, got: %d", tt.url, tt.code, rec.Code)
		}
	}
}
