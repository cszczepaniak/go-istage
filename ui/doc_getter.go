package ui

import (
	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/patch"
)

type getDocFunc func() (patch.Document, error)

func (f getDocFunc) GetDocument() (patch.Document, error) {
	return f()
}

type getFilesFunc func() ([]git.File, error)

func (f getFilesFunc) GetFiles() ([]git.File, error) {
	return f()
}
