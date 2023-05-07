package files

import (
	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/patch"
)

type RefreshMsg struct{}

type HandleFileMsg struct {
	File      git.File
	Direction patch.Direction
}

type filesMsg struct {
	files []git.File
}
