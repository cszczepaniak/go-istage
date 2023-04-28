package patch

type Hunk struct {
	Offset    int
	Length    int
	OldStart  int
	OldLength int
	NewStart  int
	NewLength int
}

func (h Hunk) LineStart() int {
	return h.Offset
}

func (h Hunk) LineEnd() int {
	return h.Offset + h.Length
}

func (h Hunk) ContainsLine(l int) bool {
	return l >= h.LineStart() && l < h.LineEnd()
}
