package window

type Window[T any] struct {
	values []T
	size   int
	start  int
}

func NewWindow[T any](values []T, size int) *Window[T] {
	if size > len(values) {
		size = len(values)
	}
	return &Window[T]{
		values: values,
		size:   size,
	}
}

func (w *Window[T]) Resize(newSize int) {
	if newSize < 0 {
		panic(`Window.Resize: negative size`)
	}
	if newSize >= len(w.values) {
		newSize = len(w.values)
	}
	w.size = newSize
	if w.start > w.maxStart() {
		w.start = w.maxStart()
	}
}

func (w *Window[T]) Size() int {
	return w.size
}

func (w *Window[T]) AbsoluteIndex(windowIndex int) int {
	if windowIndex >= w.size {
		windowIndex = w.size - 1
	}
	return w.start + windowIndex
}

func (w *Window[T]) RelativeIndex(absIndex int) int {
	if !w.ContainsAbsoluteIndex(absIndex) {
		return -1
	}
	return absIndex - w.start
}

func (w *Window[T]) ScrollUp() {
	if w.start <= 0 {
		return
	}
	w.start--
}

func (w *Window[T]) ScrollDown() {
	if w.end() >= len(w.values) {
		return
	}
	w.start++
}

func (w *Window[T]) JumpTo(sourceIndex int) {
	if sourceIndex < 0 {
		panic(`Window.JumpTo: negative sourceIndex`)
	}
	if sourceIndex > w.maxStart() {
		sourceIndex = w.maxStart()
	}
	w.start = sourceIndex
}

func (w *Window[T]) ContainsAbsoluteIndex(sourceIndex int) bool {
	return sourceIndex >= w.start && sourceIndex < w.end()
}

func (w *Window[T]) end() int {
	return w.start + w.size
}

func (w *Window[T]) maxStart() int {
	return len(w.values) - w.size
}

type Values[T any] struct {
	Values     []T
	StartIndex int
}

func (w *Window[T]) CurrentValues() Values[T] {
	end := w.end()
	if end > len(w.values) {
		end = len(w.values)
	}

	return Values[T]{
		Values:     w.values[w.start:end],
		StartIndex: w.start,
	}
}
