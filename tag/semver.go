package tag

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// SemverTagSuffix returns a set of default suggested tags
// based on the commit ref with an attached suffix.
func SemverTagSuffix(ref, suffix string, strict bool) ([]string, error) {
	tags, err := SemverTags(ref, strict)
	if err != nil {
		return nil, err
	}

	if len(suffix) == 0 {
		return tags, nil
	}

	for i, tag := range tags {
		if tag == "latest" {
			tags[i] = suffix
		} else {
			tags[i] = fmt.Sprintf("%s-%s", tag, suffix)
		}
	}

	return tags, nil
}

// SemverTags returns a set of default suggested tags based on
// the commit ref.
func SemverTags(ref string, strict bool) ([]string, error) {
	var (
		version *semver.Version
		err     error
	)

	if !strings.HasPrefix(ref, "refs/tags/") {
		return []string{"latest"}, nil
	}

	rawVersion := stripTagPrefix(ref)

	version, err = semver.NewVersion(rawVersion)
	if strict {
		version, err = semver.StrictNewVersion(rawVersion)
	}

	if err != nil {
		return []string{"latest"}, err
	}

	if version.Prerelease() != "" {
		return []string{
			version.String(),
		}, nil
	}

	if version.Major() == 0 {
		return []string{
			fmt.Sprintf("%v.%v", version.Major(), version.Minor()),
			fmt.Sprintf("%v.%v.%v", version.Major(), version.Minor(), version.Patch()),
		}, nil
	}

	return []string{
		fmt.Sprintf("%v", version.Major()),
		fmt.Sprintf("%v.%v", version.Major(), version.Minor()),
		fmt.Sprintf("%v.%v.%v", version.Major(), version.Minor(), version.Patch()),
	}, nil
}
