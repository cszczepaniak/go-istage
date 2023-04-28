package services

import (
	"errors"

	"github.com/cszczepaniak/go-istage/patch"
)

type PatchingService struct {
	gs *GitService
	ds *DocumentService
}

func NewPatchingService(gs *GitService, ds *DocumentService) *PatchingService {
	return &PatchingService{
		gs: gs,
		ds: ds,
	}
}

func (ps *PatchingService) ApplyPatch(dir patch.Direction, entireHunk bool, selectedLine int) error {
	if ps.ds.viewFiles {
		return errors.New(`ApplyPatch: viewFiles unimplemented`)
	}

	doc := ps.ds.Document
	line := doc.Lines[selectedLine]

	if !line.Kind.IsAdditionOrRemoval() {
		return nil
	}

	var lines []int
	if entireHunk {
		return errors.New(`ApplyPatch: entireHunk unimplemented`)
	} else {
		lines = []int{selectedLine}
	}

	// TODO compute patch with lines
	_ = lines
	return errors.New(`ApplyPatch: computing patch isn't implemented yet`)
	return ps.gs.ApplyPatch(``, dir)
}
