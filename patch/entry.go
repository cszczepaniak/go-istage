package patch

type Entry struct {
	Offset  int
	Length  int
	Hunks   []Hunk
	Changes Changes
}

func (pe Entry) LineStart() int {
	return pe.Offset
}

func (pe Entry) LineEnd() int {
	return pe.Offset + pe.Length
}

func (pe Entry) ContainsLine(l int) bool {
	return l >= pe.LineStart() && l < pe.LineEnd()
}

func (pe Entry) FindHunk(lineIndex int) (Hunk, bool) {
	hunkIndex := pe.FindHunkIndex(lineIndex)
	if hunkIndex <= -1 {
		return Hunk{}, false
	}
	return pe.Hunks[hunkIndex], true
}

func (pe Entry) FindHunkIndex(lineIndex int) int {
	for i, h := range pe.Hunks {
		if h.ContainsLine(lineIndex) {
			return i
		}
	}
	return -1
}
