package patch

type Direction int

const (
	Stage Direction = iota
	Unstage
	Reset
)
