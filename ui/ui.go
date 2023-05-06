package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/ui/files"
	"github.com/cszczepaniak/go-istage/ui/lines"
	"github.com/cszczepaniak/go-istage/ui/loading"
)

func RunUI(p patcher, u docUpdater, ge gitExecer, fs fileStager) error {
	v := newView(p, u, ge, fs)
	prog := tea.NewProgram(v)
	_, err := prog.Run()
	return err
}

type patcher interface {
	ApplyPatch(dir patch.Direction, doc patch.Document, selectedLines []int) error
}

type fileStager interface {
	StageFile(file git.File) error
	UnstageFile(file git.File) error
}

type docUpdater interface {
	StagedChanges() (patch.Document, error)
	UnstagedChanges() (patch.Document, error)
	StagedFiles() ([]git.File, error)
	UnstagedFiles() ([]git.File, error)
}

type gitExecer interface {
	Exec(cmd string) *git.GitExecBuilder
}

type view struct {
	patcher    patcher
	updater    docUpdater
	gitExecer  gitExecer
	fileStager fileStager

	state StateVariant

	currentModel tea.Model

	stagedLinesView   *lines.UI
	unstagedLinesView *lines.UI

	stagedFilesView   *files.UI
	unstagedFilesView *files.UI

	committing  bool
	commitInput textarea.Model

	err error

	h, w int
}

func newView(p patcher, u docUpdater, ge gitExecer, fs fileStager) view {
	v := view{
		commitInput:  textarea.New(),
		patcher:      p,
		updater:      u,
		gitExecer:    ge,
		fileStager:   fs,
		currentModel: loading.New(),
	}

	v.stagedLinesView = lines.New(
		lines.Staged,
		getDocFunc(v.updater.StagedChanges),
		lines.KeyConfig{
			HandleLineKey: unstageLineKey,
			HandleHunkKey: unstageHunkKey,
		},
		v.h,
	)

	v.unstagedLinesView = lines.New(
		lines.Unstaged,
		getDocFunc(v.updater.UnstagedChanges),
		lines.KeyConfig{
			HandleLineKey: stageLineKey,
			HandleHunkKey: stageHunkKey,
		},
		v.h,
	)

	v.stagedFilesView = files.NewView(
		files.Staged,
		files.KeyConfig{
			HandleFileKey: unstageLineKey,
		},
		getFilesFunc(v.updater.StagedFiles),
		v.h,
	)

	v.unstagedFilesView = files.NewView(
		files.Unstaged,
		files.KeyConfig{
			HandleFileKey: stageLineKey,
		},
		getFilesFunc(v.updater.UnstagedFiles),
		v.h,
	)

	v.state = ViewUnstagedLines
	v.currentModel = v.state.Model(v)
	return v
}

const (
	stageLineKey   = "s"
	stageHunkKey   = "S"
	unstageLineKey = "u"
	unstageHunkKey = "U"
)

func (v view) Init() tea.Cmd {
	return tea.Batch(
		v.stagedLinesView.Init(),
		v.unstagedLinesView.Init(),
		v.stagedFilesView.Init(),
		v.unstagedFilesView.Init(),
		textarea.Blink,
	)
}

func (v view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return v, tea.Quit
		}
	case tea.WindowSizeMsg:
		// We need to let everybody know right away about the window size.
		v.stagedFilesView.Update(msg)
		v.unstagedFilesView.Update(msg)
		v.stagedLinesView.Update(msg)
		v.unstagedLinesView.Update(msg)
		v.w = msg.Width
		v.h = msg.Height
		return v, nil
	case lines.PatchMsg:
		return v, v.handlePatch(msg)
	case files.HandleFileMsg:
		return v, v.handleFile(msg)
	}

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

	var cmd tea.Cmd
	v, cmd = v.handleStateChange(msg)
	if cmd != nil {
		return v, cmd
	}

	if v.committing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "enter":
				commitMsg := v.commitInput.Value()
				v.committing = false
				if msg.String() == "enter" {
					v.commitInput.Reset()
					return v, tea.Sequence(
						v.commit(commitMsg),
						//tea.Batch(v.updateDocs(false), v.updateDocs(true)),
					)
				}
				return v, nil
			}
			mdl, cmd := v.commitInput.Update(msg)
			v.commitInput = mdl
			return v, cmd
		}
	}

	_, cmd = v.currentModel.Update(msg)
	return v, cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.w = msg.Width
		v.commitInput.SetWidth(v.w)
	case tea.KeyMsg:
		switch msg.String() {
		// case "r":
		// 	return v, v.revertLine
		// case "R":
		// 	return v, v.revertHunk
		case "c":
			v.committing = true
			return v, v.commitInput.Focus()
		}
	case filesMsg:
		// if v.unstagedFilesView == nil {
		// 	v.unstagedFilesView = newFileView(msg.files, v.h)
		// }

		// if v.unstagedFilesView != nil {
		// 	v.unstagedFilesView.setFiles(msg.files, v.h)
		// }
	case error:
		logging.Error(`update.error`, `err`, msg)
		v.err = msg
	}

	return v, nil
}

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

	return v.currentModel.View()
}
