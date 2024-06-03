package sync

import (
	"kemono-sync/internal/api"
	"kemono-sync/internal/fs"
	"log"
	"strings"
	"sync"
)

type FetchParams struct {
	Data Data
	Done func(err error)
}

type Fetcher struct {
	Client *api.Client
	Cache  *fs.Cache
	todo   chan FetchParams
	wg     sync.WaitGroup
}

func (f *Fetcher) fetch(params FetchParams) {
	body, err := f.Client.GetData(params.Data)
	if err != nil {
		params.Done(err)
		return
	}
	if err := f.Cache.Put(strings.TrimLeft(params.Data.Path, "/"), body); err != nil {
		params.Done(err)
		return
	}
	log.Printf("%s%s cached\n", params.Data.Domain, params.Data.Path)
	params.Done(nil)
}

func (f *Fetcher) Start(workers int) {
	if f.todo == nil {
		f.todo = make(chan FetchParams)
	}
	for i := 0; i < workers; i++ {
		go func(i int) {
			for params := range f.todo {
				f.fetch(params)
				f.wg.Done()
			}
		}(i)
	}
}

func (f *Fetcher) Fetch(params FetchParams) {
	f.wg.Add(1)
	f.todo <- params
}

func (f *Fetcher) FetchSync(data Data) error {
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	f.Fetch(FetchParams{
		Data: data,
		Done: func(_err error) {
			err = _err
			wg.Done()
		},
	})
	wg.Wait()
	return err
}

func (f *Fetcher) Close() {
	close(f.todo)
	f.wg.Wait()
}
