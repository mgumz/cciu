package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"

	"cciu/internal/imagespec"
	"cciu/internal/printer"
	"cciu/internal/registry"
	"cciu/internal/registry/fetcher"
	"cciu/internal/stats"
	"cciu/internal/tag"
)

type cciuRepoTags struct {
	Tags     []string
	Duration time.Duration
	FetchErr error
}

type cciuOpts struct {
	Filter struct {
		IgnoreBeta       bool
		StrictLabels     bool
		SkipNonSemVer    bool
		KeepVersionLevel int
	}

	Fetcher registry.Fetcher
	Printer printer.Printer
	Stats   *stats.AllStats
}

func main() {

	opts := &cciuOpts{Stats: &stats.AllStats{}}

	flag.BoolVar(&opts.Filter.IgnoreBeta, "exclude-beta-tags", false, "exclude 'beta' tags")
	flag.BoolVar(&opts.Filter.StrictLabels, "strict-labels", false, "strict label matching")
	flag.BoolVar(&opts.Filter.SkipNonSemVer, "skip-non-semver", false, "skip non-semver tags")

	keepVersionLevel := flag.String("keep", "", "keep [major|minor] version")
	doPrettyPrintJSON := flag.Bool("json-pretty", false, "indent json output")
	doPrintJSON := flag.Bool("json", false, "use json output format")
	doUseSimpleMarkers := flag.Bool("simple-markers", false, "use simple ascii markers")
	doShowStats := flag.Bool("stats", false, "show stats")
	doShowOldTags := flag.Bool("show-old", false, "show older tags")
	doLimitPerRegistry := flag.Int("limit-per-registry", 0, "limit parallel fetches per registry")
	fetchTimeout := flag.Duration("timeout", 0, "timeout for fetch operations")
	//authFilePath := flag.String("auth-file", "", "path to the credential store")
	doShowVersion := flag.Bool("version", false, "show version")

	flag.Usage = printUsage
	flag.Parse()

	if *doShowVersion {
		printVersion(*doPrintJSON || *doPrettyPrintJSON)
		return
	}

	switch *keepVersionLevel {
	case "":
	case "major":
		opts.Filter.KeepVersionLevel = keepMajorLevel
	case "minor":
		opts.Filter.KeepVersionLevel = keepMinorLevel
	default:
		os.Exit(printUnsupportedMinMajorLevel(*keepVersionLevel))
		return
	}

	opts.Printer = printer.NewTextPrinter(os.Stdout, *doUseSimpleMarkers)
	if *doPrintJSON || *doPrettyPrintJSON {
		jp := printer.NewJSONPrinter(os.Stdout)
		if *doPrettyPrintJSON {
			jp.SetIndent("  ", "  ")
		}
		opts.Printer = jp
	}
	opts.Printer.SetShowOldTags(*doShowOldTags)
	opts.Printer.SetShowStats(*doShowStats)

	opts.Fetcher = fetcher.NewSimple()
	if *doLimitPerRegistry > 0 {
		opts.Fetcher = fetcher.NewPerRegistry(*doLimitPerRegistry)
	}
	opts.Fetcher.SetTimeout(*fetchTimeout)
	//opts.Fetcher.SetAuthFilePath(*authFilePath)

	ts := time.Now()
	fetchAndCompare(flag.Args(), opts)
	opts.Stats.Duration = time.Since(ts)
	opts.Printer.Flush(opts.Stats)
}

type fetchedTags map[string]*cciuRepoTags

func fetchAndCompare(names []string, opts *cciuOpts) {

	// parse input arguments to check:
	// * if repos are given once only
	// * skip repos with no tag given

	once := make(map[string]bool)
	tags := make(fetchedTags)
	specs := make(imagespec.List, 0, len(names))

	wg := sync.WaitGroup{}

	stats := opts.Stats
	stats.Asked = len(names)

	for _, ref := range names {

		if _, ok := once[ref]; ok {
			stats.Duplicates++
			continue
		}
		once[ref] = true

		spec, err := imagespec.Parse(ref)
		if err != nil {
			stats.InvalidSpec++
			opts.Printer.NewSpec(ref, time.Duration(0), fmt.Errorf(errParsingName, ref, err))
			continue
		}

		// skip images without any tag
		// TODO: decide if print something
		if spec.Tag == "" {
			stats.NonTagged++
			continue
		}

		// skip images without semver tag
		if opts.Filter.SkipNonSemVer {
			_, err := semver.NewVersion(spec.Tag)
			if err != nil {
				stats.NonSemVer++
				err = fmt.Errorf(errTagNotSemver, spec.Tag, spec, err)
				opts.Printer.NewSpec(spec.String(), time.Duration(0), err)
				continue
			}
		}

		rr := spec.RegistryRepo()

		// skip repos given multiple times
		if _, fetched := tags[rr]; !fetched {

			rt := &cciuRepoTags{Tags: []string{}}
			tags[rr] = rt

			wg.Add(1)
			go func(rt *cciuRepoTags, s imagespec.Spec) {

				registry, name := s.Registry, s.StripContext().String()
				rt.Tags, rt.Duration, rt.FetchErr = opts.Fetcher.FetchTags(registry, name)

				stats.Fetch.Fetched++
				wg.Done()

			}(rt, *spec)
		}

		specs = append(specs, spec)
	}

	wg.Wait()

	for _, spec := range specs {
		compareAndPrint(spec, tags, opts)
	}
}

func compareAndPrint(spec *imagespec.Spec, rtags map[string]*cciuRepoTags, opts *cciuOpts) {

	prt, stats := opts.Printer, opts.Stats

	v, err := semver.NewVersion(spec.Tag)
	if err != nil {
		stats.NonSemVer++
		if !opts.Filter.SkipNonSemVer {
			err = fmt.Errorf(errTagNotSemver, spec.Tag, spec, err)
			prt.NewSpec(spec.String(), time.Duration(0), err)
		}
		return
	}

	rt := rtags[spec.RegistryRepo()]

	if rt.FetchErr != nil {
		err = fmt.Errorf(errFetchTags, spec, rt.FetchErr)
		prt.NewSpec(spec.String(), rt.Duration, err)
		return
	}

	fl := fList{}
	fl = fl.filterHugeVersionGaps(v)
	fl = fl.filterBetaVersions(opts.Filter.IgnoreBeta)
	fl = fl.filterStrictLabels(spec.Label, opts.Filter.StrictLabels)
	fl = fl.filterKeepLevel(v, opts.Filter.KeepVersionLevel)

	tags := tag.NewFromStrings(rt.Tags, tag.ApplyFilterList(fl))
	tags.Sort()
	tags.Reverse()

	prt.NewSpec(spec.String(), rt.Duration, nil)

	spec.Tag, spec.Label, spec.Context = "", "", ""

	for _, tag := range tags {
		prt.PrintTag(spec.String(), v, tag)
	}

	stats.Checked++
}
