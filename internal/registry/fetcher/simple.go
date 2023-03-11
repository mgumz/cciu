package fetcher

import (
	"time"

	"github.com/mgumz/cciu/internal/repo"
)

// Simple defines a simple registry.Fetcher which is just a tiny wrapper around
// repo.FetchTags.
type Simple struct {
	timeout      time.Duration
	authFilePath string
}

// NewSimple returns a Simple registry.Fetcher
func NewSimple() *Simple { return &Simple{} }

// SetTimeout sets the timeout for the fetch operation
func (s *Simple) SetTimeout(timeout time.Duration) { s.timeout = timeout }

// SetAuthFilePath sets the path to the credential store
func (s *Simple) SetAuthFilePath(path string) { s.authFilePath = path }

// FetchTags fetches the tags for name from registry. name is a full specified
// container name which includes the registry part.
func (s *Simple) FetchTags(registry, name string) ([]string, time.Duration, error) {
	return repo.FetchTags(name, s.timeout, s.authFilePath)
}
