package registry

import "time"

// Fetcher defines the interface for a Registry.Fetcher to provide various
// ways to fetch the tags for a given name from different registries
type Fetcher interface {
	FetchTags(registry, name string) ([]string, time.Duration, error)

	// SetTimeout defines the timeout for fetch operations
	SetTimeout(timeout time.Duration)

	// SetAuthFilePath sets the path to the credential store
	SetAuthFilePath(path string)
}
