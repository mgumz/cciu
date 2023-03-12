package imagespec

// Spec describes a container image spec.
//
// A container-image name is constructed from several pieces:
//
// docker.io/library/alpine:1.13-rc0-amd64+12345
// ^-+-----^ ^--+--^ ^--+-^ ^-+^ ^+------^ ^-+-^
//   |          |       |     |   |          |
//   +- registry|       |     |   |          |
//              +- repository |   |          |
//                            |   |          |
// ····· semver ···································
//                            |   |          |
//                            +- version     |
//                            |   |          |
//                            |   +- pre-release
//                            |              |
//                            |              +- build
// ····· other ····································
//                            |
//                            +- label
//
// relevant links:
// * https://cloud.google.com/artifact-registry/docs/docker/names
//   <project>/<repo>/<image>:<tag>
// * https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux_atomic_host/7/html/recommended_practices_for_container_development/naming
//   <repository>:<tag>
//
// * short names: https://github.com/containers/shortnames
//
type Spec struct {
	Registry string
	Repo     string
	Tag      string
	Label    string
	Context  string
}

// List defines a list of ImageSpecs
type List []*Spec

const (
	sepLabelContext = "@"
	sepTagLabel     = "-"
	sepRepoTag      = ":"
	sepRegistryRepo = "/"
)

// Parse parses the given name into an ImageSpec
func Parse(name string) (*Spec, error) {

	reg, repo, tag, label, ctx := SplitImageSpec(name)

	// TODO: validation

	s := &Spec{
		Registry: reg,
		Repo:     repo,
		Tag:      tag,
		Label:    label,
		Context:  ctx,
	}
	return s, nil
}

//func MustParse(ispec string) *ImageSpec {
// FIXME: not implemented yet
//	return nil
//}

// String satisfies the Stringer interface
func (spec *Spec) String() string {

	s := spec.RegistryRepo()
	if spec.Tag != "" {
		s += ":" + spec.Tag
	}
	if spec.Label != "" {
		s += "-" + spec.Label
	}
	if spec.Context != "" {
		s += sepLabelContext + spec.Context
	}
	return s
}

// RegistryRepo returns a string based upon the registry and the repo component
// of spec
func (spec *Spec) RegistryRepo() string {

	s := spec.Registry
	if spec.Registry != "" {
		s += "/"
	}
	s += spec.Repo

	return s
}

// StripContext removes the context component of spec. The function returns
// a new ImageSpec
func (spec *Spec) StripContext() *Spec {

	n := *spec
	n.Context = ""

	return &n
}
