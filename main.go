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
	if dirContentLen > 0 {
		for _, entry := range dirContent[:dirContentLen-1] {
			entryName := entry.Name()

			if entry.IsDir() {
				fmt.Fprintf(out, "%v%v\n", prefix+"├───", entryName)

				nextPath := path2.Join(path, entryName)
				if err := recOutput(out, prefix+"│	", nextPath, printFiles); err != nil {
					return err
				}
			} else {
				fmt.Fprintf(out, "%v%v (%v)\n", prefix+"├───", entryName, formatSize(entry.Size()))
			}
		}

		lastEntry := dirContent[dirContentLen-1]
		lastEntryName := lastEntry.Name()

		if lastEntry.IsDir() {
			fmt.Fprintf(out, "%v%v\n", prefix+"└───", lastEntryName)

			nextPath := path2.Join(path, lastEntryName)
			if err := recOutput(out, prefix+"	", nextPath, printFiles); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(out, "%v%v (%v)\n", prefix+"└───", lastEntryName, formatSize(lastEntry.Size()))
		}
	}
	return nil
}

func readDir(path string, withFiles bool) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
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

func formatSize(size int64) string {
	if size == 0 {
		return "empty"
	}
	return fmt.Sprintf("%vb", size)
}
