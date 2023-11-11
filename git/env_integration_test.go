package git

import (
	"path"
	"testing"

	"github.com/cszczepaniak/go-istage/nolibgit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveGitPath(t *testing.T) {
	NewTestRepo(t)

	env, err := nolibgit.LoadEnvironment()
	require.NoError(t, err)
	assert.Contains(t, env.GitExecutable, `git`)
}

func TestResolveRepoPath(t *testing.T) {
	r := NewTestRepo(t)

	env, err := nolibgit.LoadEnvironment()
	require.NoError(t, err)

	assert.Equal(t, path.Join(r.path, `.git`), path.Join(env.RepoDir))
}
