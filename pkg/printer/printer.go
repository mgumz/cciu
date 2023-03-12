package printer

import (
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/mgumz/cciu/pkg/stats"
)

// Printer describes the interface for a cciu printer - a helper for controlled
// output of fetched container image tags + the evaluation of these.
type Printer interface {
	SetShowOldTags(bool)
	SetShowStats(bool)
	NewSpec(name string, dur time.Duration, err error)
	PrintTag(name string, base, other *semver.Version)
	Flush(stats *stats.AllStats)
}
