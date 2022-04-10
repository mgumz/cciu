package main

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"cciu/internal/tag"
)

const (
	keepIgnoreLevel = iota
	keepMajorLevel
	keepMinorLevel
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

func (list fList) filterKeepLevel(base *semver.Version, keepLevel int) fList {
	if keepLevel == keepMajorLevel {
		c, _ := semver.NewConstraint(fmt.Sprintf("~%d", base.Major()))
		return append(list, tag.ConstraintFilter(c))
	} else if keepLevel == keepMinorLevel {
		c, _ := semver.NewConstraint(fmt.Sprintf("~%d.%d", base.Major(), base.Minor()))
		return append(list, tag.ConstraintFilter(c))
	}
	return list
}
