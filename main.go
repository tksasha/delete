package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	workers, root := prepare()

	confirm(root)

	queue, wg := run(workers)

	process(root, queue)

	close(queue)

	wg.Wait()
}

func prepare() (int, string) {
	log.SetFlags(0)

	workers := flag.Int("workers", runtime.GOMAXPROCS(0), "numbers of workers")

	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("use: delete [-workers=n] directory")
	}

	root := flag.Args()[0]

	return *workers, root
}

func confirm(path string) {
	log.Printf("Are you sure to delete %s? (y/n)", path)

	var response string

	if _, err := fmt.Scanln(&response); err != nil {
		log.Fatal(err)
	}

	if strings.ToLower(response) != "y" {
		os.Exit(1)
	}
}

func run(workers int) (chan string, *sync.WaitGroup) {
	queue := make(chan string, workers)

	var wg sync.WaitGroup

	for id := range workers {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			worker(id, queue)
		}(id)
	}

	return queue, &wg
}

func process(root string, queue chan<- string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())

		if entry.IsDir() {
			dir, err := os.Open(path) //nolint:gosec
			if err != nil {
				continue
			}

			dirs, err := dir.Readdirnames(1)
			if err != nil && !errors.Is(err, io.EOF) {
				log.Printf("failed to read dir names in directory %s: %v", dir.Name(), err)

				continue
			}

			if err := dir.Close(); err != nil {
				log.Printf("failed to close directory %s: %v", dir.Name(), err)
			}

			if len(dirs) == 0 {
				queue <- path

				continue // dir is empty, go to the next entry
			}

			process(path, queue)

			continue // the entry is dir, skip it
		}

		queue <- path
	}
}

func remove(path string) (string, error) {
	if err := os.Remove(path); err != nil {
		return ".", err
	}

	return filepath.Dir(path), nil
}

func worker(id int, queue <-chan string) {
	for path := range queue {
		for {
			start := time.Now()

			var err error

			path, err = remove(path)
			if err != nil {
				break
			}

			log.Printf("%d -> %s %v\n", id, path, time.Since(start))

			if path == "." {
				log.Printf("%d -> done\n", id)

				break
			}
		}
	}
}
