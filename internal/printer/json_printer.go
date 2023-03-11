package printer

import (
	"encoding/json"
	"io"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/mgumz/cciu/internal/stats"
)

// JSONPrinter collects the requested images, the fetched tags and creates
// a JSON document which is at the end printed.
type JSONPrinter struct {
	w         io.Writer
	enc       *json.Encoder
	jo        jsonOutput
	cur       *jsonImage
	showOld   bool
	showStats bool
}

// NewJSONPrinter returns a JSONPrinter which prints the JSON to w
func NewJSONPrinter(w io.Writer) *JSONPrinter {
	p := &JSONPrinter{w: w, jo: jsonOutput{}}
	p.enc = json.NewEncoder(w)
	return p
}

// SetIndent activates indentation when printing the JSON document
func (p *JSONPrinter) SetIndent(prefix, indent string) *JSONPrinter {
	p.enc.SetIndent(prefix, indent)
	return p
}

// Flush encodes the internal JSON document to the configured Writer
func (p *JSONPrinter) Flush(stats *stats.AllStats) {

	if p.cur != nil {
		p.jo.Images = append(p.jo.Images, *p.cur)
		if p.showStats {
			p.jo.Stats = stats
		}
		p.cur = nil
	}

	_ = p.enc.Encode(p.jo)
}

// SetShowOldTags activates showing old tags while printing
func (p *JSONPrinter) SetShowOldTags(s bool) {
	p.showOld = s
}

// SetShowStats activates putting a statistic object into the printed JSON
func (p *JSONPrinter) SetShowStats(s bool) {
	p.showStats = s
}

type jsonOutput struct {
	Images []jsonImage     `json:"images"`
	Stats  *stats.AllStats `json:"stats,omitempty"`
}

type jsonImage struct {
	Requested string `json:"requested"`
	Verdict   string `json:"verdict"` // "ahead", "current", "outdated"

	Tags     []jsonTag     `json:"tags"`
	Duration time.Duration `json:"duration"`
	Err      error         `json:"error,omitempty"`
}

type jsonTag struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Verdict string `json:"verdict"` // "ahead", "current", "outdated"
}

// NewSpec starts collecting the tags for "name"
func (p *JSONPrinter) NewSpec(name string, dur time.Duration, err error) {

	if p.cur != nil {
		p.jo.Images = append(p.jo.Images, *p.cur)
	}

	p.cur = &jsonImage{
		Requested: name,
		Tags:      []jsonTag{},
		Duration:  dur,
	}
	if err != nil {
		p.cur.Err = err
	}
}

// PrintTag stores the tag with version "other" for the requested "name" which
// was started via PrintSpec
func (p *JSONPrinter) PrintTag(name string, base, other *semver.Version) {

	if !p.showOld && len(p.cur.Tags) > 0 {
		return
	}

	// force the "label" part to match base.
	// if we dont do this, library "semver" considers "-label" to be
	// a pre-release version and then "1.0.0-label" to be below "1.0.0"
	// still, we want to print the label later on. so, a copy of "other"
	// without the pre-release will do.
	o := *other
	o, _ = o.SetPrerelease("")

	vbase, vother := "", ""
	if o.GreaterThan(base) {
		vbase, vother = "outdated", "ahead"
	} else if o.Equal(base) {
		vbase, vother = "equal", "equal"
	} else if o.LessThan(base) {
		vbase, vother = "ahead", "outdated"
	}

	if p.cur.Verdict == "" {
		p.cur.Verdict = vbase
	}

	tag := jsonTag{
		Name:    name + ":" + other.String(),
		Version: other.String(),
		Verdict: vother,
	}
	p.cur.Tags = append(p.cur.Tags, tag)
}
