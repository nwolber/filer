// Copyright (c) 2014 Niklas Wolber
// This file is licensed under the MIT license.
// See the LICENSE file for more information.

package filer

import (
    "net/http"
    "net/http/httptest"
    "strings"

    "testing"
    )

    func TestNewFileSystemFiler(t *testing.T) {
        s, _ := NewFileSystemFiler("")
        if s == nil {
            t.Error("A static server should be creatable at all times.")
        }
    }

    var tests  = []struct {
        dir  string
        url  string
        mime string
        code int
        }{
            {"tests", "../filer_test.go", "", http.StatusForbidden},
            {"tests", "invalid.js", "", http.StatusNotFound},
            {"tests", "test.js", "application/javascript", http.StatusOK},
            // should redirect to index.html
            {"tests", "/", "text/html", http.StatusOK},
            {"tests", "", "text/html", http.StatusOK},
            {"tests", "test.css", "text/css", http.StatusOK},
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

                if mime := rec.Header().Get("Content-Type"); rec.Code == http.StatusOK && !strings.Contains(mime, tt.mime) {
                    t.Errorf("url: %s, want: %s, got: %s,", tt.url, tt.mime, mime)
                }
            }
        }
