package api

import (
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (m *Time) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("\"2006-01-02T15:04:05.999999\"", string(b))
	if err != nil {
		return err
	}
	m.Time = t
	return nil
}

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (f File) Hash() string {
	path := strings.Split(f.Path, "/")
	if len(path) < 3 {
		return ""
	}
	return path[2]
}

type Post struct {
	ID          string `json:"id"`
	Service     string `json:"service"`
	User        string `json:"user"`
	Title       string `json:"title"`
	Added       Time   `json:"added"`
	File        *File  `json:"file"`
	Attachments []File `json:"attachments"`
	Next        string `json:"next"`
}
