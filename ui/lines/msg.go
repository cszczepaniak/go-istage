package lines

import "github.com/cszczepaniak/go-istage/patch"

type RefreshMsg struct{}

type docMsg struct {
	d patch.Document
}

type PatchMsg struct {
	Direction patch.Direction
	Doc       patch.Document
	Lines     []int
}
