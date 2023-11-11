package git

import (
	"strings"
	"testing"

	"github.com/cszczepaniak/go-istage/nolibgit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	NewTestRepo(t)

	env, err := nolibgit.LoadEnvironment()
	require.NoError(t, err)

	gs, err := NewClient(env)
	require.NoError(t, err)

	err = gs.Exec(`status`).Run()
	require.NoError(t, err)
}

func TestExecWithStdout(t *testing.T) {
	NewTestRepo(t)

	env, err := nolibgit.LoadEnvironment()
	require.NoError(t, err)

	gs, err := NewClient(env)
	require.NoError(t, err)

	sb := &strings.Builder{}

	err = gs.Exec(`status`).WithStdout(sb).Run()
	require.NoError(t, err)

	assert.Equal(t, "On branch master\nnothing to commit, working tree clean\n", sb.String())
}

func TestExecWithStderr(t *testing.T) {
	NewTestRepo(t)

	env, err := nolibgit.LoadEnvironment()
	require.NoError(t, err)

	gs, err := NewClient(env)
	require.NoError(t, err)

	sb := &strings.Builder{}

	err = gs.Exec(`what`).WithStderr(sb).Run()
	assert.Error(t, err)

	assert.Equal(t, "git: 'what' is not a git command. See 'git --help'.\n\nThe most similar command is\n\tmktag\n", sb.String())
}
