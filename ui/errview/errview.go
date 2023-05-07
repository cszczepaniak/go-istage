package errview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	err error
}

func New() *UI {
	return &UI{}
}

func (*UI) Init() tea.Cmd {
	return nil
}

func (u *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return u, tea.Quit
		case "esc", "enter":
			return u, func() tea.Msg {
				return ExitMsg{}
			}
		}
	case error:
		u.err = msg
	}
	return u, nil
}

var errMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(`#777777`))

func (u *UI) View() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		"An error occurred:",
		errMessageStyle.Render(u.err.Error()),
		"Press enter or escape to continue",
	)
}
