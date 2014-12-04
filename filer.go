// Copyright (c) 2014 Niklas Wolber
// This file is licensed under the MIT license.
// See the LICENSE file for more information.

// Package filer contains a simple HTTP server for static resources.
package filer

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	Asset(name string) (io.Reader, error)
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

// NewFileSystemFiler creates a new filer. It will serve files relative to the current working directory.
func NewFileSystemFiler(d string) (*Filer, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	p := path.Join(wd, d)
	return &Filer{a: &fs{dir: p}}, nil
}

func (f *Filer) serveFile(w http.ResponseWriter, r *http.Request) error {
	var file io.Reader
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

	reader := bufio.NewReader(file)
	reader.WriteTo(w)
	return nil
}

type fs struct {
	dir string
}

func (fs *fs) Asset(name string) (io.Reader, error) {
	p := path.Join(fs.dir, name)

	p = filepath.Clean(p)

	if !strings.HasPrefix(p, fs.dir) {
		return nil, errorForbidden
	}

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	d, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if d.IsDir() {
		return nil, errorNoFile
	}

	return file, nil
}

func (fs *fs) IsDir(name string) (bool, error) {
	p := path.Join(fs.dir, name)
	p = filepath.Clean(p)

	file, err := os.Open(p)
	if err != nil {
		return false, err
	}

	defer file.Close()

	d, err := file.Stat()
	if err != nil {
		return false, err
	}

	return d.IsDir(), nil
}
