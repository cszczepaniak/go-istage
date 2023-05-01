package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/window"
)

type fileView struct {
	w      *window.Window[git.File]
	cursor int
}

func newFileView(files []git.File, windowSize int) *fileView {
	return &fileView{
		w: window.NewWindow(files, windowSize),
	}
}

func (dv *fileView) resize(size int) {
	dv.w.Resize(size)
}

func (dv *fileView) setFiles(files []git.File, h int) {
	if dv.w == nil {
		dv.w = window.NewWindow(files, h)
	} else {
		dv.w.SetData(files)
	}
}

var fileStatusToColor = map[git.FileStatus]lipgloss.Style{
	git.FileStatusAdded:     lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FF00`)),
	git.FileStatusUntracked: lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FF00`)),
	git.FileStatusDeleted:   lipgloss.NewStyle().Foreground(lipgloss.Color(`#FF0000`)),
	git.FileStatusModified:  lipgloss.NewStyle().Foreground(lipgloss.Color(`#DAA520`)),
	git.FileStatusRenamed:   lipgloss.NewStyle().Foreground(lipgloss.Color(`#00FFFF`)),
}

func (dv *fileView) view() string {
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
			s = s.Inherit(selectedStyle)
		}

		fmt.Fprintln(sb, s.Render(l.Path))
	}
	return sb.String()
}

func (dv *fileView) navigate(dir navigationDirection) {
	if dv == nil {
		return
	}

	switch dir {
	case navigateUp:
		if dv.cursor == 0 {
			dv.w.ScrollUp()
		} else {
			dv.cursor--
		}
	case navigateDown:
		if dv.cursor == dv.w.Size()-1 {
			dv.w.ScrollDown()
		} else {
			dv.cursor++
		}
	}
}
