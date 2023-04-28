package patch

type Direction int

func (d Direction) IsUndo() bool {
	return d == Reset || d == Unstage
}

const (
	Stage Direction = iota
	Unstage
	Reset
)

//go:generate stringer -type=LineKind
type LineKind int

func (lk LineKind) IsAdditionOrRemoval() bool {
	return lk == AdditionLine || lk == RemovalLine
}

const (
	DiffLine LineKind = iota
	HeaderLine
	HunkLine
	ContextLine
	AdditionLine
	RemovalLine
	NoEndOfLineLine
)
