// Copyright (c) 2014 Niklas Wolber

package filer

import (
	"bufio"
	"mime"
	"net/http"
	"os"
	"path"
)

// Filer serves files from a local directory.
type Filer struct {
	dir string
}

// ServeHTTP returns files from the file system. The file has to be located in the directory configured to the
// static server or a subdirectory. HTTP errors are raised if the requested file does not exist or another error
// occured.
func (f *Filer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := path.Join(f.dir, r.URL.String())
	file, err := os.Open(p)

	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if os.IsPermission(err) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
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
}

// New creates a new filer.
func New(d string) *Filer {
	return &Filer{dir: d}
}
