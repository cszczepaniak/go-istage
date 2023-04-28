package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/patch"
)

func RunUI(doc patch.Document, p patcher, u docUpdater) error {
	v := newView(doc, p, u)
	prog := tea.NewProgram(v)
	_, err := prog.Run()
	return err
}

type patcher interface {
	ApplyPatch(dir patch.Direction, entireHunk bool, selectedLines []int) error
}

type docUpdater interface {
	UpdateDocument() (patch.Document, error)
}

type view struct {
	doc     patch.Document
	patcher patcher
	updater docUpdater

	cursorLine int
	h, w       int
}

func newView(doc patch.Document, p patcher, u docUpdater) view {
	return view{
		doc:     doc,
		patcher: p,
		updater: u,
	}
}

func (v view) Init() tea.Cmd {
	return nil
}

func (v view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.h = msg.Height
		v.w = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		case "up":
			if v.cursorLine > 0 {
				v.cursorLine--
			}
		case "down":
			if v.cursorLine < len(v.doc.Lines)-1 {
				v.cursorLine++
			}
		case "s":
			return v, v.stageLine
		}
	case refreshMsg:
		return v, v.updateDoc
	case docMsg:
		v.doc = msg.d
	case error:
		return v, tea.Quit
	}

	return v, nil
}

var kindToColor = map[patch.LineKind]lipgloss.Style{
	patch.AdditionLine: lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FF00`)),
	patch.RemovalLine:  lipgloss.NewStyle().Foreground(lipgloss.Color(`#FF0000`)),
}

var selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color(`#555555`))

func (v view) View() string {
	sb := &strings.Builder{}

	numLines := v.h
	if numLines > len(v.doc.Lines) {
		numLines = len(v.doc.Lines)
	}

	for i, l := range v.doc.Lines[:numLines] {
		s := lipgloss.NewStyle()
		c, ok := kindToColor[l.Kind]
		if ok {
			s = s.Inherit(c)
		}
		if v.cursorLine == i {
			s = s.Inherit(selectedStyle)
		}

		sb.WriteString(s.Render(l.Text))
		sb.WriteString(l.LineBreak)
	}
	return sb.String()
}
