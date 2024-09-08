package printer

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/mgumz/cciu/pkg/stats"
)

const (
	markNeutral  = 0
	markAhead    = 1
	markEqual    = 2
	markOutdated = 3
)

var (
	verdictMarkersUnicode = []string{" ", "▲", "=", "▼"}
	verdictMarkersSimple  = []string{" ", "^", "=", "v"}
)

// TextPrinter is a small helper to print the requested images, fetched tags
// etc to create useful output
type TextPrinter struct {
	w              *tabwriter.Writer
	showOld        bool
	showStats      bool
	printedTag     bool
	verdictMarkers []string
}

// NewTextPrinter returns a TextPrinter which prints to w
// if simpleMarkers is true, use simple ASCII chars as vertdict markers
func NewTextPrinter(w io.Writer, simpleMarkers bool) *TextPrinter {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	p := &TextPrinter{
		w:              tw,
		verdictMarkers: verdictMarkersUnicode,
	}
	if simpleMarkers {
		p.verdictMarkers = verdictMarkersSimple
	}
	return p
}

// Flush finishes the output of the TextPrinter p
func (p *TextPrinter) Flush(stats *stats.AllStats) {
	if p.showStats && stats != nil {
		fmt.Fprintln(p.w, "---")
		fmt.Fprintf(p.w, "duration:\t%s\n", stats.Duration)
		fmt.Fprintf(p.w, "asked:\t%d\n", stats.Asked)
		fmt.Fprintf(p.w, "checked:\t%d\n", stats.Checked)
		fmt.Fprintf(p.w, "non-semver:\t%d\n", stats.NonSemVer)
		fmt.Fprintf(p.w, "duplicates:\t%d\n", stats.Duplicates)
	}
	p.w.Flush()
}

// SetShowOldTags configures TextPrinter p to also show older, outdated tags
func (p *TextPrinter) SetShowOldTags(s bool) {
	p.showOld = s
}

// SetShowStats activates printing statistics after the the list of checked
// tags
func (p *TextPrinter) SetShowStats(s bool) {
	p.showStats = s
}

// NewSpec starts printing the tags for "name" - its like a headline
func (p *TextPrinter) NewSpec(name string, dur time.Duration, err error) {
	p.printedTag = false
	comment := "\t# skipped"
	if dur > 0 {
		comment = fmt.Sprintf("\t# fetched in %s", dur.Round(time.Millisecond))
	}
	fmt.Fprintln(p.w, name, comment)
	if err != nil {
		fmt.Fprintf(p.w, "     %s\n", err)
		return
	}
}

// PrintTag prints the tag with version "other" for the requested "name" which
// was started via PrintSpec.
func (p *TextPrinter) PrintTag(name string, base, other *semver.Version) {

	if !p.showOld && p.printedTag {
		return
	}

	verdict := " "

	// force the "label" part to match base.
	// if we dont do this, library "semver" considers "-label" to be
	// a pre-release version and then "1.0.0-label" to be below "1.0.0"
	// still, we want to print the label later on. so, a copy of "other"
	// without the pre-release will do.
	o := *other
	o, _ = o.SetPrerelease("")
	b := *base

	// in case, "base" was given as "8.4" … the verdict
	// should be equal upon 8.4.1 or 8.4.99.
	bc, _ := semver.NewConstraint(b.Original())
	if bc.Check(&o) {
		verdict = p.verdictMarkers[markEqual]
	} else if o.GreaterThan(&b) {
		verdict = p.verdictMarkers[markAhead]
	} else if o.Equal(&b) {
		verdict = p.verdictMarkers[markEqual]
	} else if o.LessThan(&b) {
		verdict = p.verdictMarkers[markOutdated]
	}

	fmt.Fprintf(p.w, "%s    %s:%s\t\n", verdict, name, other)

	p.printedTag = true
}
