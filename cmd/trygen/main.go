package main

import (
	"log"
	"os"

	"github.com/goghcrow/go-try/rewriter"
)

func main() {
	goFile := os.Getenv("GOFILE")
	if goFile == "" {
		panic("Must run in go:generate mode")
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flags := log.Flags()
	defer log.SetFlags(flags)
	log.SetFlags(0)
	log.SetPrefix("[rewrite-try] ")

	rewriter.Rewrite(cwd)
}
