package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func dirTree(out io.Writer, path string, printFiles bool) error {

	err := printDirTree(out, path, printFiles, 0)
	return err
}

func printDirTree(out io.Writer, path string, printFiles bool, level int) error {
	curDir, err := os.Getwd()
	fp := filepath.Join(curDir, path)

	file, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	list, err := file.Readdirnames(0) //0 - to read all files and folders
	list = sort.StringSlice(list)

	for i, name := range list {
		fInfo, _ := os.Stat(filepath.Join(fp, name))

		isRoot := false
		if i == len(list)-1 || (!printFiles && i == len(list)-2) {
			isRoot = true
		}

		newPath := path + string(os.PathSeparator) + name

		printTreeItem(fInfo, name, printFiles, level, isRoot)

		level++
		printDirTree(out, newPath, printFiles, level)
		level--

	}

	return nil
}

func printTreeItem(fInfo os.FileInfo, name string, printFiles bool, level int, isRoot bool) {

	prefix := ""
	if level > 0 {
		prefix = strings.Repeat("\t|", level-1) + "\t"
	}
	if isRoot {
		prefix += "└"
	} else {
		prefix += "├"
	}

	if isRoot {
		prefix += "└"
	} else {
		prefix += "├"
	}

	if fInfo.IsDir() {
		fmt.Printf("%v───%v\n", prefix, name)

	} else if printFiles {
		fSize := "empty"
		if fInfo.Size() > 0 {
			fSize = fmt.Sprint(fInfo.Size()) + "b"
		}
		fmt.Printf("%v───%v (%v)\n", prefix, name, fSize)
	}
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
