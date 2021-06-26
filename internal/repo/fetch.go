package repo

import (
	"context"
	"time"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
)

// FetchTags fetches the tags for the given repo as identified by name
func FetchTags(name string, timeout time.Duration, authFilePath string) ([]string, time.Duration, error) {

	ref, err := docker.ParseReference("//" + name)
	if err != nil {
		return []string{}, time.Duration(0), err
	}

	ctx, cancel := context.Background(), func() {}
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	}
	defer cancel()

	sys := types.SystemContext{}
	// TODO: use authFilePath in some form to log into the registry

	ts := time.Now()
	tags, err := docker.GetRepositoryTags(ctx, &sys, ref)

	return tags, time.Since(ts), err
}
