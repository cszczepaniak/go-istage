package loading

import tea "github.com/charmbracelet/bubbletea"

type LoadingView struct{}

func New() LoadingView {
	return LoadingView{}
}

func (LoadingView) Init() tea.Cmd {
	return nil
}

func (lv LoadingView) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return lv, nil
}

func (lv LoadingView) View() string {
	return `Loading...`
}
