package main

import (
	"cciu/internal/tag"

	"github.com/Masterminds/semver/v3"
)

type fList []tag.FilterFunc

func (list fList) filterHugeVersionGaps(base *semver.Version) fList {
	f := tag.HugeVersionHeuristicFilter(base, 1000)
	return append(list, f)
}

func (list fList) filterBetaVersions(doFilter bool) fList {
	if !doFilter {
		return list
	}
	f := func(v *semver.Version) bool {
		return tag.IgnoreBetaVersions(v)
	}
	return append(list, f)
}

func (list fList) filterStrictLabels(label string, doFilter bool) fList {
	if !doFilter {
		return list
	}
	f := func(v *semver.Version) bool {
		return v.Prerelease() == label
	}
	return append(list, f)
}
