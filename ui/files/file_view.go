package files

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/git"
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

type fileGetter interface {
	GetFiles() ([]git.File, error)
}

type KeyConfig struct {
	HandleFileKey string
}

type UI struct {
	w *window.Window[git.File]

	files []git.File
	fg    fileGetter

	docType DocType

	keyCfg KeyConfig

	h      int
	cursor int
}

func NewView(docType DocType, keyCfg KeyConfig, fg fileGetter, windowSize int) *UI {
	return &UI{
		w:       window.NewWindow[git.File](nil, windowSize),
		files:   nil,
		docType: docType,
		keyCfg:  keyCfg,
		fg:      fg,
		h:       windowSize,
		cursor:  0,
	}
}

func (u *UI) Init() tea.Cmd {
	return u.UpdateFiles
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
		case u.keyCfg.HandleFileKey:
			return u, u.handleFile
		}
	case RefreshMsg:
		return u, u.UpdateFiles
	case filesMsg:
		logging.Info(`filesMsg`, `len(files)`, len(msg.files))
		u.setFiles(msg.files, u.h)
	case error:
		logging.Error(msg.Error())
	}
	return u, nil
}

type navigationDirection int

const (
	navigateUp navigationDirection = iota
	navigateDown
	navigateLeft
	navigateRight
)

func (u *UI) navigate(dir navigationDirection) {
	if u == nil {
		return
	}

	switch dir {
	case navigateUp:
		if u.cursor == 0 {
			u.w.ScrollUp()
		} else {
			u.cursor--
		}
	case navigateDown:
		if u.cursor == u.w.Size()-1 {
			u.w.ScrollDown()
		} else {
			u.cursor++
		}
	}
}

func (dv *UI) resize(size int) {
	dv.w.Resize(size)
}

func (dv *UI) setFiles(files []git.File, h int) {
	dv.files = files
	if dv.w == nil {
		dv.w = window.NewWindow(files, h)
	} else {
		dv.w.SetData(files)
		dv.w.Resize(len(files))
	}
}

var fileStatusToColor = map[git.FileStatus]lipgloss.Style{
	git.FileStatusAdded:     globalstyles.AdditionColor,
	git.FileStatusDeleted:   globalstyles.RemovalColor,
	git.FileStatusUntracked: globalstyles.AdditionColor,
	git.FileStatusModified:  lipgloss.NewStyle().Foreground(lipgloss.Color(`#DAA520`)),
	git.FileStatusRenamed:   lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FFFF`)),
}

func (dv *UI) View() string {
	if dv == nil || dv.w == nil {
		return ``
	}

	sb := &strings.Builder{}

	viewableLines := dv.w.CurrentValues()

	for i, l := range viewableLines.Values {
		s := lipgloss.NewStyle()
		c, ok := fileStatusToColor[l.Status]
		if ok {
			s = s.Inherit(c)
		}
		if dv.cursor == i {
			s = s.Inherit(globalstyles.SelectedBackground)
		}

		fmt.Fprintln(sb, s.Render(l.Path))
	}
	return sb.String()
}

func (u *UI) currentFile() git.File {
	return u.files[u.cursor]
}

func (u *UI) handleFile() tea.Msg {
	msg := HandleFileMsg{
		File: u.currentFile(),
	}
	msg.Direction = patch.Stage
	if u.docType == Staged {
		msg.Direction = patch.Unstage
	}
	return msg
}

func (u *UI) UpdateFiles() tea.Msg {
	var files []git.File
	var err error

	files, err = u.fg.GetFiles()
	if err != nil {
		return err
	}

	return filesMsg{files: files}
}
