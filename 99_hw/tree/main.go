package main

import (
	"fmt"
	"io"
	"os"
	// "path/filepath"
	// "strings"
)

func dirTree(out io.Writer, path string, printFiles  bool) error {
	if printFiles {
		x := os.DirFS(path)
		fmt.Printf("x: %T, %v\n", x, x)
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
