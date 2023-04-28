package patch

import (
	"fmt"
	"strings"
)

type entryKey struct {
	Offset int
	Length int
}

func Compute(doc Document, lineIndices []int, dir Direction) (string, error) {
	newPatch := &strings.Builder{}

	linesByEntry := map[entryKey][]int{}
	entryKeyToEntry := map[entryKey]Entry{}

	for _, idx := range lineIndices {
		e, ok := doc.FindEntry(idx)
		if !ok {
			return ``, fmt.Errorf(`dev error: entry not found for line index %d`, idx)
		}

		k := entryKey{
			Offset: e.Offset,
			Length: e.Length,
		}

		linesByEntry[k] = append(linesByEntry[k], idx)
		entryKeyToEntry[k] = e
	}

	for e, idxs := range linesByEntry {
		ent, ok := entryKeyToEntry[e]
		if !ok {
			return ``, fmt.Errorf(`dev error: entry not found for entry key %+v`, e)
		}

		linesByHunk := map[Hunk][]int{}
		for _, idx := range idxs {
			h, ok := ent.FindHunk(idx)
			if !ok {
				fmt.Errorf(`dev error: hunk not found for line index %d`, idx)
			}

			linesByHunk[h] = append(linesByHunk[h], idx)
		}

		for hunk, lines := range linesByHunk {
			lineSet := newSet(lines...)

			oldStart := hunk.OldStart
			newStart := hunk.NewStart

			oldLength := 0
			for i := hunk.Offset; i < hunk.Offset+hunk.Length; i++ {
				line := doc.Lines[i]
				kind := line.Kind

				wasPresent := kind == ContextLine ||
					kind == RemovalLine && (!dir.IsUndo() || lineSet.contains(i)) ||
					kind == AdditionLine && (dir.IsUndo() && !lineSet.contains(i))

				if wasPresent {
					oldLength++
				}
			}

			delta := 0
			for _, l := range lines {
				ln := doc.Lines[l]
				if ln.Kind == AdditionLine {
					delta += 1
				} else if ln.Kind == RemovalLine {
					delta -= 1
				}
			}

			newLength := oldLength + delta

			changes := ent.Changes
			oldPath := changes.OldPath
			oldExists := oldLength != 0 || changes.OldMode != ``
			path := changes.Path

			if oldExists {
				fmt.Fprintf(newPatch, "--- a/%s\n", oldPath)
			} else {
				fmt.Fprintf(newPatch, "new file mode %s\n", changes.Mode)
				fmt.Fprintf(newPatch, "--- /dev/null\n", oldPath)
			}

			fmt.Fprintf(newPatch, "+++ b/%s\n", path)

			fmt.Fprint(newPatch, `@@ -`)
			fmt.Fprintf(newPatch, `%d`, oldStart)
			if oldLength != 1 {
				fmt.Fprintf(newPatch, `,%d`, oldLength)
			}

			fmt.Fprintf(newPatch, ` +%d`, newStart)
			if newLength != 1 {
				fmt.Fprintf(newPatch, `,%d`, newLength)
			}

			fmt.Fprintln(newPatch, ` @@`)

			previousIncluded := false
			for i := hunk.Offset; i < hunk.Offset+hunk.Length; i++ {
				line := doc.Lines[i]
				kind := line.Kind

				if lineSet.contains(i) || previousIncluded && kind == NoEndOfLineLine {
					newPatch.WriteString(line.Text)
					newPatch.WriteString(line.LineBreak)
					previousIncluded = true
				} else if !dir.IsUndo() && kind == RemovalLine || dir.IsUndo() && kind == AdditionLine {
					newPatch.WriteString(` `)
					newPatch.WriteString(line.Text[1:len(line.Text)])
					newPatch.WriteString(line.LineBreak)
					previousIncluded = true
				} else {
					previousIncluded = false
				}
			}
		}
	}

	return newPatch.String(), nil
}

type set[K comparable] map[K]struct{}

func newSet[K comparable](from ...K) set[K] {
	s := make(set[K], len(from))
	for _, val := range from {
		s[val] = struct{}{}
	}
	return s
}

func (s set[K]) contains(val K) bool {
	_, ok := s[val]
	return ok
}
