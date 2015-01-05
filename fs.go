// Copyright (c) 2014 Niklas Wolber
// This file is licensed under the MIT license.
// See the LICENSE file for more information.

package filer

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// NewFileSystemFiler creates a new filer. It will serve files relative to the current working directory.
func NewFileSystemFiler(d string) (*Filer, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	p := path.Join(wd, d)
	return &Filer{a: &fs{dir: p}}, nil
}

type fs struct {
	dir string
}

func (fs *fs) Asset(name string) (io.ReadCloser, error) {
	p := path.Join(fs.dir, name)
	p = filepath.Clean(p)

	if !strings.HasPrefix(p, fs.dir) {
		return nil, errorForbidden
	}

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}

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

	if !strings.HasPrefix(p, fs.dir) {
		return false, errorForbidden
	}

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
