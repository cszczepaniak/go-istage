package services

import (
	"github.com/cszczepaniak/go-istage/patch"
	git "github.com/libgit2/git2go/v34"
)

type DocumentService struct {
	gs *GitService

	viewFiles    bool
	viewStage    bool
	fullFileDiff bool
	contextLines int
}

func NewDocumentService(gs *GitService) (*DocumentService, error) {
	ds := &DocumentService{
		gs:           gs,
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
	changes, err := ds.stagedChanges()
	if err != nil {
		return patch.Document{}, err
	}

	return patch.ParseDocument(changes), nil
}

func (ds *DocumentService) UnstagedChanges() (patch.Document, error) {
	changes, err := ds.unstagedChanges()
	if err != nil {
		return patch.Document{}, err
	}

	return patch.ParseDocument(changes), nil
}

func (ds *DocumentService) unstagedChanges() ([]string, error) {
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	opts.Flags |= git.DiffShowUntrackedContent
	opts.Flags |= git.DiffRecurseUntracked

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
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}

	headRef, err := ds.gs.repo.Head()
	if err != nil {
		return nil, err
	}
	commit, err := ds.gs.repo.LookupCommit(headRef.Target())
	if err != nil {
		return nil, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	diff, err := ds.gs.repo.DiffTreeToIndex(tree, nil, &opts)
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
