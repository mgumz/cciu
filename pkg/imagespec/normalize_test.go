package imagespec

import (
	"testing"
)

func TestNormalize(t *testing.T) {

	fixtures := [...]struct {
		InSpec       Spec
		ExpectedSpec Spec
	}{
		{Spec{"", "walkerlee/nsenter", "", "", ""}, Spec{"docker.io", "walkerlee/nsenter", "latest", "", ""}},
		{Spec{"", "alpine", "latest", "", ""}, Spec{"docker.io", "library/alpine", "latest", "", ""}},
		{Spec{"", "alpine", "", "", ""}, Spec{"docker.io", "library/alpine", "latest", "", ""}},
	}

	for _, f := range fixtures {

		in, expected := f.InSpec, f.ExpectedSpec
		normalized := in.Normalize()

		t.Logf("%q => %q %q", &f.InSpec, normalized, &expected)

	}

}
