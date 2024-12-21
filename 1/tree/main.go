package main

import (
	"fmt"
	"io"
	// "io/ioutil"
	"os"
	// "path/filepath"
	// "strings"
)


// func printDir (layer int, path string) error {
// 	fmt.Println("├")
// 	fmt.Println("─")
// 	fmt.Println("│")
	
// }

func dirTree(out io.Writer, path string, printFiles bool) error {
	// retrieve files here
	dirsAndFiles, err := os.ReadDir(path) // alrady sorted
	if err != nil {
		panic(err.Error())
	}

	for _, subPath := range dirsAndFiles{
		if subPath.IsDir(){
			fmt.Println(subPath)
			
		} else if printFiles {
			info, _ := subPath.Info()
			fmt.Println( subPath, info.Size()  )
		}
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
