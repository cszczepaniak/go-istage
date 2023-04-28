package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cszczepaniak/go-istage/patch"
)

func (v view) stageLine() tea.Msg {
	err := v.patcher.ApplyPatch(patch.Stage, false, []int{v.cursorLine})
	if err != nil {
		log.Println(`ERROR:`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) unstageLine() tea.Msg {
	err := v.patcher.ApplyPatch(patch.Unstage, false, []int{v.cursorLine})
	if err != nil {
		log.Println(`ERROR:`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) updateDoc() tea.Msg {
	doc, err := v.updater.UpdateDocument()
	if err != nil {
		log.Println(`ERROR:`, err)
		return err
	}
	return docMsg{d: doc}
}

func (v view) cursorLeft() tea.Msg {
	if v.cursorLine <= 0 {
		return nil
	}
	for i := v.cursorLine - 1; i >= 0; i-- {
		l := v.doc.Lines[i]
		if l.Kind == patch.HunkLine {
			return windowJumpMsg{
				index: i,
			}
		}
	}
	return windowJumpMsg{index: 0}
}

func (v view) cursorRight() tea.Msg {
	if v.cursorLine >= len(v.doc.Lines)-1 {
		return nil
	}
	for i := v.cursorLine + 1; i < len(v.doc.Lines); i++ {
		l := v.doc.Lines[i]
		if l.Kind == patch.HunkLine {
			return windowJumpMsg{
				index: i,
			}
		}
	}
	return windowJumpMsg{index: len(v.doc.Lines) - 1}
}
