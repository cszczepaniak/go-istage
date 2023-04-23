package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveGitPath(t *testing.T) {
	p, err := resolveGitPath()
	require.NoError(t, err)
	assert.NotEmpty(t, p)
	assert.Contains(t, p, `git`)
}

func TestResolveRepoPath(t *testing.T) {
	p, err := resolveRepoPath()
	require.NoError(t, err)
	assert.NotEmpty(t, p)
	assert.Contains(t, p, `.git`)
}
