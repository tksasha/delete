package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func main() {
	workers := flag.Int("workers", runtime.GOMAXPROCS(0), "number of workers")

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("use: delete [-workers=N] PATH")
	}

	path := flag.Args()[0]

	var wg sync.WaitGroup

	files := make(chan string)

	for worker := range *workers {
		wg.Add(1)

		go func(worker int, files chan string) {
			defer wg.Done()

			delete(worker, files)
		}(worker, files)
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

func delete(worker int, files chan string) {
	for file := range files {
		start := time.Now()

		log.Printf("%d -> %s %v", worker, file, time.Since(start))

		if err := os.Remove(file); err != nil {
			log.Fatalf("failed to remove file %s, error: %v", file, err)
		}
	}
}
