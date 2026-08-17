package main

import (
	"fmt"
	"os"

	"github.com/drobilica/tarlink-registry/internal/indexgen"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(os.Args) == 2 && os.Args[1] == "--check" {
		err = indexgen.Check(root)
	} else if len(os.Args) == 1 {
		err = indexgen.Write(root)
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/generate-index [--check]")
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
