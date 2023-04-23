package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	env, err := NewGitEnvironment(``, ``)
	require.NoError(t, err)

	gs, err := NewGitService(env)
	require.NoError(t, err)

	err = gs.Exec(`status`).Run()
	require.NoError(t, err)
}
