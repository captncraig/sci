// Package sci contains core types and interfaces for Sierra Creative Interpreter tool
package sci

import (
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
)

type Loader interface {
	GetFile(name string) ([]byte, error)
}

func NewFromURL(base string) Loader {
	return &cachingLoader{
		cache: map[string][]byte{},
		inner: fromURL(base),
	}
}

func NewFromDir(base string) Loader {
	return &cachingLoader{
		cache: map[string][]byte{},
		inner: directory(base),
	}
}

type cachingLoader struct {
	cache map[string][]byte
	inner Loader
}

func (c *cachingLoader) GetFile(name string) ([]byte, error) {
	if d, ok := c.cache[name]; ok {
		return d, nil
	}
	return c.inner.GetFile(name)
}

type fromURL string

func (f fromURL) GetFile(name string) ([]byte, error) {
	resp, err := http.Get(path.Join(string(f), name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type directory string

func (d directory) GetFile(name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(string(d), name))
}
