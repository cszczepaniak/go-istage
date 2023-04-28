package patch

type Entry struct {
	Offset int
	Length int
	Hunks  []Hunk
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
		if h.Offset <= lineIndex && lineIndex <= h.Offset+h.Length {
			return i
		}
	}
	return -1
}
