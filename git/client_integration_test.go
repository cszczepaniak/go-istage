package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnstagedChanges(t *testing.T) {
	r := NewTestRepo(t)

	f := r.MakeFile(t, `a.txt`).AddLine(`abc`).Build()

	gc, err := NewClient(r.env)
	require.NoError(t, err)

	c, err := gc.UnstagedChanges()
	require.NoError(t, err)
	require.Len(t, c, 1)

	assert.Equal(t, `diff --git a/a.txt b/a.txt
new file mode 100644
index 0000000..8baef1b
--- /dev/null
+++ b/a.txt
@@ -0,0 +1 @@
+abc
`, c[0])

	r.AddAll()
	r.Commit(`adding file`)

	f.Append("abc\n")

	c, err = gc.UnstagedChanges()
	require.NoError(t, err)
	require.Len(t, c, 1)

	assert.Equal(t, `diff --git a/a.txt b/a.txt
index 8baef1b..5d8a556 100644
--- a/a.txt
+++ b/a.txt
@@ -1 +1,2 @@
 abc
+abc
`, c[0])

	r.ClearUnstagedChanges()

	f.Replace("def")

	c, err = gc.UnstagedChanges()
	require.NoError(t, err)
	require.Len(t, c, 1)

	assert.Equal(t, `diff --git a/a.txt b/a.txt
index 8baef1b..0c00383 100644
--- a/a.txt
+++ b/a.txt
@@ -1 +1 @@
-abc
+def
\ No newline at end of file
`, c[0])

	r.ClearUnstagedChanges()

	f.Remove()

	c, err = gc.UnstagedChanges()
	require.NoError(t, err)
	require.Len(t, c, 1)

	assert.Equal(t, `diff --git a/a.txt b/a.txt
deleted file mode 100644
index 8baef1b..0000000
--- a/a.txt
+++ /dev/null
@@ -1 +0,0 @@
-abc
`, c[0])
}

func TestStagedChanges(t *testing.T) {
	r := NewTestRepo(t)

	f := r.MakeFile(t, `b.txt`).AddLine(`def`).ShouldStage().Build()
	_ = f

	gc, err := NewClient(r.env)
	require.NoError(t, err)

	c, err := gc.StagedChanges()
	require.NoError(t, err)
	require.Len(t, c, 1)

	assert.Equal(t, `diff --git a/b.txt b/b.txt
new file mode 100644
index 0000000..24c5735
--- /dev/null
+++ b/b.txt
@@ -0,0 +1 @@
+def
`, c[0])

	r.Commit(`added b.txt`)

	f.Append("abc\n")
	r.AddAll()

	c, err = gc.StagedChanges()
	require.NoError(t, err)
	require.Len(t, c, 1)

	assert.Equal(t, `diff --git a/b.txt b/b.txt
index 24c5735..7320698 100644
--- a/b.txt
+++ b/b.txt
@@ -1 +1,2 @@
 def
+abc
`, c[0])

	r.Commit(`change`)
	f.Remove()
	r.AddAll()

	c, err = gc.StagedChanges()
	require.NoError(t, err)
	require.Len(t, c, 1)

	assert.Equal(t, `diff --git a/b.txt b/b.txt
deleted file mode 100644
index 7320698..0000000
--- a/b.txt
+++ /dev/null
@@ -1,2 +0,0 @@
-def
-abc
`, c[0])
}
