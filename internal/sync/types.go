package sync

import (
	"fmt"
	"kemono-sync/internal/api"
	"kemono-sync/internal/fs"
)

type Data = api.DataParams

type File struct {
	Domain string
	fs.File
}

type Post struct {
	Domain string
	api.Post
}

func (p Post) EachFile(f func(File) error) error {
	if p.File != nil {
		if err := f(File{
			Domain: p.Domain,
			File: fs.File{
				Path:    p.File.Path,
				Name:    p.File.Name,
				Service: p.Service,
				User:    p.User,
				ID:      p.ID,
				Title:   p.Title,
				Added:   p.Added.Time,
			},
		}); err != nil {
			return fmt.Errorf("file error: %w", err)
		}
	}
	for i, file := range p.Attachments {
		if err := f(File{
			Domain: p.Domain,
			File: fs.File{
				Path:    file.Path,
				Name:    file.Name,
				Service: p.Service,
				User:    p.User,
				ID:      p.ID,
				Title:   p.Title,
				Added:   p.Added.Time,
			},
		}); err != nil {
			return fmt.Errorf("attachment #%d error: %w", i, err)
		}
	}
	return nil
}
