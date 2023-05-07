package files

import (
	"testing"

	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/ui/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFileGetter []git.File

func (fg testFileGetter) GetFiles() ([]git.File, error) {
	return []git.File(fg), nil
}

func TestNavigation(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	files := testFileGetter{{
		Path:   `a`,
		Status: git.FileStatusAdded,
	}, {
		Path:   `b`,
		Status: git.FileStatusAdded,
	}, {
		Path:   `c`,
		Status: git.FileStatusAdded,
	}}

	fv := NewView(Unstaged, KeyConfig{}, files, 40)

	fv = testutils.InitializeModel(t, fv)

	assert.Equal(t, 0, fv.cursor)

	for _, exp := range []int{1, 2, 2, 2} {
		fv = testutils.ExecKeyPressCycle(fv, `down`)
		assert.Equal(t, exp, fv.cursor)
	}

	for _, exp := range []int{1, 0, 0, 0} {
		fv = testutils.ExecKeyPressCycle(fv, `up`)
		assert.Equal(t, exp, fv.cursor)
	}
}

func TestHandleFiles(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	files := testFileGetter{{
		Path:   `a`,
		Status: git.FileStatusAdded,
	}, {
		Path:   `b`,
		Status: git.FileStatusAdded,
	}, {
		Path:   `c`,
		Status: git.FileStatusAdded,
	}}

	fv := NewView(Unstaged, KeyConfig{
		HandleFileKey: `s`,
	}, files, 40)
	fv = testutils.InitializeModel(t, fv)

	assert.EqualValues(t, files, fv.files)

	fv, msg := testutils.ExecKeyPress(fv, `s`)
	assert.Equal(t, HandleFileMsg{
		File:      files[0],
		Direction: patch.Stage,
	}, msg)

	fv = testutils.ExecKeyPressCycle(fv, `down`)

	fv, msg = testutils.ExecKeyPress(fv, `s`)
	assert.Equal(t, HandleFileMsg{
		File:      files[1],
		Direction: patch.Stage,
	}, msg)

	fv.docType = Staged

	_, msg = testutils.ExecKeyPress(fv, `s`)
	assert.Equal(t, HandleFileMsg{
		File:      files[1],
		Direction: patch.Unstage,
	}, msg)
}

func TestUpdateFiles(t *testing.T) {
	err := logging.Init(logging.Config{})
	require.NoError(t, err)

	files := testFileGetter{{
		Path:   `a`,
		Status: git.FileStatusAdded,
	}, {
		Path:   `b`,
		Status: git.FileStatusAdded,
	}, {
		Path:   `c`,
		Status: git.FileStatusAdded,
	}}

	fv := NewView(Unstaged, KeyConfig{
		HandleFileKey: `s`,
	}, files, 40)
	fv = testutils.InitializeModel(t, fv)

	assert.EqualValues(t, files, fv.files)

	files = append(files, git.File{
		Path:   `d`,
		Status: git.FileStatusModified,
	})
	fv.fg = files

	fv = testutils.RunUpdateCycle[*UI](fv.Update(RefreshMsg{}))

	assert.EqualValues(t, files, fv.files)
}
