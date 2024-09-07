package tag

import (
	"math"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// FilterFunc describes a function which filters out semver-versions:
// a return value of
// * true - the entry stays
// * false - the entry gets filtered out
type FilterFunc func(*semver.Version) bool

// HugeVersionHeuristicFilter generates a tag.FilterFunc to filter out
// unnatural gaps the version numbering in some repos which use a date-format:
// 20210530. these versions are parsed into 20210530.0.0 and as a result are
// "bigger" than regular versions and trump those.
//
// n == 1000 seems to be a major, "unnatural" jump in
// the major version. lets call it a "heuristic".
func HugeVersionHeuristicFilter(a *semver.Version, limit int) FilterFunc {

	filter := func(b *semver.Version) bool {

		if limit < 0 {
			return false
		}

		if b.Major() <= a.Major() {
			return true
		}

		// https://github.com/securego/gosec/issues/1185
		// https://github.com/ccoVeille/go-safecast
		if (a.Major() >= math.MaxInt) || (b.Major() >= math.MaxInt) {
			return false
		}
		delta := int(b.Major()) - int(a.Major()) // #nosec G115

		return delta <= limit
	}

	return filter
}

// ConstraintFilter generates a tag.FilterFunc based on a semver.Constraint.
// The Check() method already has the same interface as required for FilterFunc,
// hence, just return the function.
func ConstraintFilter(c *semver.Constraints) FilterFunc {

	return c.Check
}

// IgnoreBetaVersions filters all labels which start with
// * "rc" - for release-candidate
// * "beta" - for beta-versions
// * "alpha" - for alpha-versions
func IgnoreBetaVersions(v *semver.Version) bool {

	p := v.Prerelease()

	switch {
	case strings.HasPrefix(p, "rc"):
	case strings.HasPrefix(p, "beta"):
	case strings.HasPrefix(p, "alpha"):
		return false
	}

	return true
}

// ApplyFilterList returns a tag.FilterFunc which executes all filters in
// list sequentially. The first `false` return value stops processing
// the list. Thus, logical "AND" is applied here.
func ApplyFilterList(list []FilterFunc) FilterFunc {

	filter := func(v *semver.Version) bool {
		for _, f := range list {
			if !f(v) {
				return false
			}
		}
		return true
	}

	return filter
}
