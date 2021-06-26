package imagespec

import (
	"testing"
)

func TestSplitLabelContext(t *testing.T) {

	fixtures := [...]struct {
		LabelContext    string
		ExpectedLabel   string
		ExpectedContext string
	}{
		{"label", "label", ""},
		{"label@context", "label", "context"},
		{"", "", ""},
	}

	for _, f := range fixtures {
		label, ctx := SplitLabelContext(f.LabelContext)
		if f.ExpectedLabel != label {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedLabel, label)
		}
		if f.ExpectedContext != ctx {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedContext, ctx)
		}
	}
}

func TestSplitTagLabel(t *testing.T) {

	fixtures := [...]struct {
		TagLabel      string
		ExpectedTag   string
		ExpectedLabel string
	}{
		{"1.0.0", "1.0.0", ""},
		{"1.0.0-latest", "1.0.0", "latest"},
		{"1.0.0-really-latest", "1.0.0", "really-latest"},
		{"", "", ""},
	}

	for _, f := range fixtures {
		tag, label := SplitTagLabel(f.TagLabel)
		if f.ExpectedTag != tag {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedTag, tag)
		}
		if f.ExpectedLabel != label {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedLabel, label)
		}
	}
}

func TestSplitRepoTag(t *testing.T) {
	fixtures := [...]struct {
		RepoTag      string
		ExpectedRepo string
		ExpectedTag  string
	}{
		{"alpine", "alpine", ""},
		{"alpine:latest", "alpine", "latest"},
		{"alpine:3", "alpine", "3"},
		{"alpine:3.13", "alpine", "3.13"},
		{"alpine:3-dev", "alpine", "3-dev"},
		{"", "", ""},
	}

	for _, f := range fixtures {
		repo, tag := SplitRepoTag(f.RepoTag)

		t.Logf("%q => %q %q", f.RepoTag, repo, tag)

		if f.ExpectedRepo != repo {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedRepo, repo)
		}
		if f.ExpectedTag != tag {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedTag, tag)
		}
	}
}

func TestSplitRegistryRepo(t *testing.T) {
	fixtures := [...]struct {
		RegistryRepoTag  string
		ExpectedRegistry string
		ExpectedRepoTag  string
	}{
		{"alpine", "", "alpine"},
		{"alpine/latest", "", "alpine/latest"},
		{"repo/alpine:3", "", "repo/alpine:3"},
		{"example.com/alpine", "example.com", "alpine"},
		{"example.com/repo/alpine:3", "example.com", "repo/alpine:3"},
		{"example.com/repo/alpine:3.13", "example.com", "repo/alpine:3.13"},
		{"example.com:8080/repo/alpine:3.13", "example.com:8080", "repo/alpine:3.13"},
		{"", "", ""},
	}

	for _, f := range fixtures {
		registry, repotag := SplitRegistryRepo(f.RegistryRepoTag)

		t.Logf("%q => %q %q", f.RegistryRepoTag, registry, repotag)

		if f.ExpectedRegistry != registry {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedRegistry, registry)
		}
		if f.ExpectedRepoTag != repotag {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedRepoTag, repotag)
		}
	}
}

func TestSplitImageSpec(t *testing.T) {
	fixtures := [...]struct {
		Image            string
		ExpectedRegistry string
		ExpectedRepo     string
		ExpectedTag      string
		ExpectedLabel    string
		ExpectedContext  string
	}{
		{"alpine", "", "alpine", "", "", ""},
		{"alpine:3", "", "alpine", "3", "", ""},
		{"alpine:3.12", "", "alpine", "3.12", "", ""},
		{"alpine/latest", "", "alpine/latest", "", "", ""},
		{"repo/alpine:3", "", "repo/alpine", "3", "", ""},
		{"example.com/alpine", "example.com", "alpine", "", "", ""},
		{"example.com/repo/alpine:3", "example.com", "repo/alpine", "3", "", ""},
		{"example.com/repo/alpine:3.13", "example.com", "repo/alpine", "3.13", "", ""},
		{"example.com/repo/alpine:3.13-label@sample-ctx", "example.com", "repo/alpine", "3.13", "label", "sample-ctx"},
		{"example.com/repo/alpine:3.13@sample-ctx", "example.com", "repo/alpine", "3.13", "", "sample-ctx"},
		{"", "", "", "", "", ""},
	}

	for _, f := range fixtures {
		registry, repo, tag, label, ctx := SplitImageSpec(f.Image)

		t.Logf("%q => %q %q %q %q", f.Image, registry, repo, tag, label)

		if f.ExpectedRegistry != registry {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedRegistry, registry)
		}
		if f.ExpectedRepo != repo {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedRepo, repo)
		}
		if f.ExpectedTag != tag {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedTag, tag)
		}
		if f.ExpectedLabel != label {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedLabel, label)
		}

		if f.ExpectedContext != ctx {
			t.Fatalf("expected: %q, actual: %q", f.ExpectedContext, ctx)

		}
	}
}
