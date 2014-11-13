// Copyright (c) 2014 Niklas Wolber

package filer

import "testing"

func TestNew(t *testing.T) {
	if New("") == nil {
		t.Error("A static server should be creatable at all times.")
	}
}
