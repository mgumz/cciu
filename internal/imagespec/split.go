package imagespec

import (
	"strings"
)

// SplitLabelContext splits labelCtx into Label and Context (the part after
// "@", see sepLabelContext)
//
// FIXME: image@sha256:<sha256> is called image-by-digest. obviously this
// collides with the concept of "concept". either fix the concept / implementation
// of "a context" (which helps to remap the found container-images to the
// runnable entities they might refer to) OR detect a image-by-digest
// and ignore it coz we can't really compare such a digest against "newer"
func SplitLabelContext(labelCtx string) (label, ctx string) {

	if len(labelCtx) == 0 {
		return label, ctx
	}

	parts := strings.SplitN(labelCtx, sepLabelContext, 2)

	if len(parts) > 0 {
		label = parts[0]
	}
	if len(parts) > 1 {
		ctx = parts[1]
	}

	return label, ctx
}

// SplitTagLabel splits taglabel into tag and label, eg: 1.0-label into "1.0" and
// "label"
func SplitTagLabel(taglabel string) (tag, label string) {

	if len(taglabel) == 0 {
		return tag, label
	}

	parts := strings.SplitN(taglabel, sepTagLabel, 2)

	if len(parts) > 0 {
		tag = parts[0]
	}
	if len(parts) > 1 {
		label = parts[1]
	}

	return tag, label
}

// SplitRepoTag splits a string which contains a repo and a
// tag combination.
// SplitRepoTag("repo/image:tag-label")
// * repo: "repo/image"
// * tag:  "tag-label"
func SplitRepoTag(repotag string) (repo, tag string) {

	if len(repotag) == 0 {
		return repo, tag
	}

	parts := strings.SplitN(repotag, sepRepoTag, 2)

	if len(parts) > 0 {
		repo = parts[0]
	}
	if len(parts) > 1 {
		tag = parts[1]
	}

	return repo, tag
}

// SplitRegistryRepo splits regrepo into registry and the repo, eg: "quay.io/project/repo"
// into "quay.io" and "project/repo"
func SplitRegistryRepo(regrepo string) (registry, repotag string) {

	if len(regrepo) == 0 {
		return registry, repotag
	}

	parts := strings.SplitN(regrepo, sepRegistryRepo, 2)

	// if first part looks like a domain-name (contains '.'), do the actual
	// split. otherwise, the "registry"-part stays unknown.
	if len(parts) > 1 && strings.ContainsRune(parts[0], '.') {
		registry = parts[0]
		repotag = parts[1]
	} else {
		repotag = regrepo
	}

	return registry, repotag
}

// SplitImageSpec splits the given image specification into the various components
// of an ImageSpec
func SplitImageSpec(image string) (registry, repo, tag, label, ctx string) {

	// example.com:port/repo/image:tag1.2.3-label@ctx

	registry, repotag := SplitRegistryRepo(image)
	repo, taglabelctx := SplitRepoTag(repotag)
	taglabel, ctx := SplitLabelContext(taglabelctx)
	tag, label = SplitTagLabel(taglabel)

	return registry, repo, tag, label, ctx
}
