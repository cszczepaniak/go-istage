package patch

type entryKey struct {
	Offset int
	Length int
}

// func Compute(doc Document, lineIndices []int, dir Direction) (string, error) {
// 	newPatch := &strings.Builder{}

// 	linesByEntry := map[entryKey][]int{}
// 	entryKeyToEntry := map[entryKey]Entry{}

// 	for _, idx := range lineIndices {
// 		e, ok := doc.FindEntry(idx)
// 		if !ok {
// 			return ``, fmt.Errorf(`dev error: entry not found for line index %d`, idx)
// 		}

// 		k := entryKey{
// 			Offset: e.Offset,
// 			Length: e.Length,
// 		}

// 		linesByEntry[k] = append(linesByEntry[k], idx)
// 		entryKeyToEntry[k] = e
// 	}

// 	for e, idxs := range linesByEntry {
// 		ent, ok := entryKeyToEntry[e]
// 		if !ok {
// 			return ``, fmt.Errorf(`dev error: entry not found for entry key %+v`, e)
// 		}

// 		linesByHunk := map[Hunk][]int{}
// 		for _, idx := range idxs {
// 			h, ok := ent.FindHunk(idx)
// 			if !ok {
// 				fmt.Errorf(`dev error: hunk not found for line index %d`, idx)
// 			}

// 			linesByHunk[h] = append(linesByHunk[h], idx)
// 		}

// 		for hunk, lines := range linesByHunk {
// 			lineSet := newSet(lines...)

// 			oldStart := hunk.OldStart
// 			newStart := hunk.NewStart

// 			oldLength := 0
// 			for i := hunk.Offset; i < hunk.Offset+hunk.Length; i++ {
// 				line := doc.Lines[i]
// 				kind := line.Kind

// 				wasPresent := kind == ContextLine ||
// 					kind == RemovalLine && (!dir.IsUndo() || lineSet.contains(i)) ||
// 					kind == AdditionLine && (dir.IsUndo() && !lineSet.contains(i))

// 				if wasPresent {
// 					oldLength++
// 				}
// 			}

// 			delta := 0
// 			for _, l := range lines {
// 				ln := doc.Lines[l]
// 				if ln.Kind == AdditionLine {
// 					delta += 1
// 				} else if ln.Kind == RemovalLine {
// 					delta -= 1
// 				}
// 			}

// 			newLength := oldLength + delta

// 		}
// 	}
// }

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
