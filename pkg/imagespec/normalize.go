package imagespec

import "strings"

// Normalize tries to transform a short image name into a normalized form
// A short image name is something like "alpine" which totally lacks the
// registry part.
func (spec *Spec) Normalize() *Spec {

	if spec.Registry == "" {
		spec.Registry = "docker.io" // TODO: configurable default
	}

	// TODO: apply registry-specific "normalisation" mechanisms in a generic
	// fashion
	if spec.Registry == "docker.io" {

		// "alpine" -> "library/alpine"
		if !strings.ContainsRune(spec.Repo, '/') {
			spec.Repo = "library/" + spec.Repo
		}

		if spec.Tag == "" {
			spec.Tag = "latest"
		}
	}

	return spec
}
