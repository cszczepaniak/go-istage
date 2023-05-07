//go:build !prod

package testutils

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func InitializeModel[T tea.Model](t testing.TB, m T) T {
	var mm tea.Model = m
	mm = RunUpdateCycle[T](mm, mm.Init())
	m = RunUpdateCycle[T](mm.Update(tea.WindowSizeMsg{
		Width:  100,
		Height: 40,
	}))
	return m
}

func RunUpdateCycle[T tea.Model](m tea.Model, cmd tea.Cmd) T {
	var mm tea.Model
	for cmd != nil {
		mm, cmd = m.Update(cmd())
		m = mm.(T)
	}
	return m.(T)
}

func ExecKeyPressCycle[T tea.Model](m T, key string) T {
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(key),
	}
	return RunUpdateCycle[T](m.Update(msg))
}

func ExecKeyPress[T tea.Model](m T, key string) (T, tea.Msg) {
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(key),
	}
	mm, cmd := m.Update(msg)

	m = mm.(T)
	if cmd == nil {
		return m, nil
	}

	return m, cmd()
}
