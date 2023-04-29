package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
)

func RunUI(p patcher, u docUpdater, ge gitExecer) error {
	v := newView(p, u, ge)
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

type gitExecer interface {
	Exec(cmd string) *git.GitExecBuilder
}

type view struct {
	patcher   patcher
	updater   docUpdater
	gitExecer gitExecer

	viewStage    bool
	stagedView   *documentView
	unstagedView *documentView

	committing  bool
	commitInput textarea.Model

	err error

	h, w int
}

func newView(p patcher, u docUpdater, ge gitExecer) view {
	return view{
		commitInput: textarea.New(),
		patcher:     p,
		updater:     u,
		gitExecer:   ge,
	}
}

func (v view) currentView() *documentView {
	if v.viewStage {
		return v.stagedView
	}
	return v.unstagedView
}

func (v view) Init() tea.Cmd {
	return tea.Batch(v.updateDocs(false), v.updateDocs(true), textarea.Blink)
}

func (v view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if v.err != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return v, tea.Quit
			case "esc", "enter":
				v.err = nil
				return v, nil
			}

			return v, nil
		}
	}

	if v.committing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return v, tea.Quit
			case "esc", "enter":
				commitMsg := v.commitInput.Value()
				v.committing = false
				if msg.String() == "enter" {
					v.commitInput.Reset()
					return v, tea.Sequence(
						v.commit(commitMsg),
						tea.Batch(v.updateDocs(false), v.updateDocs(true)),
					)
				}
				return v, nil
			}
			mdl, cmd := v.commitInput.Update(msg)
			v.commitInput = mdl
			return v, cmd
		}
	}

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

		v.commitInput.SetWidth(v.w)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return v, tea.Quit
		case "up":
			v.currentView().navigate(navigateUp)
		case "down":
			v.currentView().navigate(navigateDown)
		case "left":
			v.currentView().navigate(navigateLeft)
		case "right":
			v.currentView().navigate(navigateRight)
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
		case "r":
			return v, v.revertLine
		case "R":
			return v, v.revertHunk
		case "c":
			v.committing = true
			return v, v.commitInput.Focus()
		}
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
		v.err = msg
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
var errMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(`#777777`))

func (v view) View() string {
	if v.err != nil {
		return fmt.Sprintf("%s\n\n%s\n\n%s",
			"An error occurred:",
			errMessageStyle.Render(v.err.Error()),
			"Press enter to continue",
		)
	}
	if v.committing {
		return fmt.Sprintf("Enter a commit message:\n\n%s\n\n%s", v.commitInput.View(), "(enter to commit; esc to abort)")
	}
	return v.currentView().view()
}
