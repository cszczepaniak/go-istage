package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	NewTestRepo(t)

	env, err := NewEnvironment(``, ``)
	require.NoError(t, err)

	gs, err := NewClient(env)
	require.NoError(t, err)

	err = gs.Exec(`status`).Run()
	require.NoError(t, err)
}
