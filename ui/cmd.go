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
	err := v.patcher.ApplyPatch(patch.Stage, v.unstaged, []int{lineIdx})
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
	h, ok := findHunk(v.unstaged, lineIdx)
	if !ok {
		logging.Warn(`stageHunk hunk not found`, `index`, lineIdx)
		return nil
	}

	var lines []int
	for l := h.LineStart(); l < h.LineEnd(); l++ {
		dl := v.unstaged.Lines[l]
		if dl.Kind.IsAdditionOrRemoval() {
			lines = append(lines, l)
		}
	}

	err := v.patcher.ApplyPatch(patch.Stage, v.unstaged, lines)
	if err != nil {
		logging.Error(`stageHunk failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) unstageLine() tea.Msg {
	lineIdx := v.window.AbsoluteIndex(v.cursorLine)
	err := v.patcher.ApplyPatch(patch.Unstage, v.staged, []int{lineIdx})
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
	h, ok := findHunk(v.staged, lineIdx)
	if !ok {
		logging.Warn(`unstageHunk hunk not found`, `index`, lineIdx)
		return nil
	}

	var lines []int
	for l := h.LineStart(); l < h.LineEnd(); l++ {
		dl := v.staged.Lines[l]
		if dl.Kind.IsAdditionOrRemoval() {
			lines = append(lines, l)
		}
	}

	err := v.patcher.ApplyPatch(patch.Unstage, v.staged, lines)
	if err != nil {
		logging.Error(`unstageHunk failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func findHunk(doc patch.Document, idx int) (patch.Hunk, bool) {
	e, ok := doc.FindEntry(idx)
	if !ok {
		return patch.Hunk{}, false
	}

	return e.FindHunk(idx)
}

func (v view) updateDocs(staged bool) tea.Cmd {
	return func() tea.Msg {
		var doc patch.Document
		var err error

		if v.viewStage {
			doc, err = v.updater.StagedChanges()
		} else {
			doc, err = v.updater.UnstagedChanges()
		}
		if err != nil {
			logging.Error(`updateDocs failed`, `err`, err)
			return err
		}

		return docMsg{d: doc}
	}
}

func (v view) cursorLeft() tea.Msg {
	start := v.window.AbsoluteIndex(v.cursorLine)
	if start <= 0 {
		return nil
	}

	for i := start - 1; i >= 0; i-- {
		l := v.currentDoc().Lines[i]
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
	if start >= len(v.currentDoc().Lines)-1 {
		return nil
	}
	for i := start + 1; i < len(v.currentDoc().Lines); i++ {
		l := v.currentDoc().Lines[i]
		if l.Kind == patch.HunkLine {
			return jumpToDocLineIndexMsg{
				index: i,
			}
		}
	}
	return jumpToDocLineIndexMsg{index: len(v.currentDoc().Lines) - 1}
}
