package main

import (
	"fmt"
	"io"
	"os"
	path2 "path"
	"sort"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	if err := recOutput(out, "", path, printFiles); err != nil {
		return err
	}

	return nil
}

func recOutput(out io.Writer, prefix string, path string, printFiles bool) error {
	dirContent, err := readDir(path, printFiles)
	if err != nil {
		return err
	}
	dirContentLen := len(dirContent)
	if dirContentLen == 0 {
		lastPrefix := prefix + "└───"
		fmt.Fprintf(out, "%v%v\n", lastPrefix, path2.Dir(path))
	} else if dirContentLen > 0 {
		for _, entry := range dirContent[:dirContentLen-1] {
			entryName := entry.Name()
			if entry.IsDir() {
				nextPrefix := prefix + "├───"
				fmt.Fprintf(out, "%v%v\n", nextPrefix, entryName)
				nextPath := path2.Join(path, entryName)
				nextPrefix = prefix + "│	"
				recOutput(out, nextPrefix, nextPath, printFiles)
			} else if printFiles {
				nextPrefix := prefix + "├───"
				fmt.Fprintf(out, "%v%v\n", nextPrefix, entryName)
			}
		}

		lastEntry := dirContent[dirContentLen-1]
		lastEntryName := lastEntry.Name()

		if lastEntry.IsDir() {
			nextPrefix := prefix + "└───"
			fmt.Fprintf(out, "%v%v\n", nextPrefix, lastEntryName)
			nextPath := path2.Join(path, lastEntry.Name())
			lastPrefix := prefix + "	"
			recOutput(out, lastPrefix, nextPath, printFiles)
		} else if printFiles {
			lastPrefix := prefix + "└───"
			fmt.Fprintf(out, "%v%v\n", lastPrefix, lastEntryName)
		}
	}
	return err
}

func readDir(path string, withFiles bool) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	if !withFiles {
		dirs := make([]os.FileInfo, 0, len(list))
		for _, entry := range list {
			if entry.IsDir() {
				dirs = append(dirs, entry)
			}
		}
		list = dirs
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}
