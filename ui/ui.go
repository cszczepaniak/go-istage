package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
)

func RunUI(p patcher, u docUpdater) error {
	v := newView(p, u)
	prog := tea.NewProgram(v)
	_, err := prog.Run()
	return err
}

type patcher interface {
	ApplyPatch(dir patch.Direction, doc patch.Document, selectedLines []int) error
}

type docUpdater interface {
	StagedChanges() (patch.Document, error)
	UnstagedChanges() (patch.Document, error)
}

type view struct {
	patcher patcher
	updater docUpdater

	viewStage    bool
	stagedView   *documentView
	unstagedView *documentView

	h, w int
}

func newView(p patcher, u docUpdater) view {
	return view{
		patcher: p,
		updater: u,
	}
}

func (v view) currentView() *documentView {
	if v.viewStage {
		return v.stagedView
	}
	return v.unstagedView
}

func (v view) Init() tea.Cmd {
	return tea.Batch(v.updateDocs(false), v.updateDocs(true))
}

func (v view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// I'm not really sure why, but it seems like to make bubbletea render properly, we need to subtract 1 from the
		// height when setting the window.

		v.h = msg.Height - 1
		v.w = msg.Width

		if v.stagedView != nil {
			v.stagedView.resize(v.h)
		}
		if v.unstagedView != nil {
			v.unstagedView.resize(v.h)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		case "up", "down", "left", "right":
			if v.currentView() != nil {
				v.currentView().update(msg)
			}
		case "t":
			v.viewStage = !v.viewStage
			return v, v.updateDocs(v.viewStage)
		case "s":
			if !v.viewStage {
				return v, v.stageLine
			}
		case "S":
			if !v.viewStage {
				return v, v.stageHunk
			}
		case "u":
			if v.viewStage {
				return v, v.unstageLine
			}
		case "U":
			if v.viewStage {
				return v, v.unstageHunk
			}
		}
	case windowScrollUpMsg, windowScrollDownMsg, jumpToDocLineIndexMsg:
		v.currentView().update(msg)
	case refreshMsg:
		return v, v.updateDocs(v.viewStage)
	case docMsg:
		if msg.staged && v.stagedView == nil {
			v.stagedView = newDocumentView(msg.d, v.h)
		} else if !msg.staged && v.unstagedView == nil {
			v.unstagedView = newDocumentView(msg.d, v.h)
		}

		if v.currentView() != nil {
			v.currentView().setDoc(msg.d, v.h)
		}
	case error:
		logging.Error(`update.error`, `err`, msg)
		return v, tea.Quit
	}

	return v, nil
}

var kindToColor = map[patch.LineKind]lipgloss.Style{
	patch.AdditionLine: lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FF00`)),
	patch.RemovalLine:  lipgloss.NewStyle().Foreground(lipgloss.Color(`#FF0000`)),
	patch.DiffLine:     lipgloss.NewStyle().Foreground(lipgloss.Color(`#FFFFFF`)),
	patch.HunkLine:     lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FFFF`)),
}

var selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color(`#555555`))

func (v view) View() string {
	return v.currentView().view()
}
