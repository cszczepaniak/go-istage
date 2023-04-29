package ui

import (
	"fmt"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
)

func (v view) stageLine() tea.Msg {
	lineIdx := v.window.AbsoluteIndex(v.cursorLine)
	err := v.patcher.ApplyPatch(patch.Stage, false, []int{lineIdx})
	if err != nil {
		logging.Error(`stageLine failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) stageHunk() (msg tea.Msg) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Errorf("panic recovered: %+v\n%s", r, debug.Stack())
		}
	}()

	lineIdx := v.window.AbsoluteIndex(v.cursorLine)
	h, ok := v.updater.FindHunk(lineIdx)
	if !ok {
		logging.Warn(`stageHunk hunk not found`, `index`, lineIdx)
		return nil
	}

	var lines []int
	for l := h.LineStart(); l < h.LineEnd(); l++ {
		dl := v.doc.Lines[l]
		if dl.Kind.IsAdditionOrRemoval() {
			lines = append(lines, l)
		}
	}

	err := v.patcher.ApplyPatch(patch.Stage, false, lines)
	if err != nil {
		logging.Error(`stageHunk failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) unstageLine() tea.Msg {
	lineIdx := v.window.AbsoluteIndex(v.cursorLine)
	err := v.patcher.ApplyPatch(patch.Unstage, false, []int{lineIdx})
	if err != nil {
		logging.Error(`unstageLine failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) unstageHunk() (msg tea.Msg) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Errorf("panic recovered: %+v\n%s", r, debug.Stack())
		}
	}()

	lineIdx := v.window.AbsoluteIndex(v.cursorLine)
	h, ok := v.updater.FindHunk(lineIdx)
	if !ok {
		logging.Warn(`unstageHunk hunk not found`, `index`, lineIdx)
		return nil
	}

	var lines []int
	for l := h.LineStart(); l < h.LineEnd(); l++ {
		dl := v.doc.Lines[l]
		if dl.Kind.IsAdditionOrRemoval() {
			lines = append(lines, l)
		}
	}

	err := v.patcher.ApplyPatch(patch.Unstage, false, lines)
	if err != nil {
		logging.Error(`unstageHunk failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) updateDoc() tea.Msg {
	doc, err := v.updater.UpdateDocument()
	if err != nil {
		logging.Error(`updateDoc failed`, `err`, err)
		return err
	}
	return docMsg{d: doc}
}

func (v view) cursorLeft() tea.Msg {
	start := v.window.AbsoluteIndex(v.cursorLine)
	if start <= 0 {
		return nil
	}
	for i := start - 1; i >= 0; i-- {
		l := v.doc.Lines[i]
		if l.Kind == patch.HunkLine {
			return jumpToDocLineIndexMsg{
				index: i,
			}
		}
	}
	return jumpToDocLineIndexMsg{index: 0}
}

func (v view) cursorRight() tea.Msg {
	start := v.window.AbsoluteIndex(v.cursorLine)
	if start >= len(v.doc.Lines)-1 {
		return nil
	}
	for i := start + 1; i < len(v.doc.Lines); i++ {
		l := v.doc.Lines[i]
		if l.Kind == patch.HunkLine {
			return jumpToDocLineIndexMsg{
				index: i,
			}
		}
	}
	return jumpToDocLineIndexMsg{index: len(v.doc.Lines) - 1}
}
