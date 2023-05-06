package ui

import (
	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/patch"
)

type refreshMsg struct{}

type docMsg struct {
	staged bool
	d      patch.Document
}

type filesMsg struct {
	files []git.File
}

type goToStateMsg struct {
	state StateVariant
}
