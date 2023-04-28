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

func (ps *PatchingService) ApplyPatch(dir patch.Direction, entireHunk bool, selectedLines []int) error {
	if ps.ds.viewFiles {
		return errors.New(`ApplyPatch: viewFiles unimplemented`)
	}

	var lines []int
	if entireHunk {
		return errors.New(`ApplyPatch: entireHunk unimplemented`)
	} else {
		doc := ps.ds.Document
		for _, l := range selectedLines {
			if doc.Lines[l].Kind.IsAdditionOrRemoval() {
				lines = append(lines, l)
			}
		}
	}

	if len(lines) == 0 {
		return nil
	}

	patch, err := patch.Compute(ps.ds.Document, lines, dir)
	if err != nil {
		return err
	}

	return ps.gs.ApplyPatch(patch, dir)
}
