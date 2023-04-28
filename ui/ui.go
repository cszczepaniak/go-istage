package ui

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/window"
)

func RunUI(doc patch.Document, p patcher, u docUpdater) error {
	f, err := tea.LogToFile(`log/debug.log`, `[LOG] `)
	if err != nil {
		return err
	}
	defer f.Close()

	v := newView(doc, p, u)
	prog := tea.NewProgram(v)
	_, err = prog.Run()
	return err
}

type patcher interface {
	ApplyPatch(dir patch.Direction, entireHunk bool, selectedLines []int) error
}

type docUpdater interface {
	UpdateDocument() (patch.Document, error)
	ToggleView()
	ViewStage() bool
}

type view struct {
	doc     patch.Document
	patcher patcher
	updater docUpdater

	cursorLine int
	h, w       int

	window *window.Window[patch.Line]
}

func newView(doc patch.Document, p patcher, u docUpdater) view {
	return view{
		doc:     doc,
		patcher: p,
		updater: u,
		window:  window.NewWindow(doc.Lines, 0),
	}
}

func (v view) Init() tea.Cmd {
	return nil
}

func (v view) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// I'm not really sure why, but it seems like to make bubbletea render properly, we need to subtract 1 from the
		// height when setting the window.

		v.h = msg.Height - 1
		v.w = msg.Width

		v.window.Resize(msg.Height - 1)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		case "up":
			if v.cursorLine == 0 {
				v.window.ScrollUp()
			} else {
				v.cursorLine--
			}
		case "down":
			if v.cursorLine == v.window.Size()-1 {
				v.window.ScrollDown()
			} else {
				v.cursorLine++
			}
		case "left":
			return v, v.cursorLeft
		case "right":
			return v, v.cursorRight
		case "t":
			v.updater.ToggleView()
			return v, v.updateDoc
		case "s":
			if !v.updater.ViewStage() {
				return v, v.stageLine
			}
		case "u":
			if v.updater.ViewStage() {
				return v, v.unstageLine
			}
		}
	case windowScrollUpMsg:
		v.window.ScrollUp()
	case windowScrollDownMsg:
		v.window.ScrollDown()
	case windowJumpMsg:
		v.cursorLine = msg.index
		v.window.JumpTo(msg.index)
	case refreshMsg:
		return v, v.updateDoc
	case docMsg:
		v.doc = msg.d
	case error:
		log.Println(`ERROR:`, msg)
		return v, tea.Quit
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

func (v view) View() string {
	sb := &strings.Builder{}

	viewableLines := v.window.CurrentValues()

	for i, l := range viewableLines.Values {
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
