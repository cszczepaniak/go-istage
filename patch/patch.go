package patch

type Direction int

const (
	Stage Direction = iota
	Unstage
	Reset
)

//go:generate stringer -type=LineKind
type LineKind int

const (
	DiffLine LineKind = iota
	HeaderLine
	HunkLine
	ContextLine
	AdditionLine
	RemovalLine
	NoEndOfLineLine
)
