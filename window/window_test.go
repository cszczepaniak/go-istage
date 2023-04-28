package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindowScrolling(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	w := NewWindow(data, 2)

	vals := w.CurrentValues()
	assert.Equal(t, 0, vals.StartIndex)
	assert.Equal(t, []int{1, 2}, vals.Values)

	w.ScrollDown()
	w.ScrollDown()

	vals = w.CurrentValues()
	assert.Equal(t, 2, vals.StartIndex)
	assert.Equal(t, []int{3, 4}, vals.Values)

	w.ScrollDown()
	w.ScrollDown()
	w.ScrollDown()
	w.ScrollDown()

	vals = w.CurrentValues()
	assert.Equal(t, 3, vals.StartIndex)
	assert.Equal(t, []int{4, 5}, vals.Values)

	w.ScrollUp()

	vals = w.CurrentValues()
	assert.Equal(t, 2, vals.StartIndex)
	assert.Equal(t, []int{3, 4}, vals.Values)

	w.ScrollUp()

	vals = w.CurrentValues()
	assert.Equal(t, 1, vals.StartIndex)
	assert.Equal(t, []int{2, 3}, vals.Values)

	w.ScrollUp()
	w.ScrollUp()
	w.ScrollUp()

	vals = w.CurrentValues()
	assert.Equal(t, 0, vals.StartIndex)
	assert.Equal(t, []int{1, 2}, vals.Values)
}

func TestWindowResize(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	w := NewWindow(data, 2)

	vals := w.CurrentValues()
	assert.Equal(t, 0, vals.StartIndex)
	assert.Equal(t, []int{1, 2}, vals.Values)

	w.Resize(4)

	vals = w.CurrentValues()
	assert.Equal(t, 0, vals.StartIndex)
	assert.Equal(t, []int{1, 2, 3, 4}, vals.Values)

	w.Resize(1)

	vals = w.CurrentValues()
	assert.Equal(t, 0, vals.StartIndex)
	assert.Equal(t, []int{1}, vals.Values)

	w.ScrollDown()
	w.ScrollDown()

	vals = w.CurrentValues()
	assert.Equal(t, 2, vals.StartIndex)
	assert.Equal(t, []int{3}, vals.Values)

	w.Resize(2)

	vals = w.CurrentValues()
	assert.Equal(t, 2, vals.StartIndex)
	assert.Equal(t, []int{3, 4}, vals.Values)

	w.Resize(1000)

	vals = w.CurrentValues()
	assert.Equal(t, 0, vals.StartIndex)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, vals.Values)

	assert.Panics(t, func() {
		w.Resize(-1)
	})
}

func TestJumpTo(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	w := NewWindow(data, 4)

	vals := w.CurrentValues()
	assert.Equal(t, 0, vals.StartIndex)
	assert.Equal(t, []int{1, 2, 3, 4}, vals.Values)

	w.JumpTo(3)

	vals = w.CurrentValues()
	assert.Equal(t, 3, vals.StartIndex)
	assert.Equal(t, []int{4, 5, 6, 7}, vals.Values)

	w.JumpTo(6)

	vals = w.CurrentValues()
	assert.Equal(t, 6, vals.StartIndex)
	assert.Equal(t, []int{7, 8, 9, 10}, vals.Values)

	w.JumpTo(20)

	vals = w.CurrentValues()
	assert.Equal(t, 6, vals.StartIndex)
	assert.Equal(t, []int{7, 8, 9, 10}, vals.Values)

	w.JumpTo(1)

	vals = w.CurrentValues()
	assert.Equal(t, 1, vals.StartIndex)
	assert.Equal(t, []int{2, 3, 4, 5}, vals.Values)
}
