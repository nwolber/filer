// Copyright (c) 2014 Niklas Wolber
// This file is licensed under the MIT license.
// See the LICENSE file for more information.

// Package filer contains a simple HTTP server for static resources.
package filer

import (
	"bufio"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

// A Filer serves static resources.
type Filer struct {
	http.FileSystem
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
func New(fs http.FileSystem) (*Filer, error) {
	return &Filer{fs}, nil
}

func (f *Filer) serveFile(w http.ResponseWriter, r *http.Request) error {
	p := r.URL.String()

	file, err := f.FileSystem.Open(p)
	if err != nil {
		log.Println(err)
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		p = path.Join(p, "index.html")
		file, err = f.FileSystem.Open(p)
		if err != nil {
			return err
		}
	}

	if file == nil {
		return os.ErrNotExist
	}

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
