package commit

import (
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
	}
	newTextInput, cmd := u.textInput.Update(msg)
	u.textInput = newTextInput

	return u, cmd
}

func (u *UI) View() string {
	return fmt.Sprintf("Enter a commit message:\n\n%s\n\n(Enter to commit, escape to abort)\n", u.textInput.View())
}

func (u *UI) OnEnter() tea.Cmd {
	return u.textInput.Focus()
}
