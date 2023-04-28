package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cszczepaniak/go-istage/patch"
)

func (v view) stageLine() tea.Msg {
	err := v.patcher.ApplyPatch(patch.Stage, false, []int{v.cursorLine})
	if err != nil {
		return err
	}
	return refreshMsg{}
}

func (v view) updateDoc() tea.Msg {
	doc, err := v.updater.UpdateDocument()
	if err != nil {
		return err
	}
	return docMsg{d: doc}
}
