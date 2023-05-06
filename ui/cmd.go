package ui

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/ui/files"
	"github.com/cszczepaniak/go-istage/ui/lines"
)

func (v view) handlePatch(msg lines.PatchMsg) tea.Cmd {
	return func() tea.Msg {
		err := v.patcher.ApplyPatch(msg.Direction, msg.Doc, msg.Lines)
		if err != nil {
			return err
		}
		return lines.RefreshMsg{}
	}
}

func (v view) handleFile(msg files.HandleFileMsg) tea.Cmd {
	return func() tea.Msg {
		var err error
		switch msg.Direction {
		case patch.Stage:
			err = v.fileStager.StageFile(msg.File)
		case patch.Unstage:
			err = v.fileStager.UnstageFile(msg.File)
		default:
			err = errors.New(`unimplemented`)
		}
		if err != nil {
			return err
		}
		return files.RefreshMsg{}
	}
}

func (v view) commit(msg string) tea.Cmd {
	return func() tea.Msg {
		err := v.gitExecer.
			Exec(`commit`).
			WithArgs(`-F`, `-`).
			WithStdin(strings.NewReader(msg)).
			Run()
		if err != nil {
			return err
		}

		return goToStateMsg{
			state: v.prevState,
		}
	}
}
