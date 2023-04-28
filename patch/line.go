package patch

import "fmt"

type Line struct {
	Kind      LineKind
	Text      string
	LineBreak string
}

func (l Line) String() string {
	return fmt.Sprintf(`[%-12s] %s%s`, l.Kind, l.Text, l.LineBreak)
}
