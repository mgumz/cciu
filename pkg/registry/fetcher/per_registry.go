package fetcher

import "time"

// PerRegistry implements a registry.Fetcher which allows only a limited amount
// of concurrent tag-fetch operations per named registry - a rate limiter
type PerRegistry struct {
	Simple
	limit    int
	fetchers map[string]chan *Simple
}

// NewPerRegistry returns a PerRegistry registry.Fetcher. limit conifgures the
// amount of concurrent tag-fetch operations per named registry
func NewPerRegistry(limit int) *PerRegistry {

	pr := &PerRegistry{
		limit:    limit,
		fetchers: map[string]chan *Simple{},
	}

	return pr
}

// SetTimeout sets the timeout for the fetch operation
func (pr *PerRegistry) SetTimeout(timeout time.Duration) { pr.timeout = timeout }

// SetAuthFilePath sets the path to the credential store
func (pr *PerRegistry) SetAuthFilePath(path string) { pr.authFilePath = path }

// FetchTags fetches the tags for the repo defined by name in the registry. It
// will limit the amount of concurrent tag-fetch operations per registry to
// what was configured via NewPerRegistry
func (pr *PerRegistry) FetchTags(registry, name string) ([]string, time.Duration, error) {

	fetchers, exists := pr.fetchers[registry]
	if !exists {
		fetchers = make(chan *Simple, pr.limit)
		for i := 0; i < pr.limit; i++ {
			s := NewSimple()
			s.SetTimeout(pr.timeout)
			s.SetAuthFilePath(pr.authFilePath)
			fetchers <- s
		}
		pr.fetchers[registry] = fetchers
	}

	simple := <-fetchers
	tags, dur, err := simple.FetchTags(registry, name)
	fetchers <- simple

	return tags, dur, err
}
