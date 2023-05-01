package services

import (
	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/patch"
)

type gitClient interface {
	StagedChanges() ([]string, error)
	UnstagedChanges() ([]string, error)
	StagedFiles() ([]git.File, error)
	UnstagedFiles() ([]git.File, error)
}

type DocumentService struct {
	gc gitClient

	viewFiles    bool
	viewStage    bool
	fullFileDiff bool
}

func NewDocumentService(gc gitClient) (*DocumentService, error) {
	ds := &DocumentService{
		gc:           gc,
		viewFiles:    false, // TODO support these
		viewStage:    false,
		fullFileDiff: false, // TODO support these
	}

	return ds, nil
}

func (ds *DocumentService) ToggleView() {
	ds.viewStage = !ds.viewStage
}

func (ds *DocumentService) ViewStage() bool {
	return ds.viewStage
}

func (ds *DocumentService) StagedChanges() (patch.Document, error) {
	changes, err := ds.gc.StagedChanges()
	if err != nil {
		return patch.Document{}, err
	}

	return patch.ParseDocument(changes), nil
}

func (ds *DocumentService) UnstagedChanges() (patch.Document, error) {
	changes, err := ds.gc.UnstagedChanges()
	if err != nil {
		return patch.Document{}, err
	}

	return patch.ParseDocument(changes), nil
}

func (ds *DocumentService) UnstagedFiles() ([]git.File, error) {
	return ds.gc.UnstagedFiles()
}

func (ds *DocumentService) StagedFiles() ([]git.File, error) {
	return ds.gc.StagedFiles()
}
