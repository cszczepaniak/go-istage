package patch

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func ParseDocument(patches []string) Document {
	if len(patches) == 0 {
		return Document{}
	}

	var d Document

	for _, p := range patches {
		changeLines := ParseLines(p)
		entryOffset := len(d.Lines)
		entryLength := len(changeLines)
		d.Lines = append(d.Lines, changeLines...)

		var hunks []Hunk
		hunkOffset := GetNextHunk(changeLines, -1)
		for hunkOffset < len(changeLines) {
			hunkEnd := GetNextHunk(changeLines, hunkOffset)
			hunkLength := hunkEnd - hunkOffset + 1
			hunkLine := changeLines[hunkOffset].Text

			info, ok := TryGetHunkInformation(hunkLine)
			if ok {
				hunk := Hunk{
					Offset:    entryOffset + hunkOffset,
					Length:    hunkLength,
					OldStart:  info.oldStart,
					OldLength: info.oldLength,
					NewStart:  info.newStart,
					NewLength: info.newLength,
				}
				hunks = append(hunks, hunk)
			}

			hunkOffset = hunkEnd + 1
		}

		entry := Entry{
			Offset: entryOffset,
			Length: entryLength,
			Hunks:  hunks,
		}
		entry.Changes = ChangesFromHeader(changeLines)

		d.Entries = append(d.Entries, entry)
	}

	return d
}

func GetNextHunk(lines []Line, index int) int {
	index++

	for index < len(lines) && lines[index].Kind != HunkLine {
		index++
	}

	return index
}

type hunkInfo struct {
	oldStart  int
	oldLength int
	newStart  int
	newLength int
}

// Hunk information looks like this:
// @@ -0,0 +1,29 @@
func TryGetHunkInformation(hunkLine string) (hunkInfo, bool) {
	rngs := strings.FieldsFunc(hunkLine, func(r rune) bool {
		return r == ' ' || r == '@'
	})
	if len(rngs) != 2 {
		return hunkInfo{}, false
	}

	oldRng, ok := TryParseRange(rngs[0], `-`)
	if !ok {
		return hunkInfo{}, false
	}

	newRng, ok := TryParseRange(rngs[1], `+`)
	if !ok {
		return hunkInfo{}, false
	}

	return hunkInfo{
		oldStart:  oldRng.start,
		oldLength: oldRng.length,
		newStart:  newRng.start,
		newLength: newRng.length,
	}, true
}

type rng struct {
	start  int
	length int
}

func TryParseRange(s, marker string) (rng, bool) {
	var r rng

	if !strings.HasPrefix(s, marker) {
		return r, false
	}

	num1, num2, ok := strings.Cut(s[1:], `,`)
	start, err := strconv.Atoi(num1)
	if err != nil {
		return rng{}, false
	}
	r.start = start

	if !ok {
		r.length = 1
	} else {
		length, err := strconv.Atoi(num2)
		if err != nil {
			return rng{}, false
		}
		r.length = length
	}

	return r, true
}

func ParseLines(content string) []Line {
	lines := GetLines(content)
	res := make([]Line, 0, len(lines))

	i := 0
	for i < len(lines) {
		l := lines[i]

		if strings.HasPrefix(l.text, `@@`) {
			break
		}

		kind := HeaderLine
		if strings.HasPrefix(l.text, `diff --git`) {
			kind = DiffLine
		}

		res = append(res, Line{
			Kind:      kind,
			Text:      l.text,
			LineBreak: l.lineBreak,
		})

		i++
	}

	for _, l := range lines[i:] {
		kind := ContextLine
		switch {
		case strings.HasPrefix(l.text, `@@`):
			kind = HunkLine
		case strings.HasPrefix(l.text, `+`):
			kind = AdditionLine
		case strings.HasPrefix(l.text, `-`):
			kind = RemovalLine
		case strings.HasPrefix(l.text, `\`):
			kind = NoEndOfLineLine
		}

		res = append(res, Line{
			Kind:      kind,
			Text:      l.text,
			LineBreak: l.lineBreak,
		})
	}

	return res
}

type line struct {
	text      string
	lineBreak string
}

func GetLines(content string) []line {
	var res []line

	beginningOfLine := 0
	i := 0

	commitLine := func(beginningOfLine, i, lineBreakWidth int) int {
		if beginningOfLine >= len(content) {
			return beginningOfLine
		}

		length := i - beginningOfLine
		lineContent := content[beginningOfLine : beginningOfLine+length]
		lineBreak := content[i : i+lineBreakWidth]
		res = append(res, line{
			text:      lineContent,
			lineBreak: lineBreak,
		})
		return i + lineBreakWidth
	}

	const (
		cr = '\r'
		lf = '\n'
	)

	for i < len(content) {
		r, sz := utf8.DecodeRuneInString(content[i:])
		var next rune
		if i+sz < len(content) {
			next, _ = utf8.DecodeRuneInString(content[i+sz:])
		}

		lineBreakWidth := 0
		if r == cr && next == lf {
			lineBreakWidth = 2
		} else if r == cr || r == lf {
			lineBreakWidth = 1
		}

		if lineBreakWidth > 0 {
			beginningOfLine = commitLine(beginningOfLine, i, lineBreakWidth)
			i += lineBreakWidth
		} else {
			i += sz
		}
	}

	commitLine(beginningOfLine, i, 0)

	return res
}
