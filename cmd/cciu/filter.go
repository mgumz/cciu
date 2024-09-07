package main

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/mgumz/cciu/pkg/tag"
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
	if keepLevel == tag.KeepMajor {
		cs := fmt.Sprintf("~%d", base.Major())
		c, _ := semver.NewConstraint(cs)
		return append(list, tag.ConstraintFilter(c))
	} else if keepLevel == tag.KeepMinor {
		cs := fmt.Sprintf("~%d.%d", base.Major(), base.Minor())
		c, _ := semver.NewConstraint(cs)
		return append(list, tag.ConstraintFilter(c))
	}
	return list
}
