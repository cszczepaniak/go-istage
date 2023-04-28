package services

import (
	"errors"
	"log"

	"github.com/cszczepaniak/go-istage/patch"
	git "github.com/libgit2/git2go/v34"
)

type DocumentService struct {
	gs *GitService

	viewFiles    bool
	viewStage    bool
	fullFileDiff bool
	contextLines int

	Document patch.Document
}

func NewDocumentService(gs *GitService) (*DocumentService, error) {
	ds := &DocumentService{
		gs:           gs,
		viewFiles:    false, // TODO support these
		viewStage:    false, // TODO support these
		fullFileDiff: false, // TODO support these
	}

	ds.gs.OnRepoChanged(func() {
		err := ds.UpdateDocument()
		if err != nil {
			log.Println("ERROR: update document failed:", err)
		}
	})

	err := ds.UpdateDocument()
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (ds *DocumentService) UpdateDocument() error {
	var changes []string
	var err error
	if ds.viewStage {
		changes, err = ds.stagedChanges()
	} else {
		changes, err = ds.unstagedChanges()
	}
	if err != nil {
		return err
	}

	if ds.viewFiles {
		return errors.New(`UpdateDocument: viewFiles unimplemented`)
	}

	ds.Document = patch.ParseDocument(changes)
	return nil
}

func (ds *DocumentService) unstagedChanges() ([]string, error) {
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	opts.Flags |= git.DiffShowUntrackedContent

	diff, err := ds.gs.repo.DiffIndexToWorkdir(nil, &opts)
	if err != nil {
		return nil, err
	}

	ndl, err := diff.NumDeltas()
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, ndl)
	for i := 0; i < ndl; i++ {
		p, err := diff.Patch(i)
		if err != nil {
			return nil, err
		}

		txt, err := p.String()
		if err != nil {
			return nil, err
		}

		if txt != `` {
			res = append(res, txt)
		}
	}

	return res, nil
}

func (ds *DocumentService) stagedChanges() ([]string, error) {
	return nil, errors.New(`stagedChanges: unimplemented`)
}
