package main

import (
	"os"
	"path/filepath"
	"strconv"
)

const (
	depth = 1
	dirs  = 1
	files = 1

	dirPerms  = 0o750 // rwxr-x---
	filePerms = 0o640 // rw-r-----

	root = "test/__fs"
)

func main() {
	if err := prepare(root, depth); err != nil {
		panic(err)
	}
}

func prepare(root string, depth int) error {
	if depth == 0 {
		return nil
	}

	for dir := range dirs {
		path := filepath.Join(root, strconv.Itoa(dir))

		if err := os.MkdirAll(path, dirPerms); err != nil {
			return err
		}

		for f := range files {
			filename := filepath.Join(
				path,
				strconv.Itoa(f)+".txt",
			)

			if err := os.WriteFile(filename, []byte("delme"), filePerms); err != nil {
				return err
			}
		}

		if err := prepare(path, depth-1); err != nil {
			return err
		}
	}

	return nil
}
