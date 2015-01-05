// Copyright (c) 2014 Niklas Wolber
// This file is licensed under the MIT license.
// See the LICENSE file for more information.

// Package filer contains a simple HTTP server for static resources.
package filer

import (
	"bufio"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
)

// A Filer serves static resources.
type Filer struct {
	a Asseter
}

var (
	errorForbidden = errors.New("Forbidden")
	errorNoFile    = errors.New("No file")
)

// An Asseter resolves resources by name.
type Asseter interface {
	IsDir(name string) (bool, error)
	Asset(name string) (io.ReadCloser, error)
}

// ServeHTTP returns files from the file system. The file has to be located in the directory configured to the
// static server or a subdirectory. HTTP errors are raised if the requested file does not exist or another error
// occured.
func (f *Filer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f.serveFile(w, r)

	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if os.IsPermission(err) || err == errorForbidden {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// New creates a new filer which uses the Asseter to retrieve assets.
func New(a Asseter) (*Filer, error) {
	return &Filer{a: a}, nil
}

func (f *Filer) serveFile(w http.ResponseWriter, r *http.Request) error {
	var file io.ReadCloser
	var err error

	p := r.URL.String()

	if d, err := f.a.IsDir(p); d && err == nil {
		p = path.Join(p, "index.html")
	} else if err != nil {
		return err
	}

	if file, err = f.a.Asset(p); err != nil {
		return os.ErrNotExist
	}
	defer file.Close()

	_, haveType := w.Header()["Content-Type"]
	if !haveType {
		if t := mime.TypeByExtension(path.Ext(p)); t != "" {
			w.Header().Set("Content-Type", t)
		}
	}

	reader := bufio.NewReader(file)
	reader.WriteTo(w)
	return nil
}
