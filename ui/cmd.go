package ui

import (
	"fmt"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
)

func (v view) stageLine() tea.Msg {
	if v.viewStage {
		return nil
	}

	err := v.patcher.ApplyPatch(
		patch.Stage,
		v.currentView().doc,
		[]int{v.currentView().currentLine()},
	)
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

	if v.viewStage {
		return nil
	}

	err := v.patcher.ApplyPatch(
		patch.Stage,
		v.currentView().doc,
		v.currentView().linesInCurrentHunk(),
	)
	if err != nil {
		logging.Error(`stageHunk failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) unstageLine() tea.Msg {
	if !v.viewStage {
		return nil
	}

	err := v.patcher.ApplyPatch(
		patch.Unstage,
		v.currentView().doc,
		[]int{v.currentView().currentLine()},
	)
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

	if !v.viewStage {
		return nil
	}

	err := v.patcher.ApplyPatch(
		patch.Unstage,
		v.currentView().doc,
		v.currentView().linesInCurrentHunk(),
	)
	if err != nil {
		logging.Error(`unstageHunk failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) revertLine() (msg tea.Msg) {
	if v.viewStage {
		return nil
	}

	err := v.patcher.ApplyPatch(
		patch.Reset,
		v.currentView().doc,
		[]int{v.currentView().currentLine()},
	)
	if err != nil {
		logging.Error(`revertLine failed`, `err`, err)
		return err
	}
	return refreshMsg{}
}

func (v view) revertHunk() (msg tea.Msg) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Errorf("panic recovered: %+v\n%s", r, debug.Stack())
		}
	}()

	if v.viewStage {
		return nil
	}

	err := v.patcher.ApplyPatch(
		patch.Reset,
		v.currentView().doc,
		v.currentView().linesInCurrentHunk(),
	)
	if err != nil {
		logging.Error(`revertHunk failed`, `err`, err)
		return err
	}
	return refreshMsg{}
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

		return docMsg{d: doc, staged: staged}
	}
}
