package main

import (
	"fmt"
	"log"

	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/services"
)

func main() {
	gitEnv, err := services.NewGitEnvironment(`/home/connor/src/go-istage`, ``)
	if err != nil {
		log.Fatalln(err)
	}

	gs, err := services.NewGitService(gitEnv)
	if err != nil {
		log.Fatalln(err)
	}

	ds, err := services.NewDocumentService(gs)
	if err != nil {
		log.Fatalln(err)
	}

	ps := services.NewPatchingService(gs, ds)

	for _, e := range ds.Document.Entries {
		start := e.Offset
		end := e.Offset + e.Length
		for i, l := range ds.Document.Lines[start:end] {
			fmt.Print(start+i, `:`, l)
		}
	}

	err = ps.ApplyPatch(patch.Stage, false, 18)
	if err != nil {
		log.Fatalln(err)
	}
}
