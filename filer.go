// Copyright (c) 2014 Niklas Wolber
// This file is licensed under the MIT license.
// See the LICENSE file for more information.

package filer

import (
	"bufio"
	"errors"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Filer serves files from a local directory.
type Filer struct {
	dir string
}

var (
	errorForbidden = errors.New("Forbidden")
)

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

// New creates a new filer. It will serve files relative to the current working directory.
func New(d string) (*Filer, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	p := path.Join(wd, d)
	return &Filer{dir: p}, nil
}

func (f *Filer) serveFile(w http.ResponseWriter, r *http.Request) error {
	p := path.Join(f.dir, r.URL.String())

	p = filepath.Clean(p)

	if !strings.HasPrefix(p, f.dir) {
		return errorForbidden
	}

	file, err := os.Open(p)
	if err != nil {
		return err
	}

	defer file.Close()

	d, err := file.Stat()
	if err != nil {
		return err
	}

	if d.IsDir() {
		file, err = os.Open(path.Join(p, "index.html"))
		if err != nil {
			return err
		}
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
