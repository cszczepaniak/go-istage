package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cszczepaniak/go-istage/logging"
)

func (v view) handleStateChange(msg tea.Msg) (view, tea.Cmd) {
	event := eventFromMsg(msg)
	if event == UnknownEvent {
		return v, nil
	}

	nextState := v.state.Next(event)
	if nextState == v.state {
		return v, nil
	}

	v.prevState = v.state
	v.state = nextState
	v.currentModel = v.state.Model(v)
	var cmd tea.Cmd
	if v.state != v.prevState {
		// We entered a new state.
		cmd = v.state.OnEnter(v)
	}

	return v, cmd
}

type Event int

func eventFromMsg(msg tea.Msg) Event {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case `t`:
			return ToggleStageEvent
		case `f`:
			return ToggleDiffEvent
		case `c`:
			return StartCommitEvent
		}
	}
	return UnknownEvent
}

const (
	UnknownEvent Event = iota
	ToggleStageEvent
	ToggleDiffEvent
	StartCommitEvent
)

type StateVariant int

const (
	ViewUnstagedLines StateVariant = iota
	ViewUnstagedFiles
	ViewStagedLines
	ViewStagedFiles
	Committing
	Error
)

var stateMap = map[Event]map[StateVariant]StateVariant{
	ToggleStageEvent: {
		ViewUnstagedLines: ViewStagedLines,
		ViewUnstagedFiles: ViewStagedFiles,
		ViewStagedLines:   ViewUnstagedLines,
		ViewStagedFiles:   ViewUnstagedFiles,
		Committing:        Committing,
		Error:             Error,
	},
	ToggleDiffEvent: {
		ViewUnstagedLines: ViewUnstagedFiles,
		ViewUnstagedFiles: ViewUnstagedLines,
		ViewStagedLines:   ViewStagedFiles,
		ViewStagedFiles:   ViewStagedLines,
		Committing:        Committing,
		Error:             Error,
	},
	StartCommitEvent: {
		ViewUnstagedLines: Committing,
		ViewUnstagedFiles: Committing,
		ViewStagedLines:   Committing,
		ViewStagedFiles:   Committing,
		Committing:        Committing,
		Error:             Error,
	},
}

func (sv StateVariant) Next(event Event) StateVariant {
	nextState, ok := stateMap[event][sv]
	if !ok {
		logging.Warn(`unknown next state`, `event`, event, `currState`, sv)
		return sv
	}
	return nextState
}

func (sv StateVariant) Model(v view) tea.Model {
	switch sv {
	case ViewUnstagedLines:
		return v.unstagedLinesView
	case ViewUnstagedFiles:
		return v.unstagedFilesView
	case ViewStagedLines:
		return v.stagedLinesView
	case ViewStagedFiles:
		return v.stagedFilesView
	case Committing:
		return v.commitView
	case Error:
		return v.errorView
	}
	panic(`unreachable`)
}

func (sv StateVariant) OnEnter(v view) tea.Cmd {
	switch sv {
	case ViewUnstagedLines:
		return v.unstagedLinesView.UpdateDoc
	case ViewUnstagedFiles:
		return v.unstagedFilesView.UpdateFiles
	case ViewStagedLines:
		return v.stagedLinesView.UpdateDoc
	case ViewStagedFiles:
		return v.stagedFilesView.UpdateFiles
	case Committing:
		return v.commitView.OnEnter()
	case Error:
		return nil
	}
	panic(`unreachable`)
}
