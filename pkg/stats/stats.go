package stats

import "time"

// AllStats is used to collect all kind of stats
type AllStats struct {
	Duration time.Duration

	Asked       int
	Checked     int
	Duplicates  int
	NonSemVer   int
	NonTagged   int
	InvalidSpec int

	Fetch FetchStats
}

// FetchStats is used to collect all kind of fetch-related stats
type FetchStats struct {
	Duration time.Duration
	Fetched  int
}
