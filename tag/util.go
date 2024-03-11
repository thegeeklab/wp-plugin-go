package tag

import "strings"

// stripHeadPrefix removes the "refs/heads/" prefix from the given ref string.
// It is used to clean up Git ref names.
func stripHeadPrefix(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

// stripTagPrefix removes the "refs/tags/" prefix and any "v" prefix
// from the given ref string. It is used to clean up Git tag ref names.
func stripTagPrefix(ref string) string {
	ref = strings.TrimPrefix(ref, "refs/tags/")
	ref = strings.TrimPrefix(ref, "v")

	return ref
}

// IsTaggable checks whether tags should be created for the specified ref.
// The function returns true if the ref either matches the default branch
// or is a tag ref.
func IsTaggable(ref, defaultBranch string) bool {
	if strings.HasPrefix(ref, "refs/tags/") {
		return true
	}

	if stripHeadPrefix(ref) == defaultBranch {
		return true
	}

	return false
}
