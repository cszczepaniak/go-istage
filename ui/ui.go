package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/ui/commit"
	"github.com/cszczepaniak/go-istage/ui/errview"
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

	prevState StateVariant
	state     StateVariant

	currentModel tea.Model

	stagedLinesView   *lines.UI
	unstagedLinesView *lines.UI

	stagedFilesView   *files.UI
	unstagedFilesView *files.UI

	commitView *commit.UI

	errorView *errview.UI

	h, w int
}

func newView(p patcher, u docUpdater, ge gitExecer, fs fileStager) view {
	v := view{
		patcher:      p,
		updater:      u,
		gitExecer:    ge,
		fileStager:   fs,
		currentModel: loading.New(),
	}

	v.stagedLinesView = lines.New(
		lines.Staged,
		getDocFunc(v.updater.StagedChanges),
		lines.Config{
			HandleLineKey: unstageLineKey,
			HandleHunkKey: unstageHunkKey,
		},
		v.h,
	)

	v.unstagedLinesView = lines.New(
		lines.Unstaged,
		getDocFunc(v.updater.UnstagedChanges),
		lines.Config{
			HandleLineKey: stageLineKey,
			HandleHunkKey: stageHunkKey,

			CanReset:     true,
			ResetLineKey: resetLineKey,
			ResetHunkKey: resetHunkKey,
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

	v.commitView = commit.New()

	v.errorView = errview.New()

	v.state = ViewUnstagedLines
	v.currentModel = v.state.Model(v)
	return v
}

const (
	stageLineKey   = "s"
	stageHunkKey   = "S"
	unstageLineKey = "u"
	unstageHunkKey = "U"
	resetLineKey   = "r"
	resetHunkKey   = "R"
)

func (v view) Init() tea.Cmd {
	return tea.Batch(
		v.stagedLinesView.Init(),
		v.unstagedLinesView.Init(),
		v.stagedFilesView.Init(),
		v.unstagedFilesView.Init(),
		v.errorView.Init(),
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
		v.commitView.Update(msg)
		v.errorView.Update(msg)
		v.w = msg.Width
		v.h = msg.Height
		return v, nil
	case lines.PatchMsg:
		return v, v.handlePatch(msg)
	case lines.ResetMsg:
		return v, v.handleResetPatch(msg)
	case files.HandleFileMsg:
		return v, v.handleFile(msg)
	case commit.DoCommitMsg:
		return v, v.commit(msg.CommitMessage)
	case errview.ExitMsg:
		// TODO this should be centralized with the other spot we update state.
		v.state = v.prevState
		v.prevState = Error
		v.currentModel = v.state.Model(v)
		return v, v.state.OnEnter(v)
	case goToStateMsg:
		// TODO this should be centralized with the other spot we update state.
		v.prevState = v.state
		v.state = msg.state
		v.currentModel = v.state.Model(v)
		return v, v.state.OnEnter(v)
	case error:
		// TODO this should be centralized with the other spot we update state.
		v.prevState = v.state
		v.state = Error

		v.currentModel = v.state.Model(v)

		_, cmd := v.currentModel.Update(msg)
		return v, cmd
	}

	var cmd tea.Cmd
	v, cmd = v.handleStateChange(msg)
	if cmd != nil {
		return v, cmd
	}

	_, cmd = v.currentModel.Update(msg)
	return v, cmd
}

func (v view) View() string {
	return v.currentModel.View()
}
