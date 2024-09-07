package tag

import (
	"slices"
	"sort"

	semver "github.com/Masterminds/semver/v3"
)

// List contains semantic versions
type List []*semver.Version

// Sort sorts the list of tags according to semantic versions
func (tags List) Sort() {
	sort.Sort(semver.Collection(tags))
}

// Reverse reverses the order of the TagList tags
func (tags List) Reverse() {
	slices.Reverse(tags)
}

// NewFromStrings creates a new List, based upon the string list "tags". In
// addition, it applies the filter function extraFilter
func NewFromStrings(tags []string, extraFilter FilterFunc) List {

	filtered := List{}

	for i := range tags {
		tv, err := semver.NewVersion(tags[i])
		if err != nil {
			continue
		}

		if !extraFilter(tv) {
			continue
		}

		filtered = append(filtered, tv)
	}

	return filtered
}
