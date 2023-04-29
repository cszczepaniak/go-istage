package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/window"
)

type documentView struct {
	doc    patch.Document
	w      *window.Window[patch.Line]
	cursor int
}

func newDocumentView(doc patch.Document, windowSize int) *documentView {
	return &documentView{
		doc: doc,
		w:   window.NewWindow(doc.Lines, windowSize),
	}
}

func (dv *documentView) resize(size int) {
	dv.w.Resize(size)
}

func (dv *documentView) setDoc(doc patch.Document, h int) {
	dv.doc = doc
	if dv.w == nil {
		dv.w = window.NewWindow(doc.Lines, h)
	} else {
		dv.w.SetData(doc.Lines)
	}
}

func (dv *documentView) view() string {
	if dv == nil || dv.w == nil {
		return ``
	}

	sb := &strings.Builder{}

	viewableLines := dv.w.CurrentValues()

	for i, l := range viewableLines.Values {
		s := lipgloss.NewStyle()
		c, ok := kindToColor[l.Kind]
		if ok {
			s = s.Inherit(c)
		}
		if dv.cursor == i {
			s = s.Inherit(selectedStyle)
		}

		sb.WriteString(s.Render(l.Text))
		sb.WriteString(l.LineBreak)
	}
	return sb.String()
}

func (dv *documentView) update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if dv.cursor == 0 {
				dv.w.ScrollUp()
			} else {
				dv.cursor--
			}
		case "down":
			if dv.cursor == dv.w.Size()-1 {
				dv.w.ScrollDown()
			} else {
				dv.cursor++
			}
		case "left":
			dv.cursorLeft()
		case "right":
			dv.cursorRight()
		}
	case windowScrollUpMsg:
		dv.w.ScrollUp()
	case windowScrollDownMsg:
		dv.w.ScrollDown()
	case jumpToDocLineIndexMsg:
		relIndex := dv.w.RelativeIndex(msg.index)
		if relIndex < 0 {
			dv.w.JumpTo(msg.index)
			relIndex = dv.w.RelativeIndex(msg.index)
		}
		dv.cursor = relIndex
	}
}

func (dv *documentView) cursorLeft() {
	start := dv.w.AbsoluteIndex(dv.cursor)
	if start <= 0 {
		return
	}

	for i := start - 1; i >= 0; i-- {
		l := dv.doc.Lines[i]
		if l.Kind == patch.HunkLine {
			dv.jumpToLine(i)
			return
		}
	}
	dv.jumpToLine(0)
}

func (dv *documentView) cursorRight() {
	start := dv.w.AbsoluteIndex(dv.cursor)
	if start >= len(dv.doc.Lines)-1 {
		return
	}
	for i := start + 1; i < len(dv.doc.Lines); i++ {
		l := dv.doc.Lines[i]
		if l.Kind == patch.HunkLine {
			dv.jumpToLine(i)
			return
		}
	}
	dv.jumpToLine(len(dv.doc.Lines) - 1)
}

func (dv *documentView) jumpToLine(index int) {
	relIndex := dv.w.RelativeIndex(index)
	if relIndex < 0 {
		dv.w.JumpTo(index)
		relIndex = dv.w.RelativeIndex(index)
	}
	dv.cursor = relIndex
}

func (dv *documentView) currentLine() int {
	return dv.w.AbsoluteIndex(dv.cursor)
}

func (dv *documentView) linesInCurrentHunk() []int {
	lineIdx := dv.currentLine()
	h, ok := findHunk(dv.doc, lineIdx)
	if !ok {
		logging.Warn(`hunk not found`, `index`, lineIdx)
		return nil
	}

	var lines []int
	for l := h.LineStart(); l < h.LineEnd(); l++ {
		dl := dv.doc.Lines[l]
		if dl.Kind.IsAdditionOrRemoval() {
			lines = append(lines, l)
		}
	}

	return lines
}

func findHunk(doc patch.Document, idx int) (patch.Hunk, bool) {
	e, ok := doc.FindEntry(idx)
	if !ok {
		return patch.Hunk{}, false
	}

	return e.FindHunk(idx)
}
