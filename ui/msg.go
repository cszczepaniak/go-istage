package ui

import "github.com/cszczepaniak/go-istage/patch"

type refreshMsg struct{}

type docMsg struct {
	staged bool
	d      patch.Document
}

type jumpToDocLineIndexMsg struct {
	index int
}

type windowScrollDownMsg struct{}
type windowScrollUpMsg struct{}
