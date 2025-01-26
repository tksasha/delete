package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("use delete PATH")
	}

	path := os.Args[1]

	var wg sync.WaitGroup

	files := make(chan string)

	for id := range 100 {
		wg.Add(1)

		go func(id int, files chan string) {
			defer wg.Done()

			for file := range files {
				log.Printf("id: %d, file: %s", id, file)

				if err := os.Remove(file); err != nil {
					log.Fatalf("failed to remove file %s, error: %v", file, err)
				}
			}
		}(id, files)
	}

	for id := range 100 {
		filename := filepath.Join(path, strconv.Itoa(id))

		if _, err := os.Create(filename); err != nil {
			log.Fatalf("failed to create file, error: %v", err)
		}
	}

	walk(path, files)

	wg.Wait()
}

func walk(path string, files chan string) {
	defer func() {
		close(files)
	}()

	callback := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("error: %v", err)

			return err
		}

		if !info.IsDir() {
			files <- path
		}

		return nil
	}

	if err := filepath.Walk(path, callback); err != nil {
		log.Fatalf("failed to walk, error: %v", err)
	}
}
