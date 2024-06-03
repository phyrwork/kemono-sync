package fs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	Path    string
	Name    string
	Service string
	User    string
	ID      string
	Title   string
	Added   time.Time
}

type Cache struct {
	BasePath string
}

func (c *Cache) Path(path string) string {
	if filepath.IsAbs(path) {
		panic(fmt.Errorf("path cannot be an absolute path: %s", path)) // Coding error.
	}
	return filepath.Join(c.BasePath, path)
}

func (c *Cache) EnsureDir(path string) error {
	dir := filepath.Join(c.BasePath, path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create directory %s: %w", dir, err)
	}
	return nil
}

func (c *Cache) Abs(path string) string {
	path, err := filepath.Abs(c.Path(path))
	if err != nil {
		panic(err)
	}
	return path
}

func (c *Cache) Has(path string) bool {
	path = c.Path(path)
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if stat.IsDir() {
		return false
	}
	if stat.Size() == 0 {
		// Assume zero-size files are a cache error.
		return false
	}
	return true
}

func (c *Cache) Put(path string, data io.ReadCloser) error {
	// Create cache directory before attempting download
	if err := c.EnsureDir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("could not ensure directory %s: %w", path, err)
	}
	// Create a tempfile to download to - if we're successful we'll move it to cache.
	tmp, err := os.CreateTemp("", "*")
	if err != nil {
		if err := tmp.Close(); err != nil {
			log.Printf("failed to close tempfile: %s", err)
		}
		if err := os.Remove(tmp.Name()); err != nil {
			log.Printf("failed to remove tempfile: %s", err)
		}
		return fmt.Errorf("failed to create tempfile: %w", err)
	}
	_, err = io.Copy(tmp, data)
	if err := tmp.Close(); err != nil {
		log.Printf("failed to close tempfile: %s", err)
	}
	if err != nil {
		if err := os.Remove(tmp.Name()); err != nil {
			log.Printf("failed to remove tempfile: %s", err)
		}
		return fmt.Errorf("failed to copy data: %w", err)
	}
	// Download OK - move to cache.
	if err := os.Rename(tmp.Name(), c.Path(path)); err != nil {
		if err := os.Remove(tmp.Name()); err != nil {
			log.Printf("failed to remove tempfile: %s", err)
		}
		return fmt.Errorf("failed to move tempfile to cache: %w", err)
	}
	return nil
}

type Library struct {
	BasePath string
	Cache    *Cache
}

func (l *Library) Path(path string) string {
	if filepath.IsAbs(path) {
		panic(fmt.Errorf("path cannot be an absolute path: %s", path)) // Coding error.
	}
	return filepath.Join(l.BasePath, path)
}

func (l *Library) EnsureDir(path string) error {
	dir := l.Path(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create directory %s: %w", dir, err)
	}
	return nil
}

func (l *Library) Abs(path string) string {
	path, err := filepath.Abs(l.Path(path))
	if err != nil {
		panic(err)
	}
	return path
}

func (l *Library) Link(file File) error {
	// ID is canonical location - create directory and link to cached file.
	id := filepath.Join(file.Service, file.User, "id", file.ID)
	if err := l.EnsureDir(id); err != nil {
		return fmt.Errorf("could not ensure directory %s: %w", filepath.Join(file.Service, file.User), err)
	}
	if err := os.Symlink(l.Cache.Abs(strings.TrimLeft(file.Path, "/")), l.Path(filepath.Join(id, file.Name))); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create symlink %s: %w", filepath.Join(id, file.Name), err)
	}

	// Others are alias locations - link to ID directory.

	// Link by title.
	title := filepath.Join(file.Service, file.User, "title", file.Title)
	if err := l.EnsureDir(filepath.Dir(title)); err != nil {
		return fmt.Errorf("could not ensure directory %s: %w", filepath.Join(title, file.Name), err)
	}
	if err := os.Symlink(l.Abs(id), l.Path(title)); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create symlink %s: %w", filepath.Join(id, file.Name), err)
	}

	// Link by time added.
	added := l.Path(filepath.Join(file.Service, file.User, "added", file.Added.String()))
	if err := l.EnsureDir(filepath.Dir(added)); err != nil {
		return fmt.Errorf("could not ensure directory %s: %w", filepath.Join(id, file.Name), err)
	}
	if err := os.Symlink(l.Abs(id), l.Path(added)); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create symlink %s: %w", filepath.Join(id, file.Name), err)
	}

	return nil
}
