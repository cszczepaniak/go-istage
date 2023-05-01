package git

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveGitPath(t *testing.T) {
	r := NewTestRepo(t)
	err := os.Chdir(r.path)
	require.NoError(t, err)

	p, err := resolveGitPath()
	require.NoError(t, err)
	assert.NotEmpty(t, p)
	assert.Contains(t, p, `git`)
}

func TestResolveRepoPath(t *testing.T) {
	r := NewTestRepo(t)
	err := os.Chdir(r.path)
	require.NoError(t, err)

	p, err := resolveRepoPath()
	require.NoError(t, err)
	assert.NotEmpty(t, p)
	assert.Equal(t, path.Join(r.path, `.git`), path.Join(p))
}
