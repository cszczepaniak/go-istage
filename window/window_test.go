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

func TestContainsSourceIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	w := NewWindow(data, 4)

	assert.True(t, w.ContainsAbsoluteIndex(0))
	assert.True(t, w.ContainsAbsoluteIndex(1))
	assert.True(t, w.ContainsAbsoluteIndex(2))
	assert.True(t, w.ContainsAbsoluteIndex(3))
	assert.False(t, w.ContainsAbsoluteIndex(4))
	assert.False(t, w.ContainsAbsoluteIndex(5))
	assert.False(t, w.ContainsAbsoluteIndex(6))
	assert.False(t, w.ContainsAbsoluteIndex(7))
	assert.False(t, w.ContainsAbsoluteIndex(8))
	assert.False(t, w.ContainsAbsoluteIndex(9))
	assert.False(t, w.ContainsAbsoluteIndex(-1))
	assert.False(t, w.ContainsAbsoluteIndex(1000))

	w.JumpTo(3)

	assert.False(t, w.ContainsAbsoluteIndex(0))
	assert.False(t, w.ContainsAbsoluteIndex(1))
	assert.False(t, w.ContainsAbsoluteIndex(2))
	assert.True(t, w.ContainsAbsoluteIndex(3))
	assert.True(t, w.ContainsAbsoluteIndex(4))
	assert.True(t, w.ContainsAbsoluteIndex(5))
	assert.True(t, w.ContainsAbsoluteIndex(6))
	assert.False(t, w.ContainsAbsoluteIndex(7))
	assert.False(t, w.ContainsAbsoluteIndex(8))
	assert.False(t, w.ContainsAbsoluteIndex(9))
	assert.False(t, w.ContainsAbsoluteIndex(-1))
	assert.False(t, w.ContainsAbsoluteIndex(1000))

	w.Resize(2)

	assert.False(t, w.ContainsAbsoluteIndex(0))
	assert.False(t, w.ContainsAbsoluteIndex(1))
	assert.False(t, w.ContainsAbsoluteIndex(2))
	assert.True(t, w.ContainsAbsoluteIndex(3))
	assert.True(t, w.ContainsAbsoluteIndex(4))
	assert.False(t, w.ContainsAbsoluteIndex(5))
	assert.False(t, w.ContainsAbsoluteIndex(6))
	assert.False(t, w.ContainsAbsoluteIndex(7))
	assert.False(t, w.ContainsAbsoluteIndex(8))
	assert.False(t, w.ContainsAbsoluteIndex(9))
	assert.False(t, w.ContainsAbsoluteIndex(-1))
	assert.False(t, w.ContainsAbsoluteIndex(1000))
}

func TestAbsoluteIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	w := NewWindow(data, 4)

	assert.Equal(t, 0, w.AbsoluteIndex(0))
	assert.Equal(t, 1, w.AbsoluteIndex(1))
	assert.Equal(t, 2, w.AbsoluteIndex(2))
	assert.Equal(t, 3, w.AbsoluteIndex(3))
	assert.Equal(t, 3, w.AbsoluteIndex(4))

	w.ScrollDown()

	assert.Equal(t, 1, w.AbsoluteIndex(0))
	assert.Equal(t, 2, w.AbsoluteIndex(1))
	assert.Equal(t, 3, w.AbsoluteIndex(2))
	assert.Equal(t, 4, w.AbsoluteIndex(3))
	assert.Equal(t, 4, w.AbsoluteIndex(4))

	w.JumpTo(10)

	assert.Equal(t, 6, w.AbsoluteIndex(0))
	assert.Equal(t, 7, w.AbsoluteIndex(1))
	assert.Equal(t, 8, w.AbsoluteIndex(2))
	assert.Equal(t, 9, w.AbsoluteIndex(3))
	assert.Equal(t, 9, w.AbsoluteIndex(4))
}

func TestRelativeIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	w := NewWindow(data, 4)

	assert.Equal(t, 0, w.RelativeIndex(0))
	assert.Equal(t, 1, w.RelativeIndex(1))
	assert.Equal(t, 2, w.RelativeIndex(2))
	assert.Equal(t, 3, w.RelativeIndex(3))
	assert.Equal(t, -1, w.RelativeIndex(4))
	assert.Equal(t, -1, w.RelativeIndex(5))
	assert.Equal(t, -1, w.RelativeIndex(6))
	assert.Equal(t, -1, w.RelativeIndex(7))
	assert.Equal(t, -1, w.RelativeIndex(8))
	assert.Equal(t, -1, w.RelativeIndex(9))

	w.ScrollDown()
	w.ScrollDown()
	w.ScrollDown()

	assert.Equal(t, -1, w.RelativeIndex(0))
	assert.Equal(t, -1, w.RelativeIndex(1))
	assert.Equal(t, -1, w.RelativeIndex(2))
	assert.Equal(t, 0, w.RelativeIndex(3))
	assert.Equal(t, 1, w.RelativeIndex(4))
	assert.Equal(t, 2, w.RelativeIndex(5))
	assert.Equal(t, 3, w.RelativeIndex(6))
	assert.Equal(t, -1, w.RelativeIndex(7))
	assert.Equal(t, -1, w.RelativeIndex(8))
	assert.Equal(t, -1, w.RelativeIndex(9))
}
