package ui

import "github.com/cszczepaniak/go-istage/patch"

type refreshMsg struct{}

type docMsg struct {
	d patch.Document
}

type cursorMsg struct {
	cursor int
}
