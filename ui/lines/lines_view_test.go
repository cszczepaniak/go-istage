package lines

import (
	"testing"

	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/ui/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDocGetter patch.Document

func (d testDocGetter) GetDocument() (patch.Document, error) {
	return patch.Document(d), nil
}

func TestNavigationUpAndDown(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	doc := patch.ParseDocument([]string{
		`diff --git a/new_file.txt b/new_file.txt
new file mode 100644
index 0000000..1fa3451
--- /dev/null
+++ b/new_file.txt
@@ -0,0 +1,1 @@
+abc`})

	lv := New(Unstaged, testDocGetter(doc), Config{}, 40)
	lv = testutils.InitializeModel(t, lv)

	assert.Equal(t, 0, lv.cursor)

	for _, exp := range []int{1, 2, 3, 4, 5, 6, 6, 6} {
		lv = testutils.ExecKeyPressCycle(lv, `down`)
		assert.Equal(t, exp, lv.cursor)
	}

	for _, exp := range []int{5, 4, 3, 2, 1, 0, 0, 0} {
		lv = testutils.ExecKeyPressCycle(lv, `up`)
		assert.Equal(t, exp, lv.cursor)
	}
}

func TestNavigationLeftAndRight(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	doc := patch.ParseDocument([]string{
		`diff --git a/new_file.txt b/new_file.txt
new file mode 100644
index 0000000..1fa3451
--- /dev/null
+++ b/new_file.txt
@@ -0,0 +1,1 @@
+abc
@@ -0,0 +1,1 @@
+abc
@@ -0,0 +1,1 @@
+abc`})

	lv := New(Unstaged, testDocGetter(doc), Config{}, 40)
	lv = testutils.InitializeModel(t, lv)

	assert.Equal(t, 0, lv.cursor)

	for _, exp := range []int{5, 7, 9, 10} {
		lv = testutils.ExecKeyPressCycle(lv, `right`)
		assert.Equal(t, exp, lv.cursor)
	}

	for _, exp := range []int{9, 7, 5, 0} {
		lv = testutils.ExecKeyPressCycle(lv, `left`)
		assert.Equal(t, exp, lv.cursor)
	}
}

func TestHandlePatchLines(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	doc := patch.ParseDocument([]string{
		`diff --git a/new_file.txt b/new_file.txt
new file mode 100644
index 0000000..1fa3451
--- /dev/null
+++ b/new_file.txt
@@ -0,0 +1,1 @@
+abc
`})

	lv := New(Unstaged, testDocGetter(doc), Config{
		HandleLineKey: `s`,
		HandleHunkKey: `S`,
	}, 40)
	lv = testutils.InitializeModel(t, lv)

	lv, msg := testutils.ExecKeyPress(lv, `s`)
	assert.Nil(t, msg)

	for i := 0; i < 6; i++ {
		lv = testutils.ExecKeyPressCycle(lv, `down`)
	}

	lv, msg = testutils.ExecKeyPress(lv, `s`)

	assert.Equal(t, PatchMsg{
		Direction: patch.Stage,
		Doc:       patch.Document(doc),
		Lines:     []int{6},
	}, msg)

	lv.docType = Staged

	_, msg = testutils.ExecKeyPress(lv, `s`)

	assert.Equal(t, PatchMsg{
		Direction: patch.Unstage,
		Doc:       patch.Document(doc),
		Lines:     []int{6},
	}, msg)
}

func TestHandlePatchHunks(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	doc := patch.ParseDocument([]string{
		`diff --git a/new_file.txt b/new_file.txt
new file mode 100644
index 0000000..1fa3451
--- /dev/null
+++ b/new_file.txt
@@ -0,0 +1,1 @@
+abc
-abc
+abc
+abc
@@ -0,0 +1,1 @@
+abc
-abc
`,
	})

	lv := New(Unstaged, testDocGetter(doc), Config{
		HandleLineKey: `s`,
		HandleHunkKey: `S`,
	}, 40)
	lv = testutils.InitializeModel(t, lv)

	lv, msg := testutils.ExecKeyPress(lv, `S`)
	assert.Nil(t, msg)

	lv = testutils.ExecKeyPressCycle(lv, `right`)

	for i := 0; i < 5; i++ {
		lv, msg = testutils.ExecKeyPress(lv, `S`)
		assert.Equal(t, PatchMsg{
			Direction: patch.Stage,
			Doc:       patch.Document(doc),
			Lines:     []int{6, 7, 8, 9},
		}, msg)

		lv = testutils.ExecKeyPressCycle(lv, `down`)
	}

	for i := 0; i < 3; i++ {
		lv, msg = testutils.ExecKeyPress(lv, `S`)
		assert.Equal(t, PatchMsg{
			Direction: patch.Stage,
			Doc:       patch.Document(doc),
			Lines:     []int{11, 12},
		}, msg)

		lv = testutils.ExecKeyPressCycle(lv, `down`)
	}

	lv = testutils.ExecKeyPressCycle(lv, `left`)
	lv = testutils.ExecKeyPressCycle(lv, `left`)
	lv.docType = Staged

	for i := 0; i < 5; i++ {
		lv, msg = testutils.ExecKeyPress(lv, `S`)
		assert.Equal(t, PatchMsg{
			Direction: patch.Unstage,
			Doc:       patch.Document(doc),
			Lines:     []int{6, 7, 8, 9},
		}, msg)

		lv = testutils.ExecKeyPressCycle(lv, `down`)
	}

	for i := 0; i < 3; i++ {
		lv, msg = testutils.ExecKeyPress(lv, `S`)
		assert.Equal(t, PatchMsg{
			Direction: patch.Unstage,
			Doc:       patch.Document(doc),
			Lines:     []int{11, 12},
		}, msg)

		lv = testutils.ExecKeyPressCycle(lv, `down`)
	}
}

func TestHandleResetLines(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	doc := patch.ParseDocument([]string{
		`diff --git a/new_file.txt b/new_file.txt
new file mode 100644
index 0000000..1fa3451
--- /dev/null
+++ b/new_file.txt
@@ -0,0 +1,1 @@
+abc
`})

	lv := New(Unstaged, testDocGetter(doc), Config{
		ResetLineKey: `r`,
		ResetHunkKey: `R`,
		CanReset:     false,
	}, 40)
	lv = testutils.InitializeModel(t, lv)

	lv, msg := testutils.ExecKeyPress(lv, `r`)
	assert.Nil(t, msg)

	for i := 0; i < 6; i++ {
		lv = testutils.ExecKeyPressCycle(lv, `down`)
	}

	lv, msg = testutils.ExecKeyPress(lv, `r`)
	assert.Nil(t, msg)

	lv.keyCfg.CanReset = true

	lv, msg = testutils.ExecKeyPress(lv, `r`)
	assert.Equal(t, ResetMsg{
		Doc:   lv.doc,
		Lines: []int{6},
	}, msg)
}

func TestHandleResetHunks(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	doc := patch.ParseDocument([]string{
		`diff --git a/new_file.txt b/new_file.txt
new file mode 100644
index 0000000..1fa3451
--- /dev/null
+++ b/new_file.txt
@@ -0,0 +1,1 @@
+abc
-abc
+abc
+abc
@@ -0,0 +1,1 @@
+abc
-abc
`,
	})

	lv := New(Unstaged, testDocGetter(doc), Config{
		ResetLineKey: `r`,
		ResetHunkKey: `R`,
		CanReset:     false,
	}, 40)
	lv = testutils.InitializeModel(t, lv)

	lv, msg := testutils.ExecKeyPress(lv, `R`)
	assert.Nil(t, msg)

	lv = testutils.ExecKeyPressCycle(lv, `right`)

	lv, msg = testutils.ExecKeyPress(lv, `R`)
	assert.Nil(t, msg)

	lv.keyCfg.CanReset = true

	for i := 0; i < 5; i++ {
		lv, msg = testutils.ExecKeyPress(lv, `R`)
		assert.Equal(t, ResetMsg{
			Doc:   patch.Document(doc),
			Lines: []int{6, 7, 8, 9},
		}, msg)

		lv = testutils.ExecKeyPressCycle(lv, `down`)
	}

	for i := 0; i < 3; i++ {
		lv, msg = testutils.ExecKeyPress(lv, `R`)
		assert.Equal(t, ResetMsg{
			Doc:   patch.Document(doc),
			Lines: []int{11, 12},
		}, msg)

		lv = testutils.ExecKeyPressCycle(lv, `down`)
	}
}

func TestUpdateDocument(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	doc1 := patch.ParseDocument([]string{
		`diff --git a/new_file.txt b/new_file.txt
new file mode 100644
index 0000000..1fa3451
--- /dev/null
+++ b/new_file.txt
@@ -0,0 +1,1 @@
+abc
`})

	doc2 := patch.ParseDocument([]string{
		`diff --git a/new_file.txt b/new_file.txt
new file mode 100644
index 0000000..1fa3451
--- /dev/null
+++ b/new_file.txt
@@ -0,0 +1,1 @@
+abc
-abc
+abc
+abc
@@ -0,0 +1,1 @@
+abc
-abc
`})

	lv := New(Unstaged, testDocGetter(doc1), Config{}, 40)
	lv = testutils.InitializeModel(t, lv)

	assert.EqualValues(t, doc1, lv.doc)

	lv.docGetter = testDocGetter(doc2)

	lv = testutils.RunUpdateCycle[*UI](lv.Update(RefreshMsg{}))

	assert.EqualValues(t, doc2, lv.doc)
}
