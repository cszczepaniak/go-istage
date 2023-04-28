package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cszczepaniak/go-istage/patch"
	"github.com/cszczepaniak/go-istage/services"
)

type lines []int

func (l *lines) UnmarshalText(bs []byte) error {
	numStrs := strings.Split(string(bs), `,`)

	ls := make(lines, 0, len(numStrs))
	for _, s := range numStrs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ls = append(ls, n)
	}

	*l = ls
	return nil
}

func (l *lines) MarshalText() ([]byte, error) {
	if l == nil {
		return nil, nil
	}
	res := &strings.Builder{}
	for i, n := range *l {
		res.WriteString(strconv.Itoa(n))
		if i < len(*l)-1 {
			res.WriteString(`,`)
		}
	}
	return []byte(res.String()), nil
}

var (
	applyPatch = flag.Bool(`apply`, false, `apply`)
)

func main() {
	linesArg := make(lines, 0)

	flag.TextVar(&linesArg, `lines`, &lines{}, ``)

	flag.Parse()

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

	if *applyPatch {
		err = ps.ApplyPatch(patch.Stage, false, linesArg)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		for _, e := range ds.Document.Entries {
			start := e.Offset
			end := e.Offset + e.Length
			for i, l := range ds.Document.Lines[start:end] {
				fmt.Print(start+i, `:`, l)
			}
		}
	}
}
