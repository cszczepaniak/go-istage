package main

import (
	"fmt"
	"log"

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

	for _, e := range ds.Document.Entries {
		start := e.Offset
		end := e.Offset + e.Length
		fmt.Println(`============== ANOTHER ENTRY =================`)
		for _, l := range ds.Document.Lines[start:end] {
			fmt.Print(l)
		}
	}
}
