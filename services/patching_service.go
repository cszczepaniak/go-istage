package services

import (
	"github.com/cszczepaniak/go-istage/patch"
)

type patchClient interface {
	ApplyPatch(string, patch.Direction) error
}

type PatchingService struct {
	pc patchClient
}

func NewPatchingService(pc patchClient) *PatchingService {
	return &PatchingService{
		pc: pc,
	}
}

func (ps *PatchingService) ApplyPatch(dir patch.Direction, doc patch.Document, selectedLines []int) error {
	var lines []int
	for _, l := range selectedLines {
		if doc.Lines[l].Kind.IsAdditionOrRemoval() {
			lines = append(lines, l)
		}
	}

	if len(lines) == 0 {
		return nil
	}

	patch, err := patch.Compute(doc, lines, dir)
	if err != nil {
		return err
	}

	return ps.pc.ApplyPatch(patch, dir)
}
