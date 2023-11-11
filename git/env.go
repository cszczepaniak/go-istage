package git

import (
	"github.com/cszczepaniak/go-istage/nolibgit"
)

type Environment struct {
	repoPath  string
	pathToGit string
}

func newEnvironmentFromNoLibGit(env nolibgit.Environment) Environment {
	return Environment{
		repoPath:  env.RepoDir,
		pathToGit: env.GitExecutable,
	}
}
