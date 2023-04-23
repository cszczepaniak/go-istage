package main

import (
	"fmt"
	"os"

	git "github.com/libgit2/git2go/v34"
)

func main() {
	repo, err := git.OpenRepository("/home/connor/src/go-istage")

	if err != nil {

		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)

	}

	cfg, err := repo.Config()
	if err != nil {
		panic(err)
	}
	s, err := cfg.LookupBool("pull.rebase")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
