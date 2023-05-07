package lines

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/ui/globalstyles"
	"github.com/cszczepaniak/go-istage/window"
)

type DocType int

const (
	Unstaged DocType = iota
	Staged
)

type docGetter interface {
	GetDocument() (patch.Document, error)
}

type UI struct {
	doc       patch.Document
	docType   DocType
	docGetter docGetter

	keyCfg Config

	window *window.Window[patch.Line]
	cursor int
	h      int
}

type Config struct {
	HandleLineKey string
	HandleHunkKey string

	CanReset     bool
	ResetLineKey string
	ResetHunkKey string
}

func New(dt DocType, dg docGetter, keyCfg Config, windowSize int) *UI {
	return &UI{
		docType:   dt,
		docGetter: dg,
		keyCfg:    keyCfg,
		h:         windowSize,
	}
}

func (u *UI) Init() tea.Cmd {
	return u.UpdateDoc
}

func (u *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// I'm not really sure why, but it seems like to make bubbletea render properly, we need to subtract 1 from the
		// height when setting the window.

		u.h = msg.Height - 1
		u.resize(u.h)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return u, tea.Quit
		case "up":
			u.navigate(navigateUp)
		case "down":
			u.navigate(navigateDown)
		case "left":
			u.navigate(navigateLeft)
		case "right":
			u.navigate(navigateRight)
		case u.keyCfg.HandleLineKey:
			return u, u.handleLine
		case u.keyCfg.HandleHunkKey:
			return u, u.handleHunk
		case u.keyCfg.ResetLineKey:
			if u.keyCfg.CanReset {
				return u, u.handleResetLine
			}
		case u.keyCfg.ResetHunkKey:
			if u.keyCfg.CanReset {
				return u, u.handleResetHunk
			}
		}
	case RefreshMsg:
		return u, u.UpdateDoc
	case docMsg:
		logging.Info(`received docMsg`, `docType`, u.docType)
		u.doc = msg.d
		u.setDoc(msg.d)
	case error:
		logging.Error(msg.Error())
	}
	return u, nil
}

func (u *UI) resize(size int) {
	u.h = size
	if u.window != nil {
		u.window.Resize(size)
	} else {
		u.window = window.NewWindow(u.doc.Lines, u.h)
	}
}

func (u *UI) setDoc(doc patch.Document) {
	u.doc = doc
	if u.window == nil {
		u.window = window.NewWindow(doc.Lines, u.h)
	} else {
		u.window.SetData(doc.Lines)
		u.window.Resize(u.h)
	}
}

func (u *UI) UpdateDoc() tea.Msg {
	var doc patch.Document
	var err error

	doc, err = u.docGetter.GetDocument()
	if err != nil {
		return err
	}

	return docMsg{d: doc}
}

var kindToColor = map[patch.LineKind]lipgloss.Style{
	patch.AdditionLine: globalstyles.AdditionColor,
	patch.RemovalLine:  globalstyles.RemovalColor,
	patch.DiffLine:     lipgloss.NewStyle().Foreground(lipgloss.Color(`#FFFFFF`)),
	patch.HunkLine:     lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FFFF`)),
}

func (dv *UI) View() string {
	if dv == nil || dv.window == nil {
		return ``
	}

	sb := &strings.Builder{}

	viewableLines := dv.window.CurrentValues()

	for i, l := range viewableLines.Values {
		s := lipgloss.NewStyle()
		c, ok := kindToColor[l.Kind]
		if ok {
			s = s.Inherit(c)
		}
		if dv.cursor == i {
			s = s.Inherit(globalstyles.SelectedBackground)
		}

		sb.WriteString(s.Render(l.Text))
		sb.WriteString(l.LineBreak)
	}
	return sb.String()
}

type navigationDirection int

const (
	navigateUp navigationDirection = iota
	navigateDown
	navigateLeft
	navigateRight
)

func (dv *UI) navigate(dir navigationDirection) {
	if dv == nil {
		return
	}

	switch dir {
	case navigateUp:
		if dv.cursor == 0 {
			dv.window.ScrollUp()
		} else {
			dv.cursor--
		}
	case navigateDown:
		if dv.cursor == dv.window.Size()-1 {
			dv.window.ScrollDown()
		} else {
			dv.cursor++
		}
	case navigateLeft:
		dv.cursorLeft()
	case navigateRight:
		dv.cursorRight()
	}
}

func (dv *UI) cursorLeft() {
	start := dv.window.AbsoluteIndex(dv.cursor)
	if start <= 0 {
		return
	}

	for i := start - 1; i >= 0; i-- {
		l := dv.doc.Lines[i]
		if l.Kind == patch.HunkLine {
			dv.jumpToLine(i)
			return
		}
	}
	dv.jumpToLine(0)
}

func (dv *UI) cursorRight() {
	start := dv.window.AbsoluteIndex(dv.cursor)
	if start >= len(dv.doc.Lines)-1 {
		return
	}
	for i := start + 1; i < len(dv.doc.Lines); i++ {
		l := dv.doc.Lines[i]
		if l.Kind == patch.HunkLine {
			dv.jumpToLine(i)
			return
		}
	}
	dv.jumpToLine(len(dv.doc.Lines) - 1)
}

func (dv *UI) jumpToLine(index int) {
	relIndex := dv.window.RelativeIndex(index)
	if relIndex < 0 {
		dv.window.JumpTo(index)
		relIndex = dv.window.RelativeIndex(index)
	}
	dv.cursor = relIndex
}

func (dv *UI) currentLine() patch.Line {
	return dv.doc.Lines[dv.currentLineIndex()]
}

func (dv *UI) currentLineIndex() int {
	return dv.window.AbsoluteIndex(dv.cursor)
}

func (dv *UI) linesInCurrentHunk() []int {
	lineIdx := dv.currentLineIndex()
	h, ok := findHunk(dv.doc, lineIdx)
	if !ok {
		logging.Warn(`hunk not found`, `index`, lineIdx)
		return nil
	}

	var lines []int
	for l := h.LineStart(); l < h.LineEnd(); l++ {
		dl := dv.doc.Lines[l]
		if dl.Kind.IsAdditionOrRemoval() {
			lines = append(lines, l)
		}
	}

	return lines
}

func findHunk(doc patch.Document, idx int) (patch.Hunk, bool) {
	e, ok := doc.FindEntry(idx)
	if !ok {
		return patch.Hunk{}, false
	}

	return e.FindHunk(idx)
}

func (u *UI) handleLine() tea.Msg {
	if ln := u.currentLine(); !ln.Kind.IsAdditionOrRemoval() {
		return nil
	}

	msg := PatchMsg{
		Doc:   u.doc,
		Lines: []int{u.currentLineIndex()},
	}
	msg.Direction = patch.Stage
	if u.docType == Staged {
		msg.Direction = patch.Unstage
	}
	return msg
}

func (u *UI) handleHunk() tea.Msg {
	e, ok := u.doc.FindEntry(u.currentLineIndex())
	if !ok {
		return nil
	}
	_, ok = e.FindHunk(u.currentLineIndex())
	if !ok {
		return nil
	}

	msg := PatchMsg{
		Doc:   u.doc,
		Lines: u.linesInCurrentHunk(),
	}
	msg.Direction = patch.Stage
	if u.docType == Staged {
		msg.Direction = patch.Unstage
	}
	return msg
}

func (u *UI) handleResetLine() tea.Msg {
	return ResetMsg{
		Doc:   u.doc,
		Lines: []int{u.currentLineIndex()},
	}
}

func (u *UI) handleResetHunk() tea.Msg {
	return ResetMsg{
		Doc:   u.doc,
		Lines: u.linesInCurrentHunk(),
	}
}
