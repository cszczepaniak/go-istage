package commit

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type UI struct {
	textInput textarea.Model
}

func New() *UI {
	return &UI{
		textInput: textarea.New(),
	}
}

func (u *UI) Init() tea.Cmd { return nil }

func (u *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		u.textInput.SetWidth(msg.Width)
		return u, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			commitMsg := u.textInput.Value()
			u.textInput.Reset()
			return u, u.doCommit(commitMsg)
		}
	}
	newTextInput, cmd := u.textInput.Update(msg)
	u.textInput = newTextInput

	return u, cmd
}

func (u *UI) View() string {
	return fmt.Sprintf("Enter a commit message:\n\n%s\n\n(ctrl+s to commit, escape to abort)\n", u.textInput.View())
}

func (u *UI) OnEnter() tea.Cmd {
	return u.textInput.Focus()
}

func (u *UI) doCommit(msg string) tea.Cmd {
	return func() tea.Msg {
		return DoCommitMsg{
			CommitMessage: msg,
		}
	}
}
