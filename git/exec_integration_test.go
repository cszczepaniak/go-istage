package git

import (
	"testing"

	"github.com/cszczepaniak/go-istage/nolibgit"
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
