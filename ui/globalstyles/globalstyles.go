package globalstyles

import "github.com/charmbracelet/lipgloss"

var (
	SelectedBackground = lipgloss.NewStyle().Background(lipgloss.Color(`#555555`))
	AdditionColor      = lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FF00`))
	RemovalColor       = lipgloss.NewStyle().Foreground(lipgloss.Color(`#FF0000`))
)
