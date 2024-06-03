package main

import (
	"flag"
	"fmt"
	"kemono-sync/internal/api"
	"kemono-sync/internal/fs"
	"kemono-sync/internal/sync"
	"log"
	"os"
	"path/filepath"
	"strings"
	stdsync "sync"
)

func PrintUsage() {
	if _, err := fmt.Fprintf(os.Stderr, "Usage: %s [options] <domain> <service> <user> <id>\n", filepath.Base(os.Args[0])); err != nil {
		panic(err)
	}
}

func main() {
	flag.Usage = PrintUsage
	flag.Parse()

	args := flag.Args()
	if len(args) != 4 {
		PrintUsage()
		os.Exit(1)
	}

	client := api.New()

	basePost, err := client.GetPost(api.PostParams{
		Domain:  args[0],
		Service: args[1],
		User:    args[2],
		ID:      args[3],
	})
	if err != nil {
		log.Fatalf("error getting post: %v", err)
	}
	post := sync.Post{
		Domain: args[0],
		Post:   *basePost,
	}

	// TODO: Conflate with get many posts impl.

	cache := &fs.Cache{BasePath: ".cache"}

	library := &fs.Library{
		BasePath: "",
		Cache:    cache,
	}

	fetcher := &sync.Fetcher{
		Client: client,
		Cache:  cache,
	}
	fetcher.Start(8)

	linker := struct {
		wg   stdsync.WaitGroup
		todo chan sync.File
	}{
		todo: make(chan sync.File, 8),
	}
	go func() {
		for file := range linker.todo {
			// TODO: Detect if is not linked - log
			if err := library.Link(file.File); err != nil {
				log.Fatalf("error linking file %s%s: %v", file.Domain, file.Path, err)
			}
			linker.wg.Done()
		}
	}()

	if err := post.EachFile(func(file sync.File) error {
		done := func() {
			linker.wg.Add(1)
			linker.todo <- file
		}

		if cache.Has(strings.TrimLeft(file.Path, "/")) {
			done()
			return nil
		}
		log.Printf("%s%s discovered", file.Domain, file.Path)

		fetcher.Fetch(sync.FetchParams{
			Data: sync.Data{
				Domain: file.Domain,
				Path:   file.Path,
			},
			Done: func(err error) {
				if err != nil {
					log.Printf("error fetching file %s%s: %v\n", file.Domain, file.Path, err)
				}
				done()
			},
		})
		return nil
	}); err != nil {
		log.Printf("error processing files: %v\n", err)
	}

	fetcher.Close()

	linker.wg.Wait()
}
