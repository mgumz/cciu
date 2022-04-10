package main

import (
	"cciu/internal/tag"
	"fmt"
	"os"
	"strings"

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

func (list fList) filterVersionLevel(base *semver.Version, level string) fList {
	var c *semver.Constraints

	switch strings.ToLower(level) {
	case "":
		return list
	case "major":
		c, _ = semver.NewConstraint(fmt.Sprintf("~%d", base.Major()))
	case "minor":
		c, _ = semver.NewConstraint(fmt.Sprintf("~%d.%d", base.Major(), base.Minor()))
	default:
		fmt.Fprintf(os.Stderr, "Ignoring unknown version Level: %s\n", level)
		return list
	}
	f := tag.ConstraintFilter(c)

	return append(list, f)
}
